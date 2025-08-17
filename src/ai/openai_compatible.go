package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAICompatibleProvider struct {
	client  *http.Client
	apiKey  string
	model   string
	baseURL string
}

type openAICompatRequest struct {
	Model    string                `json:"model"`
	Messages []openAICompatMessage `json:"messages"`
}

type openAICompatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAICompatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func NewOpenAICompatibleProvider(apiKey, modelName, baseURL string) (*OpenAICompatibleProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for OpenAI Compatible provider")
	}
	if baseURL == "" {
		return nil, fmt.Errorf("endpoint URL is required for OpenAI Compatible provider")
	}

	return &OpenAICompatibleProvider{
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
		apiKey:  apiKey,
		model:   modelName,
		baseURL: baseURL,
	}, nil
}

func (p *OpenAICompatibleProvider) Generate(ctx context.Context, prompt string) (string, error) {
	payload := openAICompatRequest{
		Model: p.model,
		Messages: []openAICompatMessage{
			{Role: "user", Content: prompt},
		},
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer "+p.apiKey)

	respBody, statusCode, err := doAPIRequest(ctx, p.client, "POST", p.baseURL+"/chat/completions", headers, payload)
	if err != nil {
		return "", fmt.Errorf("failed to send request to compatible endpoint: %w", err)
	}

	var apiResp openAICompatResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse json response: %w", err)
	}

	if statusCode != http.StatusOK {
		if apiResp.Error != nil {
			return "", fmt.Errorf("api error (type: %s): %s", apiResp.Error.Type, apiResp.Error.Message)
		}
		return "", fmt.Errorf("received non-200 status from endpoint: %d", statusCode)
	}

	if len(apiResp.Choices) == 0 || apiResp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("received an empty or invalid response from the endpoint")
	}

	return apiResp.Choices[0].Message.Content, nil
}
