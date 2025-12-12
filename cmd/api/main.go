package main

import (
	"log"
	"net/http"
	"time"

	"aiguardrails/internal/agent"
	"aiguardrails/internal/audit"
	"aiguardrails/internal/config"
	"aiguardrails/internal/mcp"
	"aiguardrails/internal/policy"
	"aiguardrails/internal/promptfw"
	"aiguardrails/internal/rag"
	"aiguardrails/internal/server"
	"aiguardrails/internal/store"
	"aiguardrails/internal/tenant"
	"aiguardrails/internal/usage"
)

func main() {
	cfg := config.Default()

	db, err := store.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	if err := store.Migrate(db, "migrations"); err != nil {
		log.Fatalf("migrations error: %v", err)
	}

	redisClient, err := store.NewRedis(cfg.RedisURL)
	if err != nil {
		log.Fatalf("redis connect error: %v", err)
	}

	tenantSvc := tenant.NewPGService(db)
	policyEng := policy.NewPGEngine(db)
	firewall := promptfw.NewFirewall(policyEng)
	llmClient := policy.NewQwenClient(cfg.QwenAPIBase, cfg.QwenAPIToken, cfg.QwenModel, cfg.QwenTimeoutSec, cfg.QwenRetries)
	llmDet := policy.NewLLMDetector(llmClient, cfg.LLMQueueSize, time.Duration(cfg.LLMCacheTTLMin)*time.Minute, cfg.QwenRPS)
	llmDet.Start(cfg.LLMWorkers)
	firewall.WithLLM(llmDet, cfg.OutputMode)
	capStore := mcp.NewStore(db)
	mcpBroker := mcp.NewBroker(policyEng, capStore)
	agentGw := agent.NewGateway(policyEng, firewall)
	ragSec := rag.NewSecurity(policyEng)
	usageMeter := usage.NewMeter()
	rateLimiter := usage.NewRateLimiter(redisClient)
	auditLog := audit.NewLogger()
	auditStore := audit.NewStore(db)
	if err := auditStore.Init(); err != nil {
		log.Fatalf("audit store init error: %v", err)
	}

	srv := server.New(cfg, tenantSvc, policyEng, firewall, agentGw, ragSec, usageMeter, rateLimiter, auditLog, auditStore, mcpBroker, capStore)
	log.Printf("starting API on %s", srv.Addr())
	if err := http.ListenAndServe(srv.Addr(), srv.Handler()); err != nil {
		log.Fatal(err)
	}
}
