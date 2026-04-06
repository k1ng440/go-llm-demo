package agents

import (
	"errors"
	"fmt"
)

// Agent errors for better error handling
type Error struct {
	Agent string
	Phase string
	Cause error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s agent failed during %s: %v", e.Agent, e.Phase, e.Cause)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// Common errors
var (
	ErrNoResponse       = errors.New("LLM returned no response")
	ErrContextCancelled = errors.New("operation cancelled")
	ErrSearchFailed     = errors.New("search tool failed")
)

// Result holds agent output with metadata
type Result struct {
	Content   string
	CharCount int
	Phase     string
}

func NewResult(content, phase string) Result {
	return Result{
		Content:   content,
		CharCount: len(content),
		Phase:     phase,
	}
}
