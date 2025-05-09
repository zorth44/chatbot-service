package openrouter

// Response represents the main response structure from OpenRouter API
type Response struct {
	ID                string   `json:"id"`
	Choices           []Choice `json:"choices"`
	Created           int64    `json:"created"` // Unix timestamp
	Model             string   `json:"model"`
	Object            string   `json:"object"` // 'chat.completion' or 'chat.completion.chunk'
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
	Usage             *Usage   `json:"usage,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Choice represents a union type of different choice types
type Choice struct {
	FinishReason       string           `json:"finish_reason"`
	NativeFinishReason string           `json:"native_finish_reason,omitempty"`
	Text               string           `json:"text,omitempty"`
	Message            *ResponseMessage `json:"message,omitempty"`
	Delta              *Delta           `json:"delta,omitempty"`
	Error              *Error           `json:"error,omitempty"`
}

// ResponseMessage represents a non-streaming message response
type ResponseMessage struct {
	Content   string     `json:"content"`
	Role      string     `json:"role"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Delta represents a streaming message delta
type Delta struct {
	Content   string     `json:"content"`
	Role      string     `json:"role,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Error represents an error response
type Error struct {
	Code     int                    `json:"code"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ToolCall represents a function call made by the model
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"` // Always "function"
	Function FunctionCall `json:"function"`
}

// FunctionCall represents the actual function call details
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string of arguments
}
