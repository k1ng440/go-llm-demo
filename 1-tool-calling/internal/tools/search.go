package tools

import (
	"context"
	"fmt"
)

type Search struct{}

func (s Search) Name() string { return "web_search" }

func (s Search) Description() string {
	return `Search the web for travel tips. Input is a search query string.`
}

func (s Search) Call(ctx context.Context, input string) (string, error) {
	return fmt.Sprintf(`{"query": "%s", "results": ["Tokyo travel guide", "Top 10 things to do in Tokyo"]}`, input), nil
}
