package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"llama-api/models"
)

type EncyclopediaService struct {
	httpClient        *http.Client
	wikipediaAPI      string
	britannicaService *BritannicaService
}

func NewEncyclopediaService() *EncyclopediaService {
	wikipediaAPI := os.Getenv("WIKIPEDIA_API_URL")
	if wikipediaAPI == "" {
		wikipediaAPI = "https://en.wikipedia.org/api/rest_v1"
	}

	return &EncyclopediaService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		wikipediaAPI:      wikipediaAPI,
		britannicaService: NewBritannicaService(),
	}
}

// SearchEncyclopedia searches for encyclopedia articles
func (s *EncyclopediaService) SearchEncyclopedia(request models.EncyclopediaSearchRequest) (*models.EncyclopediaSearchResponse, error) {
	if request.MaxResults == 0 {
		request.MaxResults = 5
	}
	if request.Language == "" {
		request.Language = "en"
	}
	if request.Source == "" {
		request.Source = "all"
	}

	var results []models.EncyclopediaSearchResult
	var totalFound int

	switch request.Source {
	case "wikipedia":
		wikiResults, err := s.searchWikipedia(request.Query, request.Language, request.MaxResults)
		if err != nil {
			return nil, fmt.Errorf("wikipedia search failed: %w", err)
		}
		results = append(results, wikiResults...)
		totalFound = len(wikiResults)

	case "britannica":
		britResults, err := s.britannicaService.SearchArticles(request.Query, request.Language, request.MaxResults)
		if err != nil {
			return nil, fmt.Errorf("britannica search failed: %w", err)
		}
		results = append(results, britResults...)
		totalFound = len(britResults)

	case "all":
		// Search both sources
		wikiResults, err := s.searchWikipedia(request.Query, request.Language, request.MaxResults/2)
		if err != nil {
			wikiResults = []models.EncyclopediaSearchResult{}
		}

		britResults, err := s.britannicaService.SearchArticles(request.Query, request.Language, request.MaxResults/2)
		if err != nil {
			britResults = []models.EncyclopediaSearchResult{}
		}

		results = append(wikiResults, britResults...)
		totalFound = len(results)
	}

	return &models.EncyclopediaSearchResponse{
		Query:      request.Query,
		Results:    results,
		TotalFound: totalFound,
		Source:     request.Source,
		Language:   request.Language,
	}, nil
}

// GetArticle retrieves a specific encyclopedia article
func (s *EncyclopediaService) GetArticle(request models.EncyclopediaArticleRequest) (*models.EncyclopediaArticleResponse, error) {
	if request.Language == "" {
		request.Language = "en"
	}
	if request.MaxLength == 0 {
		request.MaxLength = 2000
	}

	var article *models.EncyclopediaArticle
	var source string

	if request.URL != "" {
		// Extract source from URL
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
		wikiArticle, err := s.getWikipediaArticle(request.Title, request.URL, request.Language, request.MaxLength)
		if err != nil {
			return nil, fmt.Errorf("failed to get wikipedia article: %w", err)
		}
		article = wikiArticle

	case "britannica":
		britArticle, err := s.britannicaService.GetArticle(request.Title, request.URL, request.Language, request.MaxLength)
		if err != nil {
			return nil, fmt.Errorf("failed to get britannica article: %w", err)
		}
		article = britArticle

	default:
		return nil, fmt.Errorf("unsupported source: %s", source)
	}

	return &models.EncyclopediaArticleResponse{
		Article:  *article,
		Source:   source,
		Language: request.Language,
	}, nil
}

