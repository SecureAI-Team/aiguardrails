package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"aiguardrails/internal/types"
)

// QwenClient calls Aliyun DashScope moderation API.
type QwenClient struct {
	BaseURL string
	Token   string
	Model   string
	Timeout time.Duration
	Retries int
	http    *http.Client
}

// NewQwenClient constructs QwenClient.
func NewQwenClient(baseURL, token, model string, timeoutSec, retries int) *QwenClient {
	return &QwenClient{
		BaseURL: baseURL,
		Token:   token,
		Model:   model,
		Timeout: time.Duration(timeoutSec) * time.Second,
		Retries: retries,
		http:    &http.Client{},
	}
}

type qwenReq struct {
	Input map[string]string `json:"input"`
	Model string            `json:"model"`
}

type qwenResp struct {
	Output struct {
		Unsafe bool     `json:"unsafe"`
		Labels []string `json:"labels"`
	} `json:"output"`
}

// Moderate calls Qwen; if token missing, returns allow.
func (c *QwenClient) Moderate(ctx context.Context, text string) (types.GuardrailResult, error) {
	if c.Token == "" {
		return types.GuardrailResult{Allowed: true, Reason: "qwen_token_missing"}, nil
	}
	body, _ := json.Marshal(qwenReq{
		Input: map[string]string{"text": text},
		Model: c.Model,
	})
	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, bytes.NewReader(body))
	if err != nil {
		return types.GuardrailResult{Allowed: true, Reason: "qwen_req_error"}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)
	var lastErr error
	attempts := c.Retries + 1
	for i := 0; i < attempts; i++ {
		resp, err := c.http.Do(req.Clone(ctx))
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			continue
		}
		defer resp.Body.Close()
		var out qwenResp
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			continue
		}
		if out.Output.Unsafe {
			return types.GuardrailResult{Allowed: false, Reason: "qwen_block", Signals: out.Output.Labels}, nil
		}
		return types.GuardrailResult{Allowed: true, Reason: "qwen_allow", Signals: out.Output.Labels}, nil
	}
	return types.GuardrailResult{Allowed: true, Reason: "qwen_http_error"}, lastErr
}

var _ LLMClient = (*QwenClient)(nil)

