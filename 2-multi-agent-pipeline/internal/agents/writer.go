package agents

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

type WriterAgent struct {
	llm llms.Model
}

func NewWriter(llm llms.Model) *WriterAgent {
	return &WriterAgent{llm: llm}
}

func (w *WriterAgent) ExecuteWithStream(ctx context.Context, research string, handler StreamHandler) (string, error) {
	prompt := fmt.Sprintf(`Transform these research notes into an engaging article:

%s

Write with compelling introduction, clear section headings (Markdown ##), factual accuracy, strong conclusion, and professional tone.`, research)

	return w.callLLM(ctx, prompt, "You are a professional writer.", 0.7, handler)
}

func (w *WriterAgent) callLLM(ctx context.Context, prompt, system string, temp float64, handler StreamHandler) (string, error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, system),
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}

	if handler != nil {
		result, err := w.streamLLM(ctx, content, temp, handler)
		if err == nil {
			return result, nil
		}
	}

	return w.simpleLLM(ctx, content, temp)
}

func (w *WriterAgent) streamLLM(ctx context.Context, content []llms.MessageContent, temp float64, handler StreamHandler) (string, error) {
	var response strings.Builder
	chunkCount := 0

	_, err := w.llm.GenerateContent(ctx, content,
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

	log.Printf("[Writer] Chunks: %d, Bytes: %d", chunkCount, response.Len())
	return response.String(), nil
}

func (w *WriterAgent) simpleLLM(ctx context.Context, content []llms.MessageContent, temp float64) (string, error) {
	resp, err := w.llm.GenerateContent(ctx, content,
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
