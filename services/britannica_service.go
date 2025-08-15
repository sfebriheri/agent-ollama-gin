package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"llama-api/models"
)

// BritannicaService handles Britannica Encyclopedia API interactions
type BritannicaService struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
	timeout    time.Duration
}

// NewBritannicaService creates a new Britannica service instance
func NewBritannicaService() *BritannicaService {
	apiKey := os.Getenv("BRITANNICA_API_KEY")
	baseURL := os.Getenv("BRITANNICA_API_URL")
	if baseURL == "" {
		baseURL = "https://api.britannica.com"
	}

	timeout := 30 * time.Second
	if timeoutStr := os.Getenv("BRITANNICA_TIMEOUT"); timeoutStr != "" {
		if timeoutDuration, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = timeoutDuration
		}
	}

	return &BritannicaService{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		apiKey:  apiKey,
		baseURL: baseURL,
		timeout: timeout,
	}
}

// SearchArticles searches for Britannica articles
func (s *BritannicaService) SearchArticles(query, language string, maxResults int) ([]models.EncyclopediaSearchResult, error) {
	// Britannica search endpoint
	searchURL := fmt.Sprintf("%s/search?q=%s&limit=%d&language=%s",
		s.baseURL, url.QueryEscape(query), maxResults, language)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key if available
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Encyclopedia-Agent/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("britannica API error: %d - %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse Britannica search response
	var searchResp map[string]interface{}
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return s.parseSearchResults(searchResp, query, language)
}

// GetArticle retrieves a specific Britannica article
func (s *BritannicaService) GetArticle(title, articleURL, language string, maxLength int) (*models.EncyclopediaArticle, error) {
	var articleID string

	if articleURL != "" {
		// Extract article ID from URL
		articleID = s.extractArticleID(articleURL)
	} else if title != "" {
		// Search for article by title first
		results, err := s.SearchArticles(title, language, 1)
		if err != nil || len(results) == 0 {
			return nil, fmt.Errorf("article not found: %s", title)
		}
		articleID = s.extractArticleID(results[0].URL)
	} else {
		return nil, fmt.Errorf("either title or URL must be provided")
	}

	// Get article content
	contentURL := fmt.Sprintf("%s/article/%s?language=%s", s.baseURL, articleID, language)

	req, err := http.NewRequest("GET", contentURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Encyclopedia-Agent/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("britannica API error: %d - %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var articleResp map[string]interface{}
	if err := json.Unmarshal(body, &articleResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return s.parseArticleContent(articleResp, title, articleURL, language, maxLength)
}

// GetArticleByID retrieves an article directly by its ID
func (s *BritannicaService) GetArticleByID(articleID, language string, maxLength int) (*models.EncyclopediaArticle, error) {
	contentURL := fmt.Sprintf("%s/article/%s?language=%s", s.baseURL, articleID, language)

	req, err := http.NewRequest("GET", contentURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Encyclopedia-Agent/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("britannica API error: %d - %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var articleResp map[string]interface{}
	if err := json.Unmarshal(body, &articleResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return s.parseArticleContent(articleResp, "", "", language, maxLength)
}

// Helper methods
func (s *BritannicaService) parseSearchResults(resp map[string]interface{}, query, language string) ([]models.EncyclopediaSearchResult, error) {
	var results []models.EncyclopediaSearchResult

	// Parse Britannica search response structure
	// This will need to be adjusted based on the actual API response format
	if articles, ok := resp["articles"].([]interface{}); ok {
		for _, article := range articles {
			if articleData, ok := article.(map[string]interface{}); ok {
				title, _ := articleData["title"].(string)
				url, _ := articleData["url"].(string)
				snippet, _ := articleData["snippet"].(string)

				result := models.EncyclopediaSearchResult{
					Title:     title,
					URL:       url,
					Snippet:   snippet,
					Source:    "britannica",
					Language:  language,
					Relevance: 0.9, // Default relevance score
				}
				results = append(results, result)
			}
		}
	}

	// If no articles found, return a fallback result
	if len(results) == 0 {
		fallbackURL := fmt.Sprintf("https://www.britannica.com/topic/%s", url.QueryEscape(query))
		results = append(results, models.EncyclopediaSearchResult{
			Title:     query + " - Britannica",
			URL:       fallbackURL,
			Source:    "britannica",
			Language:  language,
			Relevance: 0.8,
		})
	}

	return results, nil
}

func (s *BritannicaService) parseArticleContent(resp map[string]interface{}, title, articleURL, language string, maxLength int) (*models.EncyclopediaArticle, error) {
	// Parse Britannica article response structure
	// This will need to be adjusted based on the actual API response format

	content, _ := resp["content"].(string)
	if content == "" {
		content = resp["text"].(string) // Alternative field name
	}

	// Truncate content if needed
	if len(content) > maxLength {
		content = content[:maxLength] + "..."
	}

	// Extract additional metadata
	summary, _ := resp["summary"].(string)
	if summary == "" {
		summary = content[:min(len(content), 200)] + "..."
	}

	// Count words
	wordCount := len(strings.Fields(content))

	// Extract categories if available
	var categories []string
	if cats, ok := resp["categories"].([]interface{}); ok {
		for _, cat := range cats {
			if catStr, ok := cat.(string); ok {
				categories = append(categories, catStr)
			}
		}
	}

	// Extract references if available
	var references []string
	if refs, ok := resp["references"].([]interface{}); ok {
		for _, ref := range refs {
			if refStr, ok := ref.(string); ok {
				references = append(references, refStr)
			}
		}
	}

	return &models.EncyclopediaArticle{
		Title:       title,
		URL:         articleURL,
		Source:      "britannica",
		Language:    language,
		Content:     content,
		Summary:     summary,
		Categories:  categories,
		References:  references,
		LastUpdated: time.Now().Format("2006-01-02"),
		WordCount:   wordCount,
	}, nil
}

func (s *BritannicaService) extractArticleID(urlStr string) string {
	// Extract article ID from Britannica URL
	// Example: https://www.britannica.com/topic/artificial-intelligence -> artificial-intelligence
	if strings.Contains(urlStr, "britannica.com/topic/") {
		parts := strings.Split(urlStr, "/topic/")
		if len(parts) > 1 {
			return strings.TrimSuffix(parts[1], "/")
		}
	}
	return urlStr
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Health check for Britannica service
func (s *BritannicaService) Health() error {
	// Test API connectivity
	testURL := fmt.Sprintf("%s/health", s.baseURL)
	resp, err := s.httpClient.Get(testURL)
	if err != nil {
		return fmt.Errorf("britannica API not accessible: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("britannica API health check failed: %d", resp.StatusCode)
	}

	return nil
}
