package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Search struct {
	client *http.Client
	apiKey string
}

func NewSearch() *Search {
	return &Search{
		client: &http.Client{Timeout: 10 * time.Second},
		apiKey: os.Getenv("BRAVE_API_KEY"),
	}
}

func (s *Search) Name() string { return "web_search" }

func (s *Search) Description() string {
	return "Search the web using Brave Search API. Returns titles, descriptions, and URLs."
}

type BraveResponse struct {
	Web struct {
		Results []struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			Description string `json:"description"`
		} `json:"results"`
	} `json:"web"`
}

func (s *Search) Call(ctx context.Context, input string) (string, error) {
	if s.apiKey == "" {
		return s.mockResults(input), nil
	}

	return s.braveSearch(ctx, input)
}

func (s *Search) braveSearch(ctx context.Context, input string) (string, error) {
	params := url.Values{"q": {input}, "count": {"10"}}
	url := "https://api.search.brave.com/res/v1/web/search?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("X-Subscription-Token", s.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("search request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("[Brave] Response: %.300s...", string(body))

	var braveResp BraveResponse
	if err := json.Unmarshal(body, &braveResp); err != nil {
		return "", fmt.Errorf("parse results: %w", err)
	}

	return s.formatResults(input, braveResp.Web.Results), nil
}

func (s *Search) formatResults(query string, results []struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}) string {
	formatted := make([]map[string]string, len(results))
	for i, r := range results {
		formatted[i] = map[string]string{
			"title":   r.Title,
			"snippet": r.Description,
			"source":  r.URL,
		}
	}

	output := map[string]interface{}{
		"query":   query,
		"results": formatted,
	}

	data, _ := json.Marshal(output)
	return string(data)
}

func (s *Search) mockResults(input string) string {
	return fmt.Sprintf(`{"query": "%s", "results": [
		{"title": "Guide to %s", "snippet": "Key facts about %s", "source": "example.com"},
		{"title": "%s - Wikipedia", "snippet": "Encyclopedia article", "source": "wikipedia.org"},
		{"title": "Research on %s", "snippet": "Academic papers", "source": "research-db.example"}
	]}`, input, input, input, input, input)
}
