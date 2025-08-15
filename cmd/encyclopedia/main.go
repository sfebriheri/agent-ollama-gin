package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"llama-api/models"
)

const (
	baseURL = "http://localhost:8080/api/v1"
	version = "1.0.0"
)

type EncyclopediaCLI struct {
	client *http.Client
}

func NewEncyclopediaCLI() *EncyclopediaCLI {
	return &EncyclopediaCLI{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (cli *EncyclopediaCLI) makeRequest(method, endpoint string, data interface{}) ([]byte, error) {
	var req *http.Request
	var err error

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		req, err = http.NewRequest(method, baseURL+endpoint, strings.NewReader(string(jsonData)))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, baseURL+endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
	}

	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (cli *EncyclopediaCLI) search(query, source, language string, maxResults int) error {
	request := models.EncyclopediaSearchRequest{
		Query:      query,
		Source:     source,
		Language:   language,
		MaxResults: maxResults,
	}

	body, err := cli.makeRequest("POST", "/encyclopedia/search", request)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	var response models.EncyclopediaSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Printf("\n🔍 Search Results for: %s\n", query)
	fmt.Printf("📊 Found %d results from %s\n\n", response.TotalFound, response.Source)

	for i, result := range response.Results {
		fmt.Printf("%d. 📖 %s\n", i+1, result.Title)
		fmt.Printf("   🌐 %s\n", result.URL)
		fmt.Printf("   📝 %s\n", result.Snippet)
		fmt.Printf("   🏷️  %s (%s)\n", result.Source, result.Language)
		fmt.Printf("   ⭐ Relevance: %.2f\n\n", result.Relevance)
	}

	return nil
}

func (cli *EncyclopediaCLI) getArticle(title, source, language string, maxLength int) error {
	request := models.EncyclopediaArticleRequest{
		Title:     title,
		Source:    source,
		Language:  language,
		MaxLength: maxLength,
	}

	body, err := cli.makeRequest("POST", "/encyclopedia/article", request)
	if err != nil {
		return fmt.Errorf("article retrieval failed: %w", err)
	}

	var response models.EncyclopediaArticleResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	article := response.Article
	fmt.Printf("\n📖 Article: %s\n", article.Title)
	fmt.Printf("🌐 Source: %s (%s)\n", article.Source, article.Language)
	fmt.Printf("📅 Updated: %s\n", article.LastUpdated)
	fmt.Printf("📊 Word Count: %d\n", article.WordCount)
	fmt.Printf("🏷️  Categories: %s\n", strings.Join(article.Categories, ", "))
	fmt.Printf("🔗 URL: %s\n\n", article.URL)
	fmt.Printf("📝 Summary:\n%s\n\n", article.Summary)
	fmt.Printf("📄 Content:\n%s\n", article.Content)

	if len(response.Related) > 0 {
		fmt.Printf("\n🔗 Related Articles:\n")
		for _, related := range response.Related {
			fmt.Printf("   • %s\n", related)
		}
	}

	return nil
}

func (cli *EncyclopediaCLI) generatePrompt(topic, style, length, language string) error {
	request := models.EncyclopediaPromptRequest{
		Topic:    topic,
		Style:    style,
		Length:   length,
		Language: language,
	}

	body, err := cli.makeRequest("POST", "/encyclopedia/prompt", request)
	if err != nil {
		return fmt.Errorf("prompt generation failed: %w", err)
	}

	var response models.EncyclopediaPromptResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Printf("\n✍️  Generated Prompt for: %s\n", topic)
	fmt.Printf("🎨 Style: %s\n", response.Style)
	fmt.Printf("📏 Length: %s\n", response.Length)
	fmt.Printf("🌍 Language: %s\n\n", response.Language)
	fmt.Printf("📝 Prompt:\n%s\n\n", response.Prompt)

	if len(response.Suggestions) > 0 {
		fmt.Printf("💡 Suggestions:\n")
		for _, suggestion := range response.Suggestions {
			fmt.Printf("   • %s\n", suggestion)
		}
		fmt.Println()
	}

	if len(response.Keywords) > 0 {
		fmt.Printf("🔑 Keywords:\n")
		for _, keyword := range response.Keywords {
			fmt.Printf("   • %s\n", keyword)
		}
		fmt.Println()
	}

	return nil
}

func (cli *EncyclopediaCLI) getSources() error {
	body, err := cli.makeRequest("GET", "/encyclopedia/sources", nil)
	if err != nil {
		return fmt.Errorf("failed to get sources: %w", err)
	}

	var sources map[string]interface{}
	if err := json.Unmarshal(body, &sources); err != nil {
		return fmt.Errorf("failed to parse sources: %w", err)
	}

	fmt.Printf("\n📚 Available Sources:\n")
	for source, details := range sources {
		if detailsMap, ok := details.(map[string]interface{}); ok {
			fmt.Printf("   • %s: %s\n", source, detailsMap["description"])
		} else {
			fmt.Printf("   • %s\n", source)
		}
	}
	fmt.Println()

	return nil
}

func (cli *EncyclopediaCLI) getLanguages() error {
	body, err := cli.makeRequest("GET", "/encyclopedia/languages", nil)
	if err != nil {
		return fmt.Errorf("failed to get languages: %w", err)
	}

	var languages map[string]interface{}
	if err := json.Unmarshal(body, &languages); err != nil {
		return fmt.Errorf("failed to parse languages: %w", err)
	}

	fmt.Printf("\n🌍 Supported Languages:\n")
	for code, name := range languages {
		fmt.Printf("   • %s: %s\n", code, name)
	}
	fmt.Println()

	return nil
}

func (cli *EncyclopediaCLI) health() error {
	body, err := cli.makeRequest("GET", "/encyclopedia/health", nil)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	var health map[string]interface{}
	if err := json.Unmarshal(body, &health); err != nil {
		return fmt.Errorf("failed to parse health response: %w", err)
	}

	fmt.Printf("\n🏥 Encyclopedia Service Health:\n")
	for key, value := range health {
		fmt.Printf("   • %s: %v\n", key, value)
	}
	fmt.Println()

	return nil
}

func (cli *EncyclopediaCLI) openURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}

func (cli *EncyclopediaCLI) interactiveMode() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\n🔍 Encyclopedia Agent CLI v%s\n", version)
	fmt.Printf("Type 'help' for available commands, 'quit' to exit\n\n")

	for {
		fmt.Print("📚 encyclopedia> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := strings.ToLower(parts[0])

		switch command {
		case "quit", "exit":
			fmt.Println("👋 Goodbye!")
			return
		case "help":
			cli.showHelp()
		case "search":
			if len(parts) < 2 {
				fmt.Println("❌ Usage: search <query> [source] [language] [max_results]")
				continue
			}
			query := parts[1]
			source := "all"
			language := "en"
			maxResults := 5

			if len(parts) > 2 {
				source = parts[2]
			}
			if len(parts) > 3 {
				language = parts[3]
			}
			if len(parts) > 4 {
				if max, err := fmt.Sscanf(parts[4], "%d", &maxResults); err != nil || max != 1 {
					fmt.Println("❌ Invalid max_results number")
					continue
				}
			}

			if err := cli.search(query, source, language, maxResults); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

		case "article":
			if len(parts) < 2 {
				fmt.Println("❌ Usage: article <title> [source] [language] [max_length]")
				continue
			}
			title := parts[1]
			source := "wikipedia"
			language := "en"
			maxLength := 2000

			if len(parts) > 2 {
				source = parts[2]
			}
			if len(parts) > 3 {
				language = parts[3]
			}
			if len(parts) > 4 {
				if max, err := fmt.Sscanf(parts[4], "%d", &maxLength); err != nil || max != 1 {
					fmt.Println("❌ Invalid max_length number")
					continue
				}
			}

			if err := cli.getArticle(title, source, language, maxLength); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

		case "prompt":
			if len(parts) < 2 {
				fmt.Println("❌ Usage: prompt <topic> [style] [length] [language]")
				continue
			}
			topic := parts[1]
			style := "educational"
			length := "medium"
			language := "en"

			if len(parts) > 2 {
				style = parts[2]
			}
			if len(parts) > 3 {
				length = parts[3]
			}
			if len(parts) > 4 {
				language = parts[4]
			}

			if err := cli.generatePrompt(topic, style, length, language); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

		case "sources":
			if err := cli.getSources(); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

		case "languages":
			if err := cli.getLanguages(); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

		case "health":
			if err := cli.health(); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}

		case "web":
			url := "http://localhost:8080/examples/encyclopedia_interface.html"
			fmt.Printf("🌐 Opening web interface: %s\n", url)
			if err := cli.openURL(url); err != nil {
				fmt.Printf("❌ Failed to open web interface: %v\n", err)
			}

		case "clear":
			cli.clearScreen()

		default:
			fmt.Printf("❌ Unknown command: %s\n", command)
			fmt.Println("Type 'help' for available commands")
		}
	}
}

