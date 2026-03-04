package llm

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiProvider struct {
	client     *genai.Client
	model      string
	dailyLimit int
	usageCount int
}

func NewGeminiProvider(ctx context.Context, apiKey, model string, dailyLimit int) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("create gemini client: %w", err)
	}
	return &GeminiProvider{
		client:     client,
		model:      model,
		dailyLimit: dailyLimit,
	}, nil
}

func (g *GeminiProvider) Generate(ctx context.Context, prompt string) (string, error) {
	if g.usageCount >= g.dailyLimit {
		return "", fmt.Errorf("daily limit reached")
	}
	m := g.client.GenerativeModel(g.model)
	resp, err := m.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("generate content: %w", err)
	}
	g.usageCount++
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response")
	}
	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}

func (g *GeminiProvider) Name() string {
	return "gemini"
}

func (g *GeminiProvider) ResetUsage() {
	g.usageCount = 0
}

func (g *GeminiProvider) Close() {
	g.client.Close()
}
