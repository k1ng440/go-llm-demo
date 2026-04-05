package agents

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/llms"
	langchaintools "github.com/tmc/langchaingo/tools"
)

type StreamHandler func(chunk string)

type SimpleResearcherAgent struct {
	llm   llms.Model
	tools []langchaintools.Tool
}

func NewSimpleResearcher(llm llms.Model, tools []langchaintools.Tool) *SimpleResearcherAgent {
	return &SimpleResearcherAgent{llm: llm, tools: tools}
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
	return r.callLLM(ctx, prompt, "You are a research specialist.", handler)
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

func (r *SimpleResearcherAgent) callLLM(ctx context.Context, prompt, system string, handler StreamHandler) (string, error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, system),
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}

	if handler != nil {
		result, err := r.streamLLM(ctx, content, handler)
		if err == nil {
			return result, nil
		}
	}

	return r.simpleLLM(ctx, content)
}

func (r *SimpleResearcherAgent) streamLLM(ctx context.Context, content []llms.MessageContent, handler StreamHandler) (string, error) {
	var response strings.Builder
	chunkCount := 0

	_, err := r.llm.GenerateContent(ctx, content,
		llms.WithTemperature(0.5),
		llms.WithMaxTokens(2000),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			str := string(chunk)
			response.WriteString(str)
			handler(str)
			chunkCount++
			if chunkCount == 1 {
				log.Printf("[LLM] First chunk received (%d bytes)", len(chunk))
			}
			return nil
		}),
	)

	if err != nil {
		return "", err
	}

	log.Printf("[LLM] Total chunks: %d, Total bytes: %d", chunkCount, response.Len())
	return response.String(), nil
}

func (r *SimpleResearcherAgent) simpleLLM(ctx context.Context, content []llms.MessageContent) (string, error) {
	resp, err := r.llm.GenerateContent(ctx, content,
		llms.WithTemperature(0.5),
		llms.WithMaxTokens(2000),
	)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response")
	}

	return resp.Choices[0].Content, nil
}
