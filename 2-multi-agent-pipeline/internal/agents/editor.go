package agents

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

type EditorAgent struct {
	llm llms.Model
}

func NewEditor(llm llms.Model) *EditorAgent {
	return &EditorAgent{llm: llm}
}

func (e *EditorAgent) ExecuteWithStream(ctx context.Context, draft string, handler StreamHandler) (string, error) {
	prompt := fmt.Sprintf(`Edit this article for grammar, clarity, structure, and tone:

%s

Output the polished version directly. No commentary.`, draft)

	return e.callLLM(ctx, prompt, "You are a meticulous editor.", 0.3, handler)
}

func (e *EditorAgent) callLLM(ctx context.Context, prompt, system string, temp float64, handler StreamHandler) (string, error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, system),
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}

	if handler != nil {
		result, err := e.streamLLM(ctx, content, temp, handler)
		if err == nil {
			return result, nil
		}
	}

	return e.simpleLLM(ctx, content, temp)
}

func (e *EditorAgent) streamLLM(ctx context.Context, content []llms.MessageContent, temp float64, handler StreamHandler) (string, error) {
	var response strings.Builder
	chunkCount := 0

	_, err := e.llm.GenerateContent(ctx, content,
		llms.WithTemperature(temp),
		llms.WithMaxTokens(2000),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			str := string(chunk)
			response.WriteString(str)
			handler(str)
			chunkCount++
			return nil
		}),
	)

	if err != nil {
		return "", err
	}

	log.Printf("[Editor] Chunks: %d, Bytes: %d", chunkCount, response.Len())
	return response.String(), nil
}

func (e *EditorAgent) simpleLLM(ctx context.Context, content []llms.MessageContent, temp float64) (string, error) {
	resp, err := e.llm.GenerateContent(ctx, content,
		llms.WithTemperature(temp),
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
