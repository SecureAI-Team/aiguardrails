package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client is a lightweight SDK for integrating via HTTP.
type Client struct {
	BaseURL string
	AppID   string
	Secret  string
	http    *http.Client
}

// NewClient constructs a Client.
func NewClient(baseURL, appID, secret string) *Client {
	return &Client{
		BaseURL: baseURL,
		AppID:   appID,
		Secret:  secret,
		http:    &http.Client{},
	}
}

// PromptCheck calls the prompt firewall.
func (c *Client) PromptCheck(prompt string) (map[string]any, error) {
	payload := map[string]string{"prompt": prompt}
	return c.post("/v1/guardrails/prompt-check", payload)
}

func (c *Client) post(path string, body interface{}) (map[string]any, error) {
	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Id", c.AppID)
	req.Header.Set("X-App-Secret", c.Secret)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

