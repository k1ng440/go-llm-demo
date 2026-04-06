package agents

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type WriterAgent struct {
	BaseAgent
}

func NewWriter(llm llms.Model) *WriterAgent {
	return &WriterAgent{BaseAgent: NewBaseAgent(llm)}
}

func (w *WriterAgent) ExecuteWithStream(ctx context.Context, research string, handler StreamHandler) (string, error) {
	prompt := fmt.Sprintf(`Transform these research notes into an engaging article:

%s

Write with compelling introduction, clear section headings (Markdown ##), factual accuracy, strong conclusion, and professional tone.`, research)

	return w.callLLM(ctx, prompt, "You are a professional writer.", 0.7, 2000, handler, "Writer")
}
