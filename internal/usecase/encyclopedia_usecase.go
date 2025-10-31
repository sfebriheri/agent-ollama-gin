package usecase

import (
	"context"
	"fmt"
	"strings"

	"agent-ollama-gin/internal/domain"
	"agent-ollama-gin/internal/infrastructure"
	"agent-ollama-gin/pkg/logger"
)

// EncyclopediaUsecase implements domain.EncyclopediaUsecase interface
type EncyclopediaUsecase struct {
	factory *infrastructure.EncyclopediaClientFactory
	cache   domain.Cache
	logger  *logger.Logger
}

// NewEncyclopediaUsecase creates a new encyclopedia use case
func NewEncyclopediaUsecase(
	factory *infrastructure.EncyclopediaClientFactory,
	cache domain.Cache,
	logger *logger.Logger,
) domain.EncyclopediaUsecase {
	return &EncyclopediaUsecase{
		factory: factory,
		cache:   cache,
		logger:  logger,
	}
}

// SearchEncyclopedia performs a parallel search across multiple sources
func (u *EncyclopediaUsecase) SearchEncyclopedia(
	ctx context.Context,
	request domain.EncyclopediaSearchRequest,
) (*domain.EncyclopediaSearchResponse, error) {
	// Set defaults
	if request.MaxResults == 0 {
		request.MaxResults = 5
	}
	if request.Language == "" {
		request.Language = "en"
	}
	if request.Source == "" {
		request.Source = "all"
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("search:%s:%s:%d", request.Query, request.Language, request.MaxResults)
	if cached, err := u.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		if response, ok := cached.(*domain.EncyclopediaSearchResponse); ok {
			u.logger.Debug("cache hit for search", cacheKey)
			return response, nil
		}
	}

	var results []domain.EncyclopediaSearchResult
	var totalFound int

	switch request.Source {
	case "wikipedia":
		wikiResults, err := u.factory.Wikipedia.Search(ctx, request.Query, request.Language, request.MaxResults)
		if err != nil {
			u.logger.Error("wikipedia search failed", err)
			return nil, err
		}
		results = wikiResults
		totalFound = len(wikiResults)

	case "britannica":
		britResults, err := u.factory.Britannica.Search(ctx, request.Query, request.Language, request.MaxResults)
		if err != nil {
			u.logger.Error("britannica search failed", err)
			return nil, err
		}
		results = britResults
		totalFound = len(britResults)

	case "all":
		// Use parallel search for both sources
		parallelResults, err := u.factory.SearchParallel(ctx, request.Query, request.Language, request.MaxResults)
		if err != nil {
			u.logger.Warn("parallel search partially failed", err)
		}
		results = parallelResults
		totalFound = len(parallelResults)
	}

	response := &domain.EncyclopediaSearchResponse{
		Query:      request.Query,
		Results:    results,
		TotalFound: totalFound,
		Source:     request.Source,
		Language:   request.Language,
	}

	// Cache the response for 1 hour
	_ = u.cache.Set(ctx, cacheKey, response, 3600)

	return response, nil
}

// GetArticle retrieves a specific encyclopedia article
func (u *EncyclopediaUsecase) GetArticle(
	ctx context.Context,
	request domain.EncyclopediaArticleRequest,
) (*domain.EncyclopediaArticleResponse, error) {
	// Set defaults
	if request.Language == "" {
		request.Language = "en"
	}
	if request.MaxLength == 0 {
		request.MaxLength = 2000
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("article:%s:%s", request.Title+request.URL, request.Language)
	if cached, err := u.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		if response, ok := cached.(*domain.EncyclopediaArticleResponse); ok {
			u.logger.Debug("cache hit for article", cacheKey)
			return response, nil
		}
	}

	var article *domain.EncyclopediaArticle
	var source string

	// Determine source
	if request.URL != "" {
		if strings.Contains(request.URL, "wikipedia.org") {
			source = "wikipedia"
		} else if strings.Contains(request.URL, "britannica.com") {
			source = "britannica"
		}
	} else {
		source = request.Source
	}

	switch source {
	case "wikipedia":
		wikiArticle, err := u.factory.Wikipedia.GetArticle(ctx, request.Title, request.MaxLength)
		if err != nil {
			return nil, err
		}
		article = wikiArticle

	case "britannica":
		britArticle, err := u.factory.Britannica.GetArticle(ctx, request.Title, request.MaxLength)
		if err != nil {
			return nil, err
		}
		article = britArticle

	default:
		return nil, fmt.Errorf("unsupported source: %s", source)
	}

	response := &domain.EncyclopediaArticleResponse{
		Article:  *article,
		Source:   source,
		Language: request.Language,
	}

	// Cache the response for 24 hours
	_ = u.cache.Set(ctx, cacheKey, response, 86400)

	return response, nil
}

