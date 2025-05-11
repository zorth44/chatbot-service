package openrouter

import (
	"testing"

	"github.com/zorth44/chatbot-service/internal/config"
)

func TestCreateChatCompletion(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("../../config/config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Create client
	client := NewClient(&cfg.OpenRouter)

	// Create a simple chat request
	req := &Request{
		Model: "anthropic/claude-3-opus-20240229",
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello, how are you?",
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
}