func (cli *EncyclopediaCLI) showHelp() {
	fmt.Printf("\n📚 Encyclopedia Agent CLI - Available Commands:\n")
	fmt.Printf("=============================================\n\n")
	fmt.Printf("🔍 search <query> [source] [language] [max_results]\n")
	fmt.Printf("   Search for encyclopedia articles\n")
	fmt.Printf("   Example: search 'artificial intelligence' wikipedia en 5\n\n")

	fmt.Printf("📖 article <title> [source] [language] [max_length]\n")
	fmt.Printf("   Retrieve a specific article\n")
	fmt.Printf("   Example: article 'Machine Learning' wikipedia en 1500\n\n")

	fmt.Printf("✍️  prompt <topic> [style] [length] [language]\n")
	fmt.Printf("   Generate encyclopedia-style prompts\n")
	fmt.Printf("   Example: prompt 'neural networks' academic long en\n\n")

	fmt.Printf("📚 sources\n")
	fmt.Printf("   Show available encyclopedia sources\n\n")

	fmt.Printf("🌍 languages\n")
	fmt.Printf("   Show supported languages\n\n")

	fmt.Printf("🏥 health\n")
	fmt.Printf("   Check encyclopedia service health\n\n")

	fmt.Printf("🌐 web\n")
	fmt.Printf("   Open web interface in browser\n\n")

	fmt.Printf("🧹 clear\n")
	fmt.Printf("   Clear the screen\n\n")

	fmt.Printf("❓ help\n")
	fmt.Printf("   Show this help message\n\n")

	fmt.Printf("🚪 quit/exit\n")
	fmt.Printf("   Exit the CLI\n\n")

	fmt.Printf("📝 Notes:\n")
	fmt.Printf("   • Default source: wikipedia\n")
	fmt.Printf("   • Default language: en\n")
	fmt.Printf("   • Default max_results: 5\n")
	fmt.Printf("   • Default max_length: 2000\n")
	fmt.Printf("   • Default style: educational\n")
	fmt.Printf("   • Default length: medium\n\n")
}

