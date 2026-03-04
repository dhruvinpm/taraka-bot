package llm

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Router struct {
	primary  Provider
	fallback Provider
}

func NewRouter(primary, fallback Provider) *Router {
	return &Router{primary: primary, fallback: fallback}
}

func (r *Router) Generate(ctx context.Context, prompt string) (string, error) {
	if r.isOnline() && r.primary != nil {
		result, err := r.primary.Generate(ctx, prompt)
		if err == nil {
			return result, nil
		}
	}
	if r.fallback != nil {
		return r.fallback.Generate(ctx, prompt)
	}
	return "", fmt.Errorf("no provider available")
}

func (r *Router) Name() string {
	return "router"
}

func (r *Router) isOnline() bool {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://www.google.com")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}
