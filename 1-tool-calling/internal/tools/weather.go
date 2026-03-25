package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

type WeatherInput struct {
	City string `json:"city"`
	Unit string `json:"unit"`
}

type Weather struct{}

func (w Weather) Name() string { return "get_weather" }

func (w Weather) Description() string {
	return `Get weather. Input: {"city": "Tokyo", "unit": "celsius"}`
}

func (w Weather) Call(ctx context.Context, input string) (string, error) {
	var inp WeatherInput
	if err := json.Unmarshal([]byte(input), &inp); err != nil {
		return "", fmt.Errorf("invalid JSON: use {\"city\": \"NAME\", \"unit\": \"celsius\"}")
	}
	if inp.Unit == "" {
		inp.Unit = "celsius"
	}
	return fmt.Sprintf(`{"city": "%s", "temp": 22, "unit": "%s", "cond": "sunny"}`, inp.City, inp.Unit), nil
}
