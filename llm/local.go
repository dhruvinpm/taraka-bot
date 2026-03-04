package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocalProvider struct {
	url   string
	model string
}

type localRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model,omitempty"`
}

type localResponse struct {
	Content string `json:"content"`
}

func NewLocalProvider(url, model string) *LocalProvider {
	return &LocalProvider{url: url, model: model}
}

func (l *LocalProvider) Generate(ctx context.Context, prompt string) (string, error) {
	reqBody, err := json.Marshal(localRequest{Prompt: prompt, Model: l.model})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", l.url, bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("local llm request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var lr localResponse
	if err := json.Unmarshal(body, &lr); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	return lr.Content, nil
}

func (l *LocalProvider) Name() string {
	return "local"
}
