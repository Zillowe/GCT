package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const azureAPIVersion = "2024-02-01"

type AzureProvider struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

type azureRequest struct {
	Model    string         `json:"model"`
	Messages []azureMessage `json:"messages"`
}
type azureMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type azureResponse struct {
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

func NewAzureProvider(apiKey, resourceName, deploymentName string) (*AzureProvider, error) {
	if apiKey == "" || resourceName == "" || deploymentName == "" {
		return nil, fmt.Errorf("API Key, Azure Resource Name, and Deployment Name (in 'model' field) are all required for Azure OpenAI")
	}

	url := fmt.Sprintf("https://%s.openai.azure.com/openai/deployments/%s/chat/completions?api-version=%s",
		resourceName, deploymentName, azureAPIVersion)

	return &AzureProvider{
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
		apiKey:  apiKey,
		baseURL: url,
	}, nil
}

func (p *AzureProvider) Generate(ctx context.Context, prompt string) (string, error) {
	payload := azureRequest{
		Model: "",
		Messages: []azureMessage{
			{Role: "user", Content: prompt},
		},
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("api-key", p.apiKey)

	respBody, statusCode, err := doAPIRequest(ctx, p.client, "POST", p.baseURL, headers, payload)
	if err != nil {
		return "", fmt.Errorf("failed to send request to azure: %w", err)
	}

	var apiResp azureResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse azure json response: %w", err)
	}

	if statusCode != http.StatusOK {
		if apiResp.Error != nil {
			return "", fmt.Errorf("azure api error (type: %s): %s", apiResp.Error.Type, apiResp.Error.Message)
		}
		return "", fmt.Errorf("received non-200 status from azure: %s", string(respBody))
	}

	if len(apiResp.Choices) == 0 || apiResp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("received an empty or invalid response from azure")
	}

	return apiResp.Choices[0].Message.Content, nil
}
