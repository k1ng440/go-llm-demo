package agents

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tmc/langchaingo/llms"
)

type StreamHandler func(chunk string)

// Agent is the common interface for all pipeline agents
type Agent interface {
	ExecuteWithStream(ctx context.Context, input string, handler StreamHandler) (string, error)
}

type BaseAgent struct {
	llm llms.Model
}

func NewBaseAgent(llm llms.Model) BaseAgent {
	return BaseAgent{llm: llm}
}

func (b *BaseAgent) callLLM(ctx context.Context, prompt, system string, temp float64, maxTokens int, handler StreamHandler, logPrefix string) (string, error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, system),
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}

	if handler != nil {
		result, err := b.streamLLM(ctx, content, temp, maxTokens, handler, logPrefix)
		if err == nil {
			return result, nil
		}
	}

	return b.simpleLLM(ctx, content, temp, maxTokens)
}

func (b *BaseAgent) streamLLM(ctx context.Context, content []llms.MessageContent, temp float64, maxTokens int, handler StreamHandler, logPrefix string) (string, error) {
	var response strings.Builder
	chunkCount := 0

	_, err := b.llm.GenerateContent(ctx, content,
		llms.WithTemperature(temp),
		llms.WithMaxTokens(maxTokens),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if len(chunk) == 0 {
				// Log empty chunks to help diagnose streaming issues
				if chunkCount < 5 && logPrefix != "" {
					log.Printf("[%s] Empty chunk #%d received", logPrefix, chunkCount+1)
				}
				chunkCount++
				return nil
			}
			str := string(chunk)
			response.WriteString(str)
			if handler != nil {
				handler(str)
			}
			chunkCount++
			if chunkCount == 1 && logPrefix != "" {
				log.Printf("[%s] First chunk received (%d bytes): %.50s...", logPrefix, len(chunk), str)
			}
			return nil
		}),
	)

	if err != nil {
		return "", err
	}

	if logPrefix != "" {
		log.Printf("[%s] Total chunks: %d, Total bytes: %d", logPrefix, chunkCount, response.Len())
	}
	return response.String(), nil
}

func (b *BaseAgent) simpleLLM(ctx context.Context, content []llms.MessageContent, temp float64, maxTokens int) (string, error) {
	const maxRetries = 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		resp, err := b.llm.GenerateContent(ctx, content,
			llms.WithTemperature(temp),
			llms.WithMaxTokens(maxTokens),
		)
		if err == nil {
			if len(resp.Choices) == 0 {
				return "", ErrNoResponse
			}
			return resp.Choices[0].Content, nil
		}

		lastErr = err
		if i < maxRetries-1 {
			log.Printf("[LLM] Retry %d/%d after error: %v", i+1, maxRetries, err)
			time.Sleep(time.Second * time.Duration(i+1)) // Exponential backoff
		}
	}

	return "", fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
