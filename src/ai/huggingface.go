package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const huggingFaceAPIBaseURL = "https://api-inference.huggingface.co/models/"

type HuggingFaceProvider struct {
	client  *http.Client
	apiKey  string
	model   string
	baseURL string
}

type huggingFaceRequest struct {
	Inputs string `json:"inputs"`
}

type huggingFaceResponse []struct {
	GeneratedText string `json:"generated_text"`
}

func NewHuggingFaceProvider(apiKey, modelName string) (*HuggingFaceProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Hugging Face API key is required")
	}
	return &HuggingFaceProvider{
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
		apiKey:  apiKey,
		model:   modelName,
		baseURL: huggingFaceAPIBaseURL,
	}, nil
}

func (p *HuggingFaceProvider) Generate(ctx context.Context, prompt string) (string, error) {
	payload := huggingFaceRequest{Inputs: prompt}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer "+p.apiKey)

	url := p.baseURL + p.model
	respBody, statusCode, err := doAPIRequest(ctx, p.client, "POST", url, headers, payload)
	if err != nil {
		return "", fmt.Errorf("failed to send request to huggingface: %w", err)
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status from huggingface: %s", string(respBody))
	}

	var apiResp huggingFaceResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse huggingface json response: %w", err)
	}

	if len(apiResp) == 0 || apiResp[0].GeneratedText == "" {
		return "", fmt.Errorf("received an empty response from huggingface")
	}

	return strings.TrimPrefix(apiResp[0].GeneratedText, prompt), nil
}
