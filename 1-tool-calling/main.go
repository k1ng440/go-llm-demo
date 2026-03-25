package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/k1ng440/go-llm-demo/1-tool-calling/internal/tools"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
	langchaintools "github.com/tmc/langchaingo/tools"
)

func main() {
	ctx := context.Background()

	llm, err := ollama.New(
		ollama.WithModel("qwen3.5:9b"),
		ollama.WithPredictMirostat(0),
		ollama.WithPullModel(),
	)
	if err != nil {
		log.Fatal(err)
	}

	systemPrompt := strings.TrimSpace(`
You are a helpful assistant that uses tools to answer questions.

IMPORTANT: You must follow this exact format:
For using a tool:
Thought: [your reasoning]
Action: [tool name]
Action Input: [tool input]

For final answer:
Thought: I now know the final answer
Final Answer: [your answer]

Always use "Final Answer:" to indicate your final response.

format the respones as Markdown
	`)

	agent := agents.NewConversationalAgent(
		llm,
		[]langchaintools.Tool{
			tools.Weather{},
			tools.Search{},
		},
		agents.NewOpenAIOption().WithSystemMessage(systemPrompt),
	)

	// We use agents.WithCallbacksHandler(callbacks.LogHandler{}) to see the "Thinking" process.
	executor := agents.NewExecutor(
		agent,
		agents.WithMemory(memory.NewConversationBuffer()),
		agents.WithMaxIterations(5),
		agents.WithCallbacksHandler(callbacks.LogHandler{}),
	)

	input := "What is the weather in Tokyo? Also, search for travel tips."
	fmt.Printf("User: %s\n", input)

	response, err := chains.Run(ctx, executor, input)
	if err != nil {
		log.Fatalf("Agent failed: %v", err)
	}
	fmt.Printf("\nFinal Agent Response: %s\n", response)
}
