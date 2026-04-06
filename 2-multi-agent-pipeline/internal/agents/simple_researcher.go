package agents

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	langchaintools "github.com/tmc/langchaingo/tools"
)

type SimpleResearcherAgent struct {
	BaseAgent
	tools []langchaintools.Tool
}

func NewSimpleResearcher(llm llms.Model, tools []langchaintools.Tool) *SimpleResearcherAgent {
	return &SimpleResearcherAgent{
		BaseAgent: NewBaseAgent(llm),
		tools:     tools,
	}
}

func (r *SimpleResearcherAgent) findSearchTool() langchaintools.Tool {
	for _, t := range r.tools {
		if t.Name() == "web_search" {
			return t
		}
	}
	return nil
}

func (r *SimpleResearcherAgent) ExecuteWithStream(ctx context.Context, topic string, handler StreamHandler) (string, error) {
	prompt := r.buildPrompt(ctx, topic)
	return r.callLLM(ctx, prompt, "You are a research specialist.", 0.5, 2000, handler, "Researcher")
}

func (r *SimpleResearcherAgent) buildPrompt(ctx context.Context, topic string) string {
	searchResults := r.getSearchResults(ctx, topic)

	if searchResults == "" {
		return fmt.Sprintf(`Research "%s". Provide comprehensive notes with facts, statistics, trends, and sources.`, topic)
	}

	return fmt.Sprintf(`Research "%s" based on these search results:

%s

Provide comprehensive notes with facts, statistics, trends, and sources.`, topic, searchResults)
}

func (r *SimpleResearcherAgent) getSearchResults(ctx context.Context, topic string) string {
	searchTool := r.findSearchTool()
	if searchTool == nil {
		return ""
	}

	results, err := searchTool.Call(ctx, topic)
	if err != nil {
		log.Printf("[Search] Failed: %v", err)
		return ""
	}

	log.Printf("[Search] Got %d chars", len(results))
	return results
}
