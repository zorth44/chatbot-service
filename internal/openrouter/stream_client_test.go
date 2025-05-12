package openrouter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zorth44/chatbot-service/internal/config"
)

func TestStreamClient_StreamChatCompletion(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse []string
		serverStatus   int
		expectError    bool
		expectedChunks int
	}{
		{
			name: "successful streaming",
			serverResponse: []string{
				"data: {\"id\":\"1\",\"choices\":[{\"delta\":{\"content\":\"Hello\"}}]}",
				"data: {\"id\":\"1\",\"choices\":[{\"delta\":{\"content\":\" World\"}}]}",
				"data: [DONE]",
			},
			serverStatus:   http.StatusOK,
			expectError:    false,
			expectedChunks: 2,
		},
		{
			name:           "server error",
			serverResponse: []string{},
			serverStatus:   http.StatusInternalServerError,
			expectError:    true,
			expectedChunks: 0,
		},
		{
			name: "invalid json response",
			serverResponse: []string{
				"data: invalid json",
			},
			serverStatus:   http.StatusOK,
			expectError:    true,
			expectedChunks: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request headers
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Errorf("expected Authorization header 'Bearer test-key', got '%s'", r.Header.Get("Authorization"))
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected Content-Type header 'application/json', got '%s'", r.Header.Get("Content-Type"))
				}

				// Set response status
				w.WriteHeader(tt.serverStatus)

				// Write response chunks
				for _, chunk := range tt.serverResponse {
					w.Write([]byte(chunk + "\n"))
					w.(http.Flusher).Flush()
				}
			}))
			defer server.Close()

			// Create client with test configuration
			cfg := &config.OpenRouterConfig{
				BaseURL: server.URL,
				APIKey:  "test-key",
			}
			client := NewStreamClient(cfg)

			// Create test request
			req := &Request{
				Messages: []Message{
					{
						Role:    "user",
						Content: "Hello",
					},
				},
			}

			// Track received chunks
			chunks := 0
			var receivedContent strings.Builder

			// Execute streaming request
			err := client.StreamChatCompletion(req, func(response *Response) error {
				chunks++
				for _, choice := range response.Choices {
					if choice.Delta != nil {
						receivedContent.WriteString(choice.Delta.Content)
					}
				}
				return nil
			})

			// Check error expectations
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			// Check chunk count
			if chunks != tt.expectedChunks {
				t.Errorf("expected %d chunks, got %d", tt.expectedChunks, chunks)
			}

			// For successful case, verify content
			if !tt.expectError && tt.expectedChunks > 0 {
				expectedContent := "Hello World"
				if receivedContent.String() != expectedContent {
					t.Errorf("expected content '%s', got '%s'", expectedContent, receivedContent.String())
				}
			}
		})
	}
}

func TestStreamClient_RequestValidation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read and validate request body
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		// Verify streaming is enabled
		if !req.Stream {
			t.Error("expected Stream to be true")
		}

		// Send empty successful response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data: [DONE]\n"))
	}))
	defer server.Close()

	cfg := &config.OpenRouterConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
	}
	client := NewStreamClient(cfg)

	req := &Request{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Test message",
			},
		},
	}

	err := client.StreamChatCompletion(req, func(response *Response) error {
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
