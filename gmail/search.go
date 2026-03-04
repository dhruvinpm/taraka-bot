package gmail

import (
	"encoding/base64"
	"strings"

	gmailapi "google.golang.org/api/gmail/v1"
)

func (c *Client) SearchReplies(threadID string) ([]string, error) {
	thread, err := c.GetThread(threadID)
	if err != nil {
		return nil, err
	}
	var replies []string
	for _, msg := range thread.Messages {
		body := GetMessageBody(msg)
		if body != "" {
			replies = append(replies, body)
		}
	}
	return replies, nil
}

func GetMessageBody(msg *gmailapi.Message) string {
	if msg.Payload == nil {
		return ""
	}
	if msg.Payload.Body != nil && msg.Payload.Body.Data != "" {
		return decodeBase64URL(msg.Payload.Body.Data)
	}
	for _, part := range msg.Payload.Parts {
		if part.MimeType == "text/plain" && part.Body != nil {
			return decodeBase64URL(part.Body.Data)
		}
	}
	return ""
}

func GetMessageSubject(msg *gmailapi.Message) string {
	if msg.Payload == nil {
		return ""
	}
	for _, h := range msg.Payload.Headers {
		if h.Name == "Subject" {
			return h.Value
		}
	}
	return ""
}

func GetMessageFrom(msg *gmailapi.Message) string {
	if msg.Payload == nil {
		return ""
	}
	for _, h := range msg.Payload.Headers {
		if h.Name == "From" {
			return h.Value
		}
	}
	return ""
}

func decodeBase64URL(s string) string {
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(b)
}
