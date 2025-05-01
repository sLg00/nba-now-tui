package nbaAPI

import (
	"encoding/json"
	"fmt"
	filesystemops "github.com/sLg00/nba-now-tui/cmd/nba/filesystem"
	"github.com/sLg00/nba-now-tui/cmd/nba/pathManager"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const nbaNewsURL = "https://www.nba.com/news"

type NewsClient struct {
	client     *http.Client
	FileSystem filesystemops.FileSystemHandler
	Paths      pathManager.PathManager
}

type NewsArticle struct {
	Title     string
	Timestamp string
	URL       string
}

func NewNewsClient(fs filesystemops.FileSystemHandler, paths pathManager.PathManager) *NewsClient {
	return &NewsClient{
		client:     &http.Client{},
		FileSystem: fs,
		Paths:      paths,
	}
}

func (nc *NewsClient) FetchNews() ([]NewsArticle, error) {
	log.Println("Checking for cached articles")
	cacheFile := nc.Paths.GetFullPath("newsCacheFile", "")

	dirErr := nc.FileSystem.EnsureDirectoryExists(nc.Paths.GetFullPath("newsCachePath", ""))
	if dirErr != nil {
		log.Printf("error creating news directory: %v", dirErr)
	}

	if nc.FileSystem.FileExists(cacheFile) {
		cacheData, err := nc.FileSystem.ReadFile(cacheFile)
		if err != nil {
			return nil, err
		}
		var cachedArticles []NewsArticle
		if json.Unmarshal(cacheData, &cachedArticles) == nil && len(cachedArticles) > 0 {
			log.Println("Using cached articles")
			return cachedArticles, nil
		}
	}

	log.Println("Fetching news")

	defaultArticles := []NewsArticle{
		{
			Title:     "NBA.com News - Visit the website for the latest news",
			URL:       "https://www.nba.com/news",
			Timestamp: "Current",
		},
	}

	articles, err := nc.Scrape()
	if err != nil {
		log.Printf("error scraping news: %v", err)
		return defaultArticles, err
	}

	if len(articles) == 0 {
		log.Println("No articles found, using default")
		return defaultArticles, nil
	}

	if len(articles) > 0 {
		cacheData, err := json.Marshal(articles)
		if err != nil {
			return nil, err
		}
		writeErr := nc.FileSystem.WriteFile(cacheFile, cacheData)
		if writeErr != nil {
			log.Printf("error writing news cache: %v", writeErr)
		}
	}
	return articles, nil
}

func (nc *NewsClient) Scrape() ([]NewsArticle, error) {
	resp, err := nc.client.Get(nbaNewsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch news: %s", resp.StatusCode)
	}

	log.Println("Reading response body")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	content := string(body)
	log.Printf("Response body length: %d bytes", len(content))

	var articles []NewsArticle
	articlePattern := `<a\s+[^>]*href="(/news/[^"]+)"[^>]*>\s*<div[^>]*>.*?<div[^>]*>.*?<h\d[^>]*>(.*?)</h\d>`
	articleRegex := regexp.MustCompile(articlePattern)

	matches := articleRegex.FindAllStringSubmatch(content, -1)

	seenURLs := make(map[string]bool)
	for _, match := range matches {
		if len(match) >= 3 {
			url := "https://www.nba.com" + match[1]
			title := cleanHTML(match[2])

			if seenURLs[url] {
				continue
			}
			seenURLs[url] = true

			articles = append(articles, NewsArticle{
				Title:     title,
				Timestamp: time.Now().Format("2006-01-02"),
				URL:       url,
			})

			if len(articles) >= 10 {
				break
			}
		}
	}

	log.Printf("Found %d articles", len(articles))
	return articles, nil
}

func cleanHTML(html string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	text := re.ReplaceAllString(html, "")

	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	text = strings.ReplaceAll(text, "&#x27;", "'")
	text = strings.ReplaceAll(text, "&apos;", "'")
	text = strings.ReplaceAll(text, "&#96;", "`")
	text = strings.ReplaceAll(text, "&#x60;", "`")
	text = strings.ReplaceAll(text, "&ndash;", "–")
	text = strings.ReplaceAll(text, "&mdash;", "—")
	text = strings.ReplaceAll(text, "&#8211;", "–")
	text = strings.ReplaceAll(text, "&#8212;", "—")
	text = strings.ReplaceAll(text, "&nbsp;", " ") // Non-breaking space

	text = strings.TrimSpace(text)

	return text
}
