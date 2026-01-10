package llm_guard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	APIKey   string
	Endpoint string
	Model    string
}

type Client struct {
	cfg        Config
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	}
	if cfg.Model == "" {
		cfg.Model = "qwen-turbo"
	}
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) SetAPIKey(key string) {
	c.cfg.APIKey = key
}

func (c *Client) Check(prompt string, instruction string) (bool, string, error) {
	if c.cfg.APIKey == "" {
		return true, "", nil
	}

	systemPrompt := fmt.Sprintf("You are a content safety guardrail. %s\nIf the content is unsafe, respond with 'UNSAFE: <reason>'. If safe, respond with 'SAFE'.", instruction)

	payload := map[string]interface{}{
		"model": c.cfg.Model,
		"input": map[string]interface{}{
			"messages": []map[string]string{
				{"role": "system", "content": systemPrompt},
				{"role": "user", "content": prompt},
			},
		},
	}

	jsonBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", c.cfg.Endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return false, "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, "", fmt.Errorf("LLM API error: %d", resp.StatusCode)
	}

	var result struct {
		Output struct {
			Text string `json:"text"`
		} `json:"output"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, "", err
	}

	text := result.Output.Text
	if strings.HasPrefix(text, "UNSAFE") {
		reason := "content unsafe"
		if len(text) > 6 {
			reason = strings.TrimSpace(text[6:])
			if strings.HasPrefix(reason, ":") {
				reason = strings.TrimSpace(reason[1:])
			}
		}
		return false, reason, nil
	}

	return true, "", nil
}
