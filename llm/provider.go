package llm

import "context"

type Provider interface {
	Generate(ctx context.Context, prompt string) (string, error)
	Name() string
}
