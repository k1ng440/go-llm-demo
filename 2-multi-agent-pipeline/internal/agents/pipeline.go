package agents

import (
	"context"
	"fmt"
)

// PipelineStep represents one stage in the pipeline
type PipelineStep struct {
	Name   string
	Agent  Agent
	Prompt func(string) string // Transform previous output into prompt
}

// PipelineRunner executes a sequence of agents
type PipelineRunner struct {
	steps []PipelineStep
}

func NewPipelineRunner(steps ...PipelineStep) *PipelineRunner {
	return &PipelineRunner{steps: steps}
}

// Run executes the pipeline sequentially
func (p *PipelineRunner) Run(ctx context.Context, initialInput string, onStep func(name, output string)) (string, error) {
	input := initialInput

	for i, step := range p.steps {
		select {
		case <-ctx.Done():
			return "", &Error{Agent: step.Name, Phase: "execution", Cause: ErrContextCancelled}
		default:
		}

		// Transform input through prompt function
		prompt := input
		if step.Prompt != nil {
			prompt = step.Prompt(input)
		}

		// Execute agent (no streaming for cleaner abstraction)
		output, err := step.Agent.ExecuteWithStream(ctx, prompt, nil)
		if err != nil {
			return "", &Error{Agent: step.Name, Phase: "execution", Cause: err}
		}

		if onStep != nil {
			onStep(step.Name, output)
		}

		input = output
		_ = i // step number available if needed
	}

	return input, nil
}

// RunWithStream executes with streaming for the final step only
func (p *PipelineRunner) RunWithStream(ctx context.Context, initialInput string, handler StreamHandler) (string, error) {
	if len(p.steps) == 0 {
		return "", fmt.Errorf("no steps in pipeline")
	}

	input := initialInput

	// Run all but last step without streaming
	for _, step := range p.steps[:len(p.steps)-1] {
		output, err := step.Agent.ExecuteWithStream(ctx, input, nil)
		if err != nil {
			return "", &Error{Agent: step.Name, Phase: "execution", Cause: err}
		}
		input = output
	}

	// Final step with streaming
	lastStep := p.steps[len(p.steps)-1]
	return lastStep.Agent.ExecuteWithStream(ctx, input, handler)
}
