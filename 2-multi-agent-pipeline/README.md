# Multi-Agent Pipeline Demo

This demo implements a multi-agent orchestration system in Go using LangChainGo, based on the [Pipeline Pattern](https://github.com/k1ng440/go-llm-demo/tree/main/articles/go-multi-agent-orchestration.md).

## Overview

The system uses three specialized agents working in a pipeline:

1. **Researcher Agent** - Gathers information using search tools
2. **Writer Agent** - Transforms research notes into engaging prose
3. **Editor Agent** - Reviews and polishes the final output

## Architecture

```
User Input → Researcher → Writer → Editor → Final Output
                ↓           ↓         ↓
            [Search]    [Write]   [Edit]
```

## Features

- **Sequential Pipeline**: Research → Write → Edit
- **Error Handling**: Retry logic and fallback mechanisms
- **Partial Output Recovery**: Returns intermediate results if later stages fail
- **Structured Data**: JSON-based handoffs between agents (no context bleed)
- **Observability**: Tracing and logging at each pipeline stage
- **Parallel Research**: Can research multiple topics concurrently

## Usage

```bash
cd 2-multi-agent-pipeline
go mod tidy
go run main.go
```

## Project Structure

```
2-multi-agent-pipeline/
├── main.go                    # Pipeline orchestrator
├── internal/
│   ├── agents/
│   │   ├── researcher.go      # Researcher agent
│   │   ├── writer.go        # Writer agent
│   │   └── editor.go        # Editor agent
│   ├── tools/
│   │   └── search.go        # Search tool implementation
│   └── observability/
│       └── tracer.go        # Execution tracing
└── go.mod
```

## Configuration

### LLM Setup

The demo uses Ollama with the `qwen3.5:14b` model by default. Make sure you have Ollama running locally:

```bash
ollama pull qwen3.5:14b
```

You can change the model in `main.go`:

```go
llm, err := ollama.New(
    ollama.WithModel("your-model"),
    ollama.WithPredictMirostat(0),
    ollama.WithPullModel(),
)
```

### Search API Setup (Brave Search)

The Researcher agent uses **Brave Search API** for real web searches. To use it:

1. **Get an API key** from [Brave Search API](https://brave.com/search/api/)
   - Free tier: 2,000 queries/month
   - Paid tier: $3 per 1000 queries

2. **Set the API key** as an environment variable:

```bash
export BRAVE_API_KEY="your-api-key-here"
```

3. **Run the demo** - The tool will automatically use Brave Search:

```bash
export BRAVE_API_KEY="your-api-key"
go run main.go
```

Without the API key, the tool falls back to **mock data** (simulated search results) for demo purposes.

### Using with Nix

```bash
nix develop  # Enter the dev shell
export BRAVE_API_KEY="your-api-key"
go run main.go
```

## How It Works

The demo uses a **simplified multi-agent architecture** that's more reliable with local LLMs:

### Agent Types

1. **SimpleResearcherAgent** - Uses direct LLM calls (not agent framework) for maximum compatibility with local models like qwen3.5
2. **WriterAgent** - Direct LLM calls with structured prompts for article generation
3. **EditorAgent** - Direct LLM calls for editing and polishing

### Why Not Use Agent Framework?

The original article demonstrated LangChainGo's `ConversationalAgent` with tool calling. However, local models (like qwen3.5) often struggle with the strict format requirements:

```
Thought: ...
Action: ...
Action Input: ...
```

This implementation uses **direct LLM calls** instead, which is:
- More reliable with local models
- Faster (no parsing overhead)
- Simpler to debug

The agent-based approach is still available in `researcher.go` for models that support it well (GPT-4, Claude, etc.).

```
[Pipeline] Starting research phase...
[Trace] Researcher started | input: 45 chars
[Trace] Researcher completed in 3.2s | output: 850 chars
[Pipeline] Research complete (850 chars)
[Pipeline] Starting writing phase...
[Trace] Writer started | input: 850 chars
[Trace] Writer completed in 2.8s | output: 1200 chars
[Pipeline] Draft complete (1200 chars)
[Pipeline] Starting editing phase...
[Trace] Editor started | input: 1200 chars
[Trace] Editor completed in 1.5s | output: 1150 chars
[Pipeline] Final version ready (1150 chars)

=== FINAL OUTPUT ===
[Polished article about the topic]
```

## Patterns Demonstrated

1. **Agent Specialization**: Each agent has a focused system prompt
2. **Explicit Handoffs**: Structured data (JSON) passed between agents
3. **Error Resilience**: Retry with fallback, partial output recovery
4. **Observability**: Tracing wrapper for debugging
5. **Go Concurrency**: Parallel research using goroutines
