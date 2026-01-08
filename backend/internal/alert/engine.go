package alert

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// EnhancedEngine 增强版告警引擎
type EnhancedEngine struct {
	mu         sync.RWMutex
	store      *RuleStore
	dispatcher *NotifyDispatcher
	cooldowns  map[string]time.Time // rule_id -> last_triggered
	eventCh    chan Event
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewEnhancedEngine 创建增强告警引擎
func NewEnhancedEngine(store *RuleStore, dispatcher *NotifyDispatcher) *EnhancedEngine {
	ctx, cancel := context.WithCancel(context.Background())
	return &EnhancedEngine{
		store:      store,
		dispatcher: dispatcher,
		cooldowns:  make(map[string]time.Time),
		eventCh:    make(chan Event, 500),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start 启动告警引擎
func (e *EnhancedEngine) Start() {
	go e.processLoop()
}

// Stop 停止告警引擎
func (e *EnhancedEngine) Stop() {
	e.cancel()
}

// Submit 提交事件
func (e *EnhancedEngine) Submit(event Event) {
	select {
	case e.eventCh <- event:
	default:
		// Channel full, drop event
	}
}

func (e *EnhancedEngine) processLoop() {
	for {
		select {
		case <-e.ctx.Done():
			return
		case event := <-e.eventCh:
			e.evaluate(event)
		}
	}
}

func (e *EnhancedEngine) evaluate(event Event) {
	// 获取匹配的规则
	rules, err := e.store.ListRules(nil, true)
	if err != nil {
		return
	}

	for _, rule := range rules {
		if e.matchRule(rule, event) && e.checkCooldown(rule.ID, rule.CooldownSec) {
			e.trigger(rule, event)
		}
	}
}

func (e *EnhancedEngine) matchRule(rule AlertRule, event Event) bool {
	// 检查事件类型
	if len(rule.EventTypes) > 0 {
		matched := false
		for _, et := range rule.EventTypes {
			if et == event.Type || et == event.Reason {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 检查严重等级
	if rule.SeverityThreshold != "" {
		if !severityMet(event.Severity, rule.SeverityThreshold) {
			return false
		}
	}

	return true
}

func severityMet(eventSev, threshold string) bool {
	levels := map[string]int{
		"critical": 4,
		"high":     3,
		"medium":   2,
		"low":      1,
	}
	return levels[eventSev] >= levels[threshold]
}

func (e *EnhancedEngine) checkCooldown(ruleID string, cooldownSec int) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	lastTriggered, exists := e.cooldowns[ruleID]
	if !exists || time.Since(lastTriggered) > time.Duration(cooldownSec)*time.Second {
		e.cooldowns[ruleID] = time.Now()
		return true
	}
	return false
}

func (e *EnhancedEngine) trigger(rule AlertRule, event Event) {
	eventData, _ := json.Marshal(event)

	history := &AlertHistory{
		RuleID:    &rule.ID,
		RuleName:  rule.Name,
		TenantID:  rule.TenantID,
		Severity:  event.Severity,
		Title:     rule.Name + ": " + event.Reason,
		Message:   event.Reason,
		EventData: eventData,
	}

	// 分发通知
	status := e.dispatcher.Dispatch(e.ctx, history, rule.NotifyChannels)
	statusJSON, _ := json.Marshal(status)
	history.NotifyStatus = statusJSON

	// 保存历史
	_ = e.store.SaveHistory(history)
}
