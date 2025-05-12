package openrouter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zorth44/chatbot-service/internal/config"
)

// StreamClient represents an OpenRouter API client for streaming responses
type StreamClient struct {
	config     *config.OpenRouterConfig
	httpClient *http.Client
}

// NewStreamClient creates a new OpenRouter API streaming client
func NewStreamClient(cfg *config.OpenRouterConfig) *StreamClient {
	return &StreamClient{
		config:     cfg,
		httpClient: &http.Client{},
	}
}

// StreamChatCompletion sends a streaming chat completion request to the OpenRouter API
func (c *StreamClient) StreamChatCompletion(req *Request, handler func(*Response) error) error {
	// Ensure streaming is enabled
	req.Stream = true

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", c.config.BaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	httpReq.Header.Set("Content-Type", "application/json")
	if c.config.SiteURL != "" {
		httpReq.Header.Set("HTTP-Referer", c.config.SiteURL)
	}
	if c.config.SiteName != "" {
		httpReq.Header.Set("X-Title", c.config.SiteName)
	}

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for error status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Create a scanner to read the response line by line
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Remove "data: " prefix if present
		if len(line) > 6 && line[:6] == "data: " {
			line = line[6:]
		}

		// Skip "[DONE]" message
		if line == "[DONE]" {
			continue
		}

		// Parse the response
		var response Response
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		// Call the handler with the response
		if err := handler(&response); err != nil {
			return fmt.Errorf("handler error: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	return nil
}
