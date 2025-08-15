package models

import "time"

// Message represents a chat message
type Message struct {
	Role    string `json:"role" binding:"required"` // "system", "user", "assistant"
	Content string `json:"content" binding:"required"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Messages    []Message `json:"messages" binding:"required"`
	Model       string    `json:"model,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice
type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
	Delta   Message `json:"delta,omitempty"` // For streaming
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// CompletionRequest represents a text completion request
type CompletionRequest struct {
	Prompt      string  `json:"prompt" binding:"required"`
	Model       string  `json:"model,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Stop        string  `json:"stop,omitempty"`
}

// CompletionResponse represents a text completion response
type CompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// EmbeddingRequest represents an embedding request
type EmbeddingRequest struct {
	Input string `json:"input" binding:"required"`
	Model string `json:"model,omitempty"`
}

// EmbeddingResponse represents an embedding response
type EmbeddingResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
}

// Embedding represents a text embedding
type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// Model represents a Llama model
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
	Code    int    `json:"code"`
}

// StreamResponse represents a streaming response chunk
type StreamResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

// EncyclopediaSearchRequest represents a request to search for encyclopedia articles
type EncyclopediaSearchRequest struct {
	Query          string `json:"query" binding:"required"`
	Source         string `json:"source,omitempty"`      // "wikipedia", "britannica", "all"
	MaxResults     int    `json:"max_results,omitempty"` // Default: 5
	Language       string `json:"language,omitempty"`    // Default: "en"
	IncludeSnippet bool   `json:"include_snippet,omitempty"`
}

// EncyclopediaSearchResult represents a single search result
type EncyclopediaSearchResult struct {
	Title     string  `json:"title"`
	URL       string  `json:"url"`
	Snippet   string  `json:"snippet,omitempty"`
	Source    string  `json:"source"`
	Language  string  `json:"language"`
	Relevance float64 `json:"relevance"`
}

// EncyclopediaSearchResponse represents a search response
type EncyclopediaSearchResponse struct {
	Query      string                     `json:"query"`
	Results    []EncyclopediaSearchResult `json:"results"`
	TotalFound int                        `json:"total_found"`
	Source     string                     `json:"source"`
	Language   string                     `json:"language"`
}

// EncyclopediaArticleRequest represents a request to get a specific article
type EncyclopediaArticleRequest struct {
	Title         string `json:"title,omitempty"`
	URL           string `json:"url,omitempty"`
	Source        string `json:"source,omitempty"`   // "wikipedia", "britannica"
	Language      string `json:"language,omitempty"` // Default: "en"
	IncludeImages bool   `json:"include_images,omitempty"`
	MaxLength     int    `json:"max_length,omitempty"` // Default: 2000 characters
}

// EncyclopediaArticle represents an encyclopedia article
type EncyclopediaArticle struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Source      string   `json:"source"`
	Language    string   `json:"language"`
	Content     string   `json:"content"`
	Summary     string   `json:"summary"`
	Categories  []string `json:"categories,omitempty"`
	References  []string `json:"references,omitempty"`
	LastUpdated string   `json:"last_updated,omitempty"`
	WordCount   int      `json:"word_count"`
}

// EncyclopediaArticleResponse represents an article response
type EncyclopediaArticleResponse struct {
	Article  EncyclopediaArticle `json:"article"`
	Related  []string            `json:"related_articles,omitempty"`
	Source   string              `json:"source"`
	Language string              `json:"language"`
}

// EncyclopediaPromptRequest represents a request to generate encyclopedia-style prompts
type EncyclopediaPromptRequest struct {
	Topic           string `json:"topic" binding:"required"`
	Style           string `json:"style,omitempty"`  // "academic", "casual", "educational"
	Length          string `json:"length,omitempty"` // "short", "medium", "long"
	IncludeExamples bool   `json:"include_examples,omitempty"`
	Language        string `json:"language,omitempty"` // Default: "en"
}

// EncyclopediaPromptResponse represents a generated encyclopedia prompt
type EncyclopediaPromptResponse struct {
	Topic       string   `json:"topic"`
	Prompt      string   `json:"prompt"`
	Style       string   `json:"style"`
	Length      string   `json:"length"`
	Language    string   `json:"language"`
	Suggestions []string `json:"suggestions,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
}
