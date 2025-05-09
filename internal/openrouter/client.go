package openrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL = "https://openrouter.ai/api/v1"
)

// Client represents an OpenRouter API client
type Client struct {
	apiKey     string
	siteURL    string
	siteName   string
	httpClient *http.Client
}

// NewClient creates a new OpenRouter API client
func NewClient(apiKey, siteURL, siteName string) *Client {
	return &Client{
		apiKey:     apiKey,
		siteURL:    siteURL,
		siteName:   siteName,
		httpClient: &http.Client{},
	}
}

// CreateChatCompletion sends a chat completion request to the OpenRouter API
func (c *Client) CreateChatCompletion(req *Request) (*Response, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")
	if c.siteURL != "" {
		httpReq.Header.Set("HTTP-Referer", c.siteURL)
	}
	if c.siteName != "" {
		httpReq.Header.Set("X-Title", c.siteName)
	}

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
