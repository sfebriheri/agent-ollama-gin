package domain

// LLMRequest represents a request to the LLM service
type LLMRequest struct {
	Messages    []Message
	Model       string
	Temperature float64
	MaxTokens   int
	Stream      bool
}

// Message represents a chat message
type Message struct {
	Role    string
	Content string
}

// LLMResponse represents a response from the LLM service
type LLMResponse struct {
	ID      string
	Object  string
	Created int64
	Model   string
	Choices []Choice
	Usage   Usage
	Error   error
}

// Choice represents a completion choice
type Choice struct {
	Index   int
	Message Message
	Delta   Message
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// CompletionRequest represents a text completion request
type CompletionRequest struct {
	Prompt      string
	Model       string
	Temperature float64
	MaxTokens   int
	Stop        string
}

// CompletionResponse represents a text completion response
type CompletionResponse struct {
	ID      string
	Object  string
	Created int64
	Model   string
	Choices []Choice
	Usage   Usage
}

// EmbeddingRequest represents an embedding request
type EmbeddingRequest struct {
	Input string
	Model string
}

// EmbeddingResponse represents an embedding response
type EmbeddingResponse struct {
	Object string
	Data   []Embedding
	Model  string
	Usage  Usage
}

// Embedding represents a text embedding
type Embedding struct {
	Object    string
	Embedding []float64
	Index     int
}

// LLMModel represents a Llama model
type LLMModel struct {
	ID      string
	Object  string
	Created int64
	OwnedBy string
}

// EncyclopediaSearchRequest represents a request to search for encyclopedia articles
type EncyclopediaSearchRequest struct {
	Query          string
	Source         string // "wikipedia", "britannica", "all"
	MaxResults     int
	Language       string
	IncludeSnippet bool
}

// EncyclopediaSearchResult represents a single search result
type EncyclopediaSearchResult struct {
	Title     string
	URL       string
	Snippet   string
	Source    string
	Language  string
	Relevance float64
}

// EncyclopediaSearchResponse represents a search response
type EncyclopediaSearchResponse struct {
	Query      string
	Results    []EncyclopediaSearchResult
	TotalFound int
	Source     string
	Language   string
}

// EncyclopediaArticleRequest represents a request to get a specific article
type EncyclopediaArticleRequest struct {
	Title         string
	URL           string
	Source        string // "wikipedia", "britannica"
	Language      string
	IncludeImages bool
	MaxLength     int
}

// EncyclopediaArticle represents an encyclopedia article
type EncyclopediaArticle struct {
	Title       string
	URL         string
	Source      string
	Language    string
	Content     string
	Summary     string
	Categories  []string
	References  []string
	LastUpdated string
	WordCount   int
}

// EncyclopediaArticleResponse represents an article response
type EncyclopediaArticleResponse struct {
	Article  EncyclopediaArticle
	Related  []string
	Source   string
	Language string
}

// EncyclopediaPromptRequest represents a request to generate encyclopedia-style prompts
type EncyclopediaPromptRequest struct {
	Topic           string
	Style           string // "academic", "casual", "educational"
	Length          string // "short", "medium", "long"
	IncludeExamples bool
	Language        string
}

// EncyclopediaPromptResponse represents a generated encyclopedia prompt
type EncyclopediaPromptResponse struct {
	Topic       string
	Prompt      string
	Style       string
	Length      string
	Language    string
	Suggestions []string
	Keywords    []string
}
