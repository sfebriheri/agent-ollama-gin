package container

import (
	"llama-api/internal/domain"
	"llama-api/internal/infrastructure"
	"llama-api/internal/usecase"
	"llama-api/pkg/cache"
	"llama-api/pkg/logger"
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
