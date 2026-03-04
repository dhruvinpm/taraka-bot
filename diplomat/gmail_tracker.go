package diplomat

import (
	localmail "github.com/dhruvinpm/taraka-bot/gmail"
	"github.com/dhruvinpm/taraka-bot/memory"
)

type GmailTracker struct {
	client *localmail.Client
	store  *memory.Store
}

func NewGmailTracker(client *localmail.Client, store *memory.Store) *GmailTracker {
	return &GmailTracker{client: client, store: store}
}

func (g *GmailTracker) TrackMessage(threadID, messageID string, leadID int64) error {
	return g.store.UpdateLeadEmail(leadID, threadID, messageID)
}

func (g *GmailTracker) GetThreadReplies(threadID string) ([]string, error) {
	return g.client.SearchReplies(threadID)
}