// GenerateEncyclopediaPrompt generates encyclopedia-style prompts using the LLM
func (s *EncyclopediaService) GenerateEncyclopediaPrompt(request models.EncyclopediaPromptRequest, llamaService *LlamaService) (*models.EncyclopediaPromptResponse, error) {
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
	chatRequest := models.ChatRequest{
		Messages: []models.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	response, err := llamaService.Chat(chatRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to generate prompt with LLM: %w", err)
	}

	// Extract the generated content
	generatedContent := response.Choices[0].Message.Content

	// Generate suggestions and keywords
	suggestions := s.generateSuggestions(request.Topic, request.Style)
	keywords := s.extractKeywords(generatedContent)

	return &models.EncyclopediaPromptResponse{
		Topic:       request.Topic,
		Prompt:      generatedContent,
		Style:       request.Style,
		Length:      request.Length,
		Language:    request.Language,
		Suggestions: suggestions,
		Keywords:    keywords,
	}, nil
}

// Helper methods for Wikipedia
func (s *EncyclopediaService) searchWikipedia(query, language string, maxResults int) ([]models.EncyclopediaSearchResult, error) {
	// Wikipedia search API endpoint
	searchURL := fmt.Sprintf("%s/search/page?q=%s&limit=%d&language=%s",
		s.wikipediaAPI, url.QueryEscape(query), maxResults, language)

	resp, err := s.httpClient.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse Wikipedia search response
	var searchResp map[string]interface{}
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, err
	}

	var results []models.EncyclopediaSearchResult
	if pages, ok := searchResp["pages"].([]interface{}); ok {
		for _, page := range pages {
			if pageData, ok := page.(map[string]interface{}); ok {
				title, _ := pageData["title"].(string)
				contentURLs, _ := pageData["content_urls"].(map[string]interface{})
				if webURL, ok := contentURLs["desktop"].(map[string]interface{}); ok {
					if pageURL, ok := webURL["page"].(string); ok {
						result := models.EncyclopediaSearchResult{
							Title:     title,
							URL:       pageURL,
							Source:    "wikipedia",
							Language:  language,
							Relevance: 0.9, // Default relevance score
						}
						results = append(results, result)
					}
				}
			}
		}
	}

	return results, nil
}

func (s *EncyclopediaService) getWikipediaArticle(title, articleURL, language string, maxLength int) (*models.EncyclopediaArticle, error) {
	// Extract title from URL if not provided
	if title == "" && articleURL != "" {
		title = s.extractTitleFromURL(articleURL)
	}

	// Wikipedia content API endpoint
	contentURL := fmt.Sprintf("%s/page/summary/%s", s.wikipediaAPI, url.QueryEscape(title))

	resp, err := s.httpClient.Get(contentURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pageResp map[string]interface{}
	if err := json.Unmarshal(body, &pageResp); err != nil {
		return nil, err
	}

	extract, _ := pageResp["extract"].(string)
	contentURLs, _ := pageResp["content_urls"].(map[string]interface{})

	var fullURL string
	if desktop, ok := contentURLs["desktop"].(map[string]interface{}); ok {
		fullURL, _ = desktop["page"].(string)
	}

	// Truncate content if needed
	if len(extract) > maxLength {
		extract = extract[:maxLength] + "..."
	}

	// Count words
	wordCount := len(strings.Fields(extract))

	return &models.EncyclopediaArticle{
		Title:       title,
		URL:         fullURL,
		Source:      "wikipedia",
		Language:    language,
		Content:     extract,
		Summary:     extract,
		WordCount:   wordCount,
		LastUpdated: time.Now().Format("2006-01-02"),
	}, nil
}

// Helper methods for Britannica - now handled by BritannicaService

// Helper methods
func (s *EncyclopediaService) extractTitleFromURL(urlStr string) string {
	// Extract title from Wikipedia URL
	re := regexp.MustCompile(`/wiki/(.+)$`)
	matches := re.FindStringSubmatch(urlStr)
	if len(matches) > 1 {
		return strings.ReplaceAll(matches[1], "_", " ")
	}
	return ""
}

func (s *EncyclopediaService) generateSuggestions(topic, style string) []string {
	// Generate related topics based on the main topic and style
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

func (s *EncyclopediaService) extractKeywords(content string) []string {
	// Simple keyword extraction (in a real implementation, use NLP libraries)
	words := strings.Fields(strings.ToLower(content))
	keywordMap := make(map[string]int)

	// Common stop words to filter out
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "is": true, "are": true, "was": true, "were": true,
	}

	for _, word := range words {
		// Clean the word
		cleanWord := strings.Trim(word, ".,!?;:\"()[]{}")
		if len(cleanWord) > 3 && !stopWords[cleanWord] {
			keywordMap[cleanWord]++
		}
	}

	// Get top keywords
	var keywords []string
	for word, count := range keywordMap {
		if count > 1 { // Only include words that appear more than once
			keywords = append(keywords, word)
		}
	}

	// Limit to top 10 keywords
	if len(keywords) > 10 {
		keywords = keywords[:10]
	}

	return keywords
}
