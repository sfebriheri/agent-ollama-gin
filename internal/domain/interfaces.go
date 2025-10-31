package domain

import "context"

// LLMService defines the interface for LLM operations
type LLMService interface {
	Chat(ctx context.Context, request LLMRequest) (*LLMResponse, error)
	Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)
	Embedding(ctx context.Context, request EmbeddingRequest) (*EmbeddingResponse, error)
	ListModels(ctx context.Context) ([]LLMModel, error)
	StreamChat(ctx context.Context, request LLMRequest, responseChan chan<- string) error
}

// EncyclopediaSearcher defines the interface for encyclopedia searching
type EncyclopediaSearcher interface {
	SearchWikipedia(ctx context.Context, query, language string, maxResults int) ([]EncyclopediaSearchResult, error)
	SearchBritannica(ctx context.Context, query, language string, maxResults int) ([]EncyclopediaSearchResult, error)
	SearchAll(ctx context.Context, query, language string, maxResults int) ([]EncyclopediaSearchResult, error)
}

// EncyclopediaProvider defines the interface for encyclopedia operations
type EncyclopediaProvider interface {
	GetWikipediaArticle(ctx context.Context, title, url, language string, maxLength int) (*EncyclopediaArticle, error)
	GetBritannicaArticle(ctx context.Context, title, url, language string, maxLength int) (*EncyclopediaArticle, error)
	GetArticle(ctx context.Context, request EncyclopediaArticleRequest) (*EncyclopediaArticleResponse, error)
}

// EncyclopediaUsecase defines the interface for encyclopedia use cases
type EncyclopediaUsecase interface {
	SearchEncyclopedia(ctx context.Context, request EncyclopediaSearchRequest) (*EncyclopediaSearchResponse, error)
	GetArticle(ctx context.Context, request EncyclopediaArticleRequest) (*EncyclopediaArticleResponse, error)
	GeneratePrompt(ctx context.Context, request EncyclopediaPromptRequest, llmService LLMService) (*EncyclopediaPromptResponse, error)
}

// LLMUsecase defines the interface for LLM use cases
type LLMUsecase interface {
	Chat(ctx context.Context, request LLMRequest) (*LLMResponse, error)
	Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)
	Embedding(ctx context.Context, request EmbeddingRequest) (*EmbeddingResponse, error)
	ListModels(ctx context.Context) ([]LLMModel, error)
	StreamChat(ctx context.Context, request LLMRequest, responseChan chan<- string) error
}

// Cache defines the interface for caching operations
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl int64) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

// Logger defines the interface for logging
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

// HTTPClient defines the interface for HTTP operations
type HTTPClient interface {
	Get(ctx context.Context, url string) (interface{}, error)
	Post(ctx context.Context, url string, body interface{}) (interface{}, error)
}