func (cli *EncyclopediaCLI) clearScreen() {
	switch runtime.GOOS {
	case "darwin", "linux":
		fmt.Print("\033[H\033[2J")
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	cli := NewEncyclopediaCLI()

	// Check if server is running
	if _, err := cli.makeRequest("GET", "/encyclopedia/health", nil); err != nil {
		fmt.Printf("❌ Encyclopedia service is not running: %v\n", err)
		fmt.Println("Please start the server first: go run main.go")
		os.Exit(1)
	}

	// Check command line arguments
	if len(os.Args) > 1 {
		cli.handleCommandLine(os.Args[1:])
	} else {
		cli.interactiveMode()
	}
}

func (cli *EncyclopediaCLI) handleCommandLine(args []string) {
	if len(args) == 0 {
		return
	}

	command := strings.ToLower(args[0])

	switch command {
	case "search":
		if len(args) < 2 {
			fmt.Println("❌ Usage: encyclopedia search <query> [source] [language] [max_results]")
			os.Exit(1)
		}
		query := args[1]
		source := "all"
		language := "en"
		maxResults := 5

		if len(args) > 2 {
			source = args[2]
		}
		if len(args) > 3 {
			language = args[3]
		}
		if len(args) > 4 {
			if max, err := fmt.Sscanf(args[4], "%d", &maxResults); err != nil || max != 1 {
				fmt.Println("❌ Invalid max_results number")
				os.Exit(1)
			}
		}

		if err := cli.search(query, source, language, maxResults); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

	case "article":
		if len(args) < 2 {
			fmt.Println("❌ Usage: encyclopedia article <title> [source] [language] [max_length]")
			os.Exit(1)
		}
		title := args[1]
		source := "wikipedia"
		language := "en"
		maxLength := 2000

		if len(args) > 2 {
			source = args[2]
		}
		if len(args) > 3 {
			language = args[3]
		}
		if len(args) > 4 {
			if max, err := fmt.Sscanf(args[4], "%d", &maxLength); err != nil || max != 1 {
				fmt.Println("❌ Invalid max_length number")
				os.Exit(1)
			}
		}

		if err := cli.getArticle(title, source, language, maxLength); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

	case "prompt":
		if len(args) < 2 {
			fmt.Println("❌ Usage: encyclopedia prompt <topic> [style] [length] [language]")
			os.Exit(1)
		}
		topic := args[1]
		style := "educational"
		length := "medium"
		language := "en"

		if len(args) > 2 {
			style = args[2]
		}
		if len(args) > 3 {
			length = args[3]
		}
		if len(args) > 4 {
			language = args[4]
		}

		if err := cli.generatePrompt(topic, style, length, language); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

	case "sources":
		if err := cli.getSources(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

	case "languages":
		if err := cli.getLanguages(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

	case "health":
		if err := cli.health(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

	case "web":
		url := "http://localhost:8080/examples/encyclopedia_interface.html"
		fmt.Printf("🌐 Opening web interface: %s\n", url)
		if err := cli.openURL(url); err != nil {
			fmt.Printf("❌ Failed to open web interface: %v\n", err)
			os.Exit(1)
		}

	case "help":
		cli.showHelp()

	default:
		fmt.Printf("❌ Unknown command: %s\n", command)
		fmt.Println("Available commands: search, article, prompt, sources, languages, health, web, help")
		fmt.Println("Use 'encyclopedia help' for detailed usage information")
		os.Exit(1)
	}
}