// GeneratePrompt generates encyclopedia-style prompts using the LLM
func (u *EncyclopediaUsecase) GeneratePrompt(
	ctx context.Context,
	request domain.EncyclopediaPromptRequest,
	llmService domain.LLMService,
) (*domain.EncyclopediaPromptResponse, error) {
	// Set defaults
	if request.Language == "" {
		request.Language = "en"
	}
	if request.Style == "" {
		request.Style = "educational"
	}
	if request.Length == "" {
		request.Length = "medium"
	}

	// Create a system prompt for encyclopedia-style content
	systemPrompt := fmt.Sprintf(`You are an expert encyclopedia editor. Generate a comprehensive, well-structured prompt about "%s" in %s style with %s length.

The prompt should:
- Be informative and educational
- Follow encyclopedia writing conventions
- Include key facts and concepts
- Be appropriate for %s language speakers
- Maintain academic rigor while being accessible

Style guidelines for %s:
- Academic: Formal, detailed, with citations
- Casual: Conversational, engaging, easy to understand
- Educational: Clear explanations, examples, learning objectives

Length guidelines for %s:
- Short: 100-200 words, key points only
- Medium: 300-500 words, balanced coverage
- Long: 600-1000 words, comprehensive treatment`,
		request.Topic, request.Style, request.Length, request.Language, request.Style, request.Length)

	// Create the user prompt
	userPrompt := fmt.Sprintf("Generate an encyclopedia-style prompt about: %s", request.Topic)

	// Use the LLM service to generate the prompt
	chatRequest := domain.LLMRequest{
		Messages: []domain.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	response, err := llmService.Chat(ctx, chatRequest)
	if err != nil {
		return nil, err
	}

	// Extract the generated content
	var generatedContent string
	if len(response.Choices) > 0 {
		generatedContent = response.Choices[0].Message.Content
	}

	// Generate suggestions and keywords
	suggestions := u.generateSuggestions(request.Topic, request.Style)
	keywords := u.extractKeywords(generatedContent)

	return &domain.EncyclopediaPromptResponse{
		Topic:       request.Topic,
		Prompt:      generatedContent,
		Style:       request.Style,
		Length:      request.Length,
		Language:    request.Language,
		Suggestions: suggestions,
		Keywords:    keywords,
	}, nil
}

// Helper methods

func (u *EncyclopediaUsecase) generateSuggestions(topic, style string) []string {
	suggestions := []string{
		fmt.Sprintf("History of %s", topic),
		fmt.Sprintf("Modern developments in %s", topic),
		fmt.Sprintf("Key figures in %s", topic),
	}

	switch style {
	case "academic":
		suggestions = append(suggestions,
			fmt.Sprintf("Research methodologies in %s", topic),
			fmt.Sprintf("Theoretical frameworks of %s", topic))
	case "casual":
		suggestions = append(suggestions,
			fmt.Sprintf("Fun facts about %s", topic),
			fmt.Sprintf("Everyday applications of %s", topic))
	}

	return suggestions
}

func (u *EncyclopediaUsecase) extractKeywords(content string) []string {
	words := strings.Fields(strings.ToLower(content))
	keywordMap := make(map[string]int)

	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "is": true, "are": true, "was": true, "were": true,
	}

	for _, word := range words {
		cleanWord := strings.Trim(word, ".,!?;:\"()[]{}")
		if len(cleanWord) > 3 && !stopWords[cleanWord] {
			keywordMap[cleanWord]++
		}
	}

	var keywords []string
	for word, count := range keywordMap {
		if count > 1 {
			keywords = append(keywords, word)
		}
	}

	if len(keywords) > 10 {
		keywords = keywords[:10]
	}

	return keywords
}
