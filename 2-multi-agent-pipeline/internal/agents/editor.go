package agents

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type EditorAgent struct {
	BaseAgent
}

func NewEditor(llm llms.Model) *EditorAgent {
	return &EditorAgent{BaseAgent: NewBaseAgent(llm)}
}

func (e *EditorAgent) ExecuteWithStream(ctx context.Context, draft string, handler StreamHandler) (string, error) {
	prompt := fmt.Sprintf(`Edit this article for grammar, clarity, structure, and tone:

%s

Output the polished version directly. No commentary.`, draft)

	return e.callLLM(ctx, prompt, "You are a meticulous editor.", 0.3, 2000, handler, "Editor")
}
