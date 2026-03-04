package gmail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmailapi "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Client struct {
	service *gmailapi.Service
	address string
}

func NewClient(ctx context.Context, credentialsFile, tokenFile, address string) (*Client, error) {
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("read credentials: %w", err)
	}
	config, err := google.ConfigFromJSON(b, gmailapi.GmailSendScope, gmailapi.GmailReadonlyScope, gmailapi.GmailModifyScope)
	if err != nil {
		return nil, fmt.Errorf("parse credentials: %w", err)
	}
	httpClient, err := getHTTPClient(ctx, config, tokenFile)
	if err != nil {
		return nil, fmt.Errorf("get http client: %w", err)
	}
	return NewClientWithHTTP(ctx, httpClient, address)
}

func NewClientWithHTTP(ctx context.Context, httpClient *http.Client, address string) (*Client, error) {
	svc, err := gmailapi.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("create gmail service: %w", err)
	}
	return &Client{service: svc, address: address}, nil
}

func getHTTPClient(ctx context.Context, config *oauth2.Config, tokenFile string) (*http.Client, error) {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("load token: %w", err)
	}
	return config.Client(ctx, tok), nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	return tok, json.NewDecoder(f).Decode(tok)
}

func (c *Client) Address() string {
	return c.address
}

// NewNullClient returns a Client with no underlying service (for use when Gmail is not configured).
func NewNullClient() *Client {
	return &Client{}
}
