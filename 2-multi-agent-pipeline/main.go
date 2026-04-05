package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/k1ng440/go-llm-demo/2-multi-agent-pipeline/internal/agents"
	"github.com/k1ng440/go-llm-demo/2-multi-agent-pipeline/internal/tools"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	langchaintools "github.com/tmc/langchaingo/tools"
)

type Pipeline struct {
	researcher *agents.SimpleResearcherAgent
	writer     *agents.WriterAgent
	editor     *agents.EditorAgent
}

func NewPipeline() (*Pipeline, error) {
	model := "minimax-m2.7:cloud"

	llm, err := ollama.New(
		ollama.WithModel(model),
		ollama.WithPredictMirostat(0),
	)
	if err != nil {
		return nil, err
	}

	testCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = llm.Call(testCtx, "Hi", llms.WithTemperature(0.1))
	if err != nil {
		return nil, fmt.Errorf("model %s not available (run: ollama pull %s): %w", model, model, err)
	}

	searchTools := []langchaintools.Tool{tools.NewSearch()}

	return &Pipeline{
		researcher: agents.NewSimpleResearcher(llm, searchTools),
		writer:     agents.NewWriter(llm),
		editor:     agents.NewEditor(llm),
	}, nil
}

func (p *Pipeline) Run(ctx context.Context, topic string) (string, error) {
	log.Println("[1/3] Research...")
	log.Println("[Streaming] Starting...")
	research, err := p.researcher.ExecuteWithStream(ctx, topic, func(chunk string) {
		fmt.Print(chunk)
		os.Stdout.Sync()
	})
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("research: %w", err)
	}
	log.Printf("      -> %d chars", len(research))
	log.Printf("      -> Preview: %.100s...", research)

	log.Println("")
	log.Println("[2/3] Writing...")
	log.Printf("      -> Input: %d chars of research", len(research))
	log.Println("[Streaming] Starting...")
	draft, err := p.writer.ExecuteWithStream(ctx, research, func(chunk string) {
		fmt.Print(chunk)
		os.Stdout.Sync()
	})
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("writing: %w", err)
	}
	log.Printf("      -> %d chars", len(draft))
	log.Printf("      -> Preview: %.100s...", draft)

	log.Println("")
	log.Println("[3/3] Editing...")
	log.Printf("      -> Input: %d chars of draft", len(draft))
	log.Println("[Streaming] Starting...")
	final, err := p.editor.ExecuteWithStream(ctx, draft, func(chunk string) {
		fmt.Print(chunk)
		os.Stdout.Sync()
	})
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("editing: %w", err)
	}
	log.Printf("      -> %d chars", len(final))
	log.Printf("      -> Preview: %.100s...", final)

	return final, nil
}

func main() {
	log.Println("Multi-Agent Pipeline Demo")
	log.Println("=========================")

	pipeline, err := NewPipeline()
	if err != nil {
		log.Fatal(err)
	}

	topic := "The impact of AI on software engineering workflows"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	result, err := pipeline.Run(ctx, topic)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
	fmt.Println("=========================")
	fmt.Println("FINAL OUTPUT")
	fmt.Println("=========================")
	fmt.Println(result)
	fmt.Printf("\nTotal: %d characters\n", len(result))
}
