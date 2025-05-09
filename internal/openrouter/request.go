package openrouter

// Request represents the main request structure for OpenRouter API
type Request struct {
	// Either Messages or Prompt is required
	Messages []Message `json:"messages,omitempty"`
	Prompt   string    `json:"prompt,omitempty"`

	// Optional model specification
	Model string `json:"model,omitempty"`

	// Response format configuration
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`

	Stop   interface{} `json:"stop,omitempty"` // Can be string or []string
	Stream bool        `json:"stream,omitempty"`

	// LLM Parameters
	MaxTokens         int             `json:"max_tokens,omitempty"`
	Temperature       float64         `json:"temperature,omitempty"`
	Tools             []Tool          `json:"tools,omitempty"`
	ToolChoice        interface{}     `json:"tool_choice,omitempty"` // Can be string or ToolChoiceObject
	Seed              int             `json:"seed,omitempty"`
	TopP              float64         `json:"top_p,omitempty"`
	TopK              int             `json:"top_k,omitempty"`
	FrequencyPenalty  float64         `json:"frequency_penalty,omitempty"`
	PresencePenalty   float64         `json:"presence_penalty,omitempty"`
	RepetitionPenalty float64         `json:"repetition_penalty,omitempty"`
	LogitBias         map[int]float64 `json:"logit_bias,omitempty"`
	TopLogprobs       int             `json:"top_logprobs,omitempty"`
	MinP              float64         `json:"min_p,omitempty"`
	TopA              float64         `json:"top_a,omitempty"`

	// Prediction for latency optimization
	Prediction *Prediction `json:"prediction,omitempty"`

	// OpenRouter-specific parameters
	Transforms []string             `json:"transforms,omitempty"`
	Models     []string             `json:"models,omitempty"`
	Route      string               `json:"route,omitempty"`
	Provider   *ProviderPreferences `json:"provider,omitempty"`
}

// ResponseFormat represents the response format configuration
type ResponseFormat struct {
	Type string `json:"type"`
}

// Message represents a message in the conversation
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // Can be string or []ContentPart
	Name    string      `json:"name,omitempty"`
	// For tool messages
	ToolCallID string `json:"tool_call_id,omitempty"`
}

// ContentPart represents a part of the message content
type ContentPart struct {
	Type     string        `json:"type"`
	Text     string        `json:"text,omitempty"`
	ImageURL *ImageURLPart `json:"image_url,omitempty"`
}

// ImageURLPart represents an image URL in the content
type ImageURLPart struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// Tool represents a tool that can be used by the model
type Tool struct {
	Type     string              `json:"type"`
	Function FunctionDescription `json:"function"`
}

// FunctionDescription describes a function that can be called
type FunctionDescription struct {
	Description string      `json:"description,omitempty"`
	Name        string      `json:"name"`
	Parameters  interface{} `json:"parameters"` // JSON Schema object
}

// ToolChoiceObject represents a specific tool choice
type ToolChoiceObject struct {
	Type     string `json:"type"`
	Function struct {
		Name string `json:"name"`
	} `json:"function"`
}

// Prediction represents a predicted output for latency optimization
type Prediction struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// ProviderPreferences represents provider routing preferences
type ProviderPreferences struct {
	// Add provider-specific preferences here
	// This is a placeholder as the exact structure depends on the provider
}
