# Practical Look at Tool Calling

This demo shows how to build an agent that can interact with external "tools" to solve tasks it can't handle alone (like getting real-time weather).

## How it works

I've set up two mock tools in `internal/tools/`:
- `get_weather`: Takes a JSON input for city and unit.
- `web_search`: A simple string-based search for travel tips.

The agent uses the `qwen3.5:9b` model via Ollama. It follows a system prompt that enforces a **Thought -> Action -> Action Input -> Final Answer** loop. This is a classic ReAct pattern that helps the model "think" through its steps.

## Running the demo

Make sure you have Ollama running with the model pulled:
```bash
ollama pull qwen3.5:9b
```

Then just run the main file:
```bash
go run main.go
```

## Practical observations

- **System Prompting**: The system prompt is where most of the magic happens. It tells the LLM exactly what format to use so `langchaingo`'s executor can parse the tool calls.
- **Mocking**: For learning, I'm just mocking the tool outputs. In a real scenario, you'd swap the `Call` method logic with an actual API request to OpenWeatherMap or Tavily.
- **Log Handler**: I'm using `callbacks.LogHandler{}` so you can see the model's internal reasoning steps in the terminal.

For more details on this implementation, see my post: [Building AI Agents with Go, LangChainGo, and Ollama](https://iampavel.dev/blog/go-ai-agents-langchaingo-ollama)
