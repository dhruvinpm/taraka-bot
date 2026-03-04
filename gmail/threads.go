package gmail

import (
	"fmt"

	gmailapi "google.golang.org/api/gmail/v1"
)

func (c *Client) GetThread(threadID string) (*gmailapi.Thread, error) {
	return c.service.Users.Threads.Get("me", threadID).Do()
}

func (c *Client) ListThreads(query string, maxResults int64) ([]*gmailapi.Thread, error) {
	resp, err := c.service.Users.Threads.List("me").Q(query).MaxResults(maxResults).Do()
	if err != nil {
		return nil, fmt.Errorf("list threads: %w", err)
	}
	return resp.Threads, nil
}

func (c *Client) GetMessage(messageID string) (*gmailapi.Message, error) {
	return c.service.Users.Messages.Get("me", messageID).Do()
}

func (c *Client) MarkAsRead(messageID string) error {
	_, err := c.service.Users.Messages.Modify("me", messageID, &gmailapi.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}).Do()
	return err
}

func (c *Client) ListMessages(query string, maxResults int64) ([]*gmailapi.Message, error) {
	resp, err := c.service.Users.Messages.List("me").Q(query).MaxResults(maxResults).Do()
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	return resp.Messages, nil
}
