package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	anthropicAPIBaseURL = "https://api.anthropic.com/v1"
	anthropicAPIVersion = "2023-06-01"
)

type AnthropicProvider struct {
	client  *http.Client
	apiKey  string
	model   string
	baseURL string
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	Messages  []anthropicMessage `json:"messages"`
	MaxTokens int                `json:"max_tokens"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func NewAnthropicProvider(apiKey, modelName string) (*AnthropicProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("anthropic API key is required")
	}

	return &AnthropicProvider{
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
		apiKey:  apiKey,
		model:   modelName,
		baseURL: anthropicAPIBaseURL,
	}, nil
}

func (p *AnthropicProvider) Generate(ctx context.Context, prompt string) (string, error) {
	payload := anthropicRequest{
		Model: p.model,
		Messages: []anthropicMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens: 4096,
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("x-api-key", p.apiKey)
	headers.Set("anthropic-version", anthropicAPIVersion)

	respBody, statusCode, err := doAPIRequest(ctx, p.client, "POST", p.baseURL+"/messages", headers, payload)
	if err != nil {
		return "", fmt.Errorf("failed to send request to anthropic: %w", err)
	}

	var apiResp anthropicResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse anthropic json response: %w", err)
	}

	if statusCode != http.StatusOK {
		if apiResp.Error != nil {
			return "", fmt.Errorf("anthropic api error (type: %s): %s", apiResp.Error.Type, apiResp.Error.Message)
		}
		return "", fmt.Errorf("received non-200 status from anthropic: %d", statusCode)
	}

	if len(apiResp.Content) == 0 || apiResp.Content[0].Text == "" {
		return "", fmt.Errorf("received an empty or invalid response from anthropic")
	}

	return apiResp.Content[0].Text, nil
}
