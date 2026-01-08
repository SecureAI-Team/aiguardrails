package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Notifier 通知接口
type Notifier interface {
	Send(ctx context.Context, alert *AlertHistory) error
	Type() string
}

// NotifyDispatcher 通知分发器
type NotifyDispatcher struct {
	mu        sync.RWMutex
	notifiers map[string]Notifier
	store     *RuleStore
}

// NewNotifyDispatcher 创建通知分发器
func NewNotifyDispatcher(store *RuleStore) *NotifyDispatcher {
	return &NotifyDispatcher{
		notifiers: make(map[string]Notifier),
		store:     store,
	}
}

// Register 注册通知器
func (d *NotifyDispatcher) Register(n Notifier) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.notifiers[n.Type()] = n
}

// Dispatch 分发通知
func (d *NotifyDispatcher) Dispatch(ctx context.Context, alert *AlertHistory, channels []string) map[string]string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	status := make(map[string]string)
	for _, ch := range channels {
		if notifier, ok := d.notifiers[ch]; ok {
			if err := notifier.Send(ctx, alert); err != nil {
				status[ch] = "failed: " + err.Error()
			} else {
				status[ch] = "sent"
			}
		} else {
			status[ch] = "notifier_not_found"
		}
	}
	return status
}

// SMSNotifier 短信通知器
type SMSNotifier struct {
	Provider   string
	AccessKey  string
	SecretKey  string
	SignName   string
	TemplateID string
}

func NewSMSNotifier(provider, accessKey, secretKey, signName, templateID string) *SMSNotifier {
	return &SMSNotifier{
		Provider:   provider,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		SignName:   signName,
		TemplateID: templateID,
	}
}

func (n *SMSNotifier) Type() string { return "sms" }

func (n *SMSNotifier) Send(ctx context.Context, alert *AlertHistory) error {
	// 这里实现实际的短信发送逻辑
	// 根据Provider选择阿里云或腾讯云SDK
	fmt.Printf("[SMS] Sending alert: %s - %s\n", alert.Title, alert.Severity)

	// TODO: 从notify_recipients获取手机号并发送
	// 实际实现需要调用短信SDK
	return nil
}

// WeChatNotifier 微信通知器
type WeChatNotifier struct {
	AppID       string
	AppSecret   string
	TemplateID  string
	accessToken string
	tokenExpire time.Time
	mu          sync.Mutex
}

func NewWeChatNotifier(appID, appSecret, templateID string) *WeChatNotifier {
	return &WeChatNotifier{
		AppID:      appID,
		AppSecret:  appSecret,
		TemplateID: templateID,
	}
}

func (n *WeChatNotifier) Type() string { return "wechat" }

func (n *WeChatNotifier) Send(ctx context.Context, alert *AlertHistory) error {
	// 获取access_token
	token, err := n.getAccessToken()
	if err != nil {
		return err
	}

	// 发送模板消息
	fmt.Printf("[WeChat] Sending alert with token %s...: %s\n", token[:10], alert.Title)

	// TODO: 从notify_recipients获取openid并发送模板消息
	return nil
}

func (n *WeChatNotifier) getAccessToken() (string, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.accessToken != "" && time.Now().Before(n.tokenExpire) {
		return n.accessToken, nil
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		n.AppID, n.AppSecret)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("wechat error: %s", result.ErrMsg)
	}

	n.accessToken = result.AccessToken
	n.tokenExpire = time.Now().Add(time.Duration(result.ExpiresIn-60) * time.Second)
	return n.accessToken, nil
}

// WeComNotifier 企业微信机器人通知器
type WeComNotifier struct {
	WebhookURL string
}

func NewWeComNotifier(webhookURL string) *WeComNotifier {
	return &WeComNotifier{WebhookURL: webhookURL}
}

func (n *WeComNotifier) Type() string { return "wecom" }

func (n *WeComNotifier) Send(ctx context.Context, alert *AlertHistory) error {
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": fmt.Sprintf("## 🚨 安全告警\n**级别**: %s\n**标题**: %s\n**详情**: %s\n**时间**: %s",
				alert.Severity, alert.Title, alert.Message, alert.CreatedAt.Format("2006-01-02 15:04:05")),
		},
	}
	body, _ := json.Marshal(msg)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// DingTalkNotifier 钉钉机器人通知器
type DingTalkNotifier struct {
	WebhookURL string
	Secret     string
}

func NewDingTalkNotifier(webhookURL, secret string) *DingTalkNotifier {
	return &DingTalkNotifier{WebhookURL: webhookURL, Secret: secret}
}

func (n *DingTalkNotifier) Type() string { return "dingtalk" }

func (n *DingTalkNotifier) Send(ctx context.Context, alert *AlertHistory) error {
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "安全告警: " + alert.Title,
			"text": fmt.Sprintf("## 🚨 安全告警\n- **级别**: %s\n- **标题**: %s\n- **详情**: %s\n- **时间**: %s",
				alert.Severity, alert.Title, alert.Message, alert.CreatedAt.Format("2006-01-02 15:04:05")),
		},
	}
	body, _ := json.Marshal(msg)

	// TODO: 如果有secret，需要添加签名
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// WebhookNotifier 通用Webhook通知器
type WebhookNotifier struct {
	URL     string
	Headers map[string]string
}

func NewWebhookNotifier(url string, headers map[string]string) *WebhookNotifier {
	return &WebhookNotifier{URL: url, Headers: headers}
}

func (n *WebhookNotifier) Type() string { return "webhook" }

func (n *WebhookNotifier) Send(ctx context.Context, alert *AlertHistory) error {
	body, _ := json.Marshal(alert)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range n.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
