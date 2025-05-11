package openrouter

import (
	"testing"

	"github.com/zorth44/chatbot-service/internal/config"
)

func TestCreateChatCompletion(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("../../../config/config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Create client
	client := NewClient(&cfg.OpenRouter)

	// Create a simple chat request
	req := &Request{
		Model: "google/gemini-2.5-flash-preview",
		Messages: []Message{
			{
				Role:    "user",
				Content: "can you tell me a joke?",
			},
		},
	}

	// Send request
	resp, err := client.CreateChatCompletion(req)
	if err != nil {
		t.Fatalf("Failed to create chat completion: %v", err)
	}

	// Print response for verification
	t.Logf("Response: %+v", resp)

	// Print the actual content of the response
	if len(resp.Choices) > 0 && resp.Choices[0].Message != nil {
		t.Logf("Response Content: %s", resp.Choices[0].Message.Content)
	} else {
		t.Logf("No content in response or unexpected response format")
	}
}
