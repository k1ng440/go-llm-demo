# Practical AI Agent Demos

This is a collection of hands-on demos for building AI Agents. I'm focusing on using Go with `langchaingo` and Ollama to see how these pieces actually fit together.

## Demos

- **[1-tool-calling](./1-tool-calling/)**: A basic agent that uses tools to answer questions. It uses a ReAct-style loop to decide whether it needs to fetch weather data or search the web before giving a final answer.

## Resources

For a deep dive into the concepts behind these demos, check out:
- [Building AI Agents with Go, LangChainGo, and Ollama](https://iampavel.dev/blog/go-ai-agents-langchaingo-ollama)

## Practical Setup

Most of these demos rely on:
- **Go**: 1.22+
- **Ollama**: Running locally with models like `qwen3.5:9b`.
- **langchaingo**: For the agent orchestration and tool interfaces.

## Notes on the approach

Instead of reaching for complex frameworks immediately, I'm starting with the basics of tool definition and system prompting. You'll see that a lot of the "intelligence" comes from how the system message is structured to guide the LLM's reasoning process.
