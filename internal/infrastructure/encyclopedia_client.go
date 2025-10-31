package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"agent-ollama-gin/internal/domain"
	"agent-ollama-gin/pkg/errors"
	"agent-ollama-gin/pkg/logger"
)

// WikipediaClient handles Wikipedia API communication
type WikipediaClient struct {
	httpClient *http.Client
	apiURL     string
	logger     *logger.Logger
}

// BritannicaClient handles Britannica API communication
type BritannicaClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
	logger     *logger.Logger
}

// EncyclopediaClientFactory provides access to both clients
type EncyclopediaClientFactory struct {
	Wikipedia  *WikipediaClient
	Britannica *BritannicaClient
	Logger     *logger.Logger
}

// NewWikipediaClient creates a new Wikipedia client
func NewWikipediaClient(logger *logger.Logger) *WikipediaClient {
	apiURL := os.Getenv("WIKIPEDIA_API_URL")
	if apiURL == "" {
		apiURL = "https://en.wikipedia.org/api/rest_v1"
	}

	return &WikipediaClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiURL: apiURL,
		logger: logger,
	}
}

// NewBritannicaClient creates a new Britannica client
func NewBritannicaClient(logger *logger.Logger) *BritannicaClient {
	apiKey := os.Getenv("BRITANNICA_API_KEY")
	baseURL := os.Getenv("BRITANNICA_API_URL")
	if baseURL == "" {
		baseURL = "https://api.britannica.com"
	}

	return &BritannicaClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:  apiKey,
		baseURL: baseURL,
		logger:  logger,
	}
}

// NewEncyclopediaClientFactory creates a new encyclopedia client factory
func NewEncyclopediaClientFactory(logger *logger.Logger) *EncyclopediaClientFactory {
	return &EncyclopediaClientFactory{
		Wikipedia:  NewWikipediaClient(logger),
		Britannica: NewBritannicaClient(logger),
		Logger:     logger,
	}
}

// SearchParallel searches both Wikipedia and Britannica concurrently
func (f *EncyclopediaClientFactory) SearchParallel(
	ctx context.Context,
	query, language string,
	maxResults int,
) ([]domain.EncyclopediaSearchResult, error) {
	// Create channels for concurrent results
	wikiResultsChan := make(chan []domain.EncyclopediaSearchResult, 1)
	britResultsChan := make(chan []domain.EncyclopediaSearchResult, 1)
	errChan := make(chan error, 2)

	var wg sync.WaitGroup

	// Search Wikipedia concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		results, err := f.Wikipedia.Search(ctx, query, language, maxResults/2)
		if err != nil {
			f.Logger.Warn("wikipedia search failed", err)
			errChan <- err
			wikiResultsChan <- []domain.EncyclopediaSearchResult{}
			return
		}
		wikiResultsChan <- results
	}()

	// Search Britannica concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		results, err := f.Britannica.Search(ctx, query, language, maxResults/2)
		if err != nil {
			f.Logger.Warn("britannica search failed", err)
			errChan <- err
			britResultsChan <- []domain.EncyclopediaSearchResult{}
			return
		}
		britResultsChan <- results
	}()

	// Wait for all goroutines to complete
	wg.Wait()
	close(wikiResultsChan)
	close(britResultsChan)
	close(errChan)

	// Collect results
	var allResults []domain.EncyclopediaSearchResult
	allResults = append(allResults, <-wikiResultsChan...)
	allResults = append(allResults, <-britResultsChan...)

	return allResults, nil
}

// Wikipedia client methods

// Search performs a Wikipedia search
func (c *WikipediaClient) Search(ctx context.Context, query, language string, maxResults int) ([]domain.EncyclopediaSearchResult, error) {
	searchURL := fmt.Sprintf("%s/search/page?q=%s&limit=%d",
		c.apiURL, url.QueryEscape(query), maxResults)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, errors.Wrap(errors.CodeEncyclopedia, "failed to create wikipedia request", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Encyclopedia-Agent/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(errors.CodeEncyclopedia, "failed to search wikipedia", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(errors.CodeEncyclopedia, "wikipedia API error").WithDetails(map[string]interface{}{
			"status": resp.StatusCode,
			"body":   string(body),
		})
	}

	var searchResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, errors.Wrap(errors.CodeEncyclopedia, "failed to parse wikipedia response", err)
	}

	return c.parseSearchResults(searchResp, language)
}

