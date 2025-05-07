package openrouter

// TextContent defines a text part of a message.
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
