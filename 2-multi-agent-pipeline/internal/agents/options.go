package agents

import "github.com/tmc/langchaingo/llms"

// Option configures an Agent
type Option func(*config)

type config struct {
	temperature  float64
	maxTokens    int
	systemPrompt string
	logPrefix    string
}

func defaultConfig() config {
	return config{
		temperature: 0.5,
		maxTokens:   2000,
		logPrefix:   "Agent",
	}
}

func WithTemperature(t float64) Option {
	return func(c *config) {
		c.temperature = t
	}
}

func WithMaxTokens(n int) Option {
	return func(c *config) {
		c.maxTokens = n
	}
}

func WithSystemPrompt(p string) Option {
	return func(c *config) {
		c.systemPrompt = p
	}
}

func WithLogPrefix(p string) Option {
	return func(c *config) {
		c.logPrefix = p
	}
}

// LLMCaller handles LLM interactions with configuration
type LLMCaller struct {
	llm    llms.Model
	config config
}

func NewLLMCaller(llm llms.Model, opts ...Option) LLMCaller {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return LLMCaller{llm: llm, config: cfg}
}