// GetArticle retrieves a Wikipedia article
func (c *WikipediaClient) GetArticle(ctx context.Context, title string, maxLength int) (*domain.EncyclopediaArticle, error) {
	contentURL := fmt.Sprintf("%s/page/summary/%s", c.apiURL, url.QueryEscape(title))

	req, err := http.NewRequestWithContext(ctx, "GET", contentURL, nil)
	if err != nil {
		return nil, errors.Wrap(errors.CodeEncyclopedia, "failed to create request", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Encyclopedia-Agent/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(errors.CodeEncyclopedia, "failed to get wikipedia article", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(errors.CodeEncyclopedia, "failed to get article").WithDetails(map[string]interface{}{
			"status": resp.StatusCode,
			"title":  title,
		})
	}

	var pageResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&pageResp); err != nil {
		return nil, errors.Wrap(errors.CodeEncyclopedia, "failed to parse article response", err)
	}

	return c.parseArticle(pageResp, maxLength)
}

// Britannica client methods

// Search performs a Britannica search
func (c *BritannicaClient) Search(ctx context.Context, query, language string, maxResults int) ([]domain.EncyclopediaSearchResult, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s&limit=%d&language=%s",
		c.baseURL, url.QueryEscape(query), maxResults, language)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, errors.Wrap(errors.CodeEncyclopedia, "failed to create britannica request", err)
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Encyclopedia-Agent/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(errors.CodeEncyclopedia, "failed to search britannica", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(errors.CodeEncyclopedia, "britannica API error").WithDetails(map[string]interface{}{
			"status": resp.StatusCode,
			"body":   string(body),
		})
	}

	var searchResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, errors.Wrap(errors.CodeEncyclopedia, "failed to parse britannica response", err)
	}

	return c.parseSearchResults(searchResp, language)
}

// GetArticle retrieves a Britannica article
func (c *BritannicaClient) GetArticle(ctx context.Context, title string, maxLength int) (*domain.EncyclopediaArticle, error) {
	contentURL := fmt.Sprintf("%s/article/%s", c.baseURL, url.QueryEscape(title))

	req, err := http.NewRequestWithContext(ctx, "GET", contentURL, nil)
	if err != nil {
		return nil, errors.Wrap(errors.CodeEncyclopedia, "failed to create request", err)
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Encyclopedia-Agent/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(errors.CodeEncyclopedia, "failed to get article", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(errors.CodeEncyclopedia, "britannica article not found").WithDetails(map[string]interface{}{
			"status": resp.StatusCode,
			"body":   string(body),
		})
	}

	var articleResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&articleResp); err != nil {
		return nil, errors.Wrap(errors.CodeEncyclopedia, "failed to parse article response", err)
	}

	return c.parseArticleContent(articleResp, title, maxLength)
}

// Helper methods

func (c *WikipediaClient) parseSearchResults(resp map[string]interface{}, language string) ([]domain.EncyclopediaSearchResult, error) {
	var results []domain.EncyclopediaSearchResult

	pagesData, err := errors.SafeSliceAssert(resp["pages"], "pages")
	if err != nil {
		return results, nil // Return empty results if no pages
	}

	for _, pageData := range pagesData {
		pageMap, err := errors.SafeMapAssert(pageData, "page")
		if err != nil {
			c.logger.Warn("skipping invalid page entry", err)
			continue
		}

		title, err := errors.SafeStringAssert(pageMap["title"], "title")
		if err != nil {
			c.logger.Warn("skipping page without title", err)
			continue
		}

		contentURLs, err := errors.SafeMapAssert(pageMap["content_urls"], "content_urls")
		if err != nil {
			c.logger.Warn("skipping page without content_urls", err)
			continue
		}

		desktop, err := errors.SafeMapAssert(contentURLs["desktop"], "desktop")
		if err != nil {
			c.logger.Warn("skipping page without desktop url", err)
			continue
		}

		pageURL, err := errors.SafeStringAssert(desktop["page"], "page")
		if err != nil {
			c.logger.Warn("skipping page without page url", err)
			continue
		}

		result := domain.EncyclopediaSearchResult{
			Title:     title,
			URL:       pageURL,
			Source:    "wikipedia",
			Language:  language,
			Relevance: 0.9,
		}
		results = append(results, result)
	}

	return results, nil
}

