package container

import (
	"agent-ollama-gin/internal/domain"
	"agent-ollama-gin/internal/infrastructure"
	"agent-ollama-gin/internal/usecase"
	"agent-ollama-gin/pkg/cache"
	"agent-ollama-gin/pkg/logger"
)

// Container holds all application dependencies
type Container struct {
	Logger              *logger.Logger
	Cache               *cache.MemoryCache
	OllamaClient        *infrastructure.OllamaClient
	EncyclopediaFactory *infrastructure.EncyclopediaClientFactory
	LLMUsecase          domain.LLMUsecase
	EncyclopediaUsecase domain.EncyclopediaUsecase
}

// New creates and initializes a new dependency container
func New() *Container {
	// Initialize logger
	appLogger := logger.New("App")

	// Initialize cache
	memCache := cache.NewMemoryCache()

	// Initialize infrastructure clients
	ollamaClient := infrastructure.NewOllamaClient(appLogger)
	encyclopediaFactory := infrastructure.NewEncyclopediaClientFactory(appLogger)

	// Initialize use cases
	llmUsecase := usecase.NewLLMUsecase(ollamaClient, memCache, appLogger)
	encyclopediaUsecase := usecase.NewEncyclopediaUsecase(encyclopediaFactory, memCache, appLogger)

	return &Container{
		Logger:              appLogger,
		Cache:               memCache,
		OllamaClient:        ollamaClient,
		EncyclopediaFactory: encyclopediaFactory,
		LLMUsecase:          llmUsecase,
		EncyclopediaUsecase: encyclopediaUsecase,
	}
}
