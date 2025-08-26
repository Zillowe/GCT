package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	openRouterAPIBaseURL = "https://openrouter.ai/api/v1"
	gctAppName           = "https://zillowe.qzz.io/docs/zds/gct"
)

type OpenRouterProvider struct {
	client  *http.Client
	apiKey  string
	model   string
	baseURL string
}

type openRouterRequest struct {
	Model    string              `json:"model"`
	Messages []openRouterMessage `json:"messages"`
}

type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func NewOpenRouterProvider(apiKey, modelName string) (*OpenRouterProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenRouter API key is required")
	}

	return &OpenRouterProvider{
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
		apiKey:  apiKey,
		model:   modelName,
		baseURL: openRouterAPIBaseURL,
	}, nil
}

func (p *OpenRouterProvider) Generate(ctx context.Context, prompt string) (string, error) {
	payload := openRouterRequest{
		Model: p.model,
		Messages: []openRouterMessage{
			{Role: "user", Content: prompt},
		},
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer "+p.apiKey)
	headers.Set("HTTP-Referer", gctAppName)
	headers.Set("X-Title", "GCT AI Commit")

	respBody, statusCode, err := doAPIRequest(ctx, p.client, "POST", p.baseURL+"/chat/completions", headers, payload)
	if err != nil {
		return "", fmt.Errorf("failed to send request to openrouter: %w", err)
	}

	var apiResp openRouterResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse openrouter json response: %w", err)
	}

	if statusCode != http.StatusOK {
		if apiResp.Error != nil {
			return "", fmt.Errorf("openrouter api error (%d): %s", statusCode, apiResp.Error.Message)
		}
		return "", fmt.Errorf("received non-200 status from openrouter: %d", statusCode)
	}

	if len(apiResp.Choices) == 0 || apiResp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("received an empty or invalid response from openrouter")
	}

	return apiResp.Choices[0].Message.Content, nil
}