func (c *WikipediaClient) parseArticle(resp map[string]interface{}, maxLength int) (*domain.EncyclopediaArticle, error) {
	title, _ := errors.SafeStringAssert(resp["title"], "title")
	extract, _ := errors.SafeStringAssert(resp["extract"], "extract")

	// Truncate if needed
	if len(extract) > maxLength {
		extract = extract[:maxLength] + "..."
	}

	wordCount := len(strings.Fields(extract))

	contentURLs, _ := errors.SafeMapAssert(resp["content_urls"], "content_urls")
	desktop, _ := errors.SafeMapAssert(contentURLs["desktop"], "desktop")
	fullURL, _ := errors.SafeStringAssert(desktop["page"], "page")

	return &domain.EncyclopediaArticle{
		Title:       title,
		URL:         fullURL,
		Source:      "wikipedia",
		Language:    "en",
		Content:     extract,
		Summary:     extract,
		WordCount:   wordCount,
		LastUpdated: time.Now().Format("2006-01-02"),
	}, nil
}

func (c *BritannicaClient) parseSearchResults(resp map[string]interface{}, language string) ([]domain.EncyclopediaSearchResult, error) {
	var results []domain.EncyclopediaSearchResult

	articlesData, err := errors.SafeSliceAssert(resp["articles"], "articles")
	if err != nil {
		return results, nil
	}

	for _, articleData := range articlesData {
		articleMap, err := errors.SafeMapAssert(articleData, "article")
		if err != nil {
			c.logger.Warn("skipping invalid article entry", err)
			continue
		}

		title, err := errors.SafeStringAssert(articleMap["title"], "title")
		if err != nil {
			c.logger.Warn("skipping article without title", err)
			continue
		}

		articleURL, err := errors.SafeStringAssert(articleMap["url"], "url")
		if err != nil {
			c.logger.Warn("skipping article without url", err)
			continue
		}

		result := domain.EncyclopediaSearchResult{
			Title:     title,
			URL:       articleURL,
			Source:    "britannica",
			Language:  language,
			Relevance: 0.9,
		}
		results = append(results, result)
	}

	return results, nil
}

func (c *BritannicaClient) parseArticleContent(resp map[string]interface{}, title string, maxLength int) (*domain.EncyclopediaArticle, error) {
	content, _ := errors.SafeStringAssert(resp["content"], "content")
	if content == "" {
		content, _ = errors.SafeStringAssert(resp["text"], "text")
	}

	if len(content) > maxLength {
		content = content[:maxLength] + "..."
	}

	summary, _ := errors.SafeStringAssert(resp["summary"], "summary")
	if summary == "" && len(content) > 200 {
		summary = content[:200] + "..."
	}

	wordCount := len(strings.Fields(content))

	var categories []string
	if catsData, err := errors.SafeSliceAssert(resp["categories"], "categories"); err == nil {
		for _, cat := range catsData {
			if catStr, err := errors.SafeStringAssert(cat, "category"); err == nil {
				categories = append(categories, catStr)
			}
		}
	}

	return &domain.EncyclopediaArticle{
		Title:       title,
		Source:      "britannica",
		Language:    "en",
		Content:     content,
		Summary:     summary,
		Categories:  categories,
		WordCount:   wordCount,
		LastUpdated: time.Now().Format("2006-01-02"),
	}, nil
}
