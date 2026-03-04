package analyst

import (
	"strings"

	"github.com/dhruvinpm/taraka-bot/memory"
)

type Deduplicator struct {
	store *memory.Store
}

func NewDeduplicator(store *memory.Store) *Deduplicator {
	return &Deduplicator{store: store}
}

func (d *Deduplicator) IsDuplicate(email string) bool {
	exists, err := d.store.EmailExists(email)
	if err != nil {
		return false
	}
	return exists
}

func (d *Deduplicator) FindSimilar(company string) ([]*memory.Lead, error) {
	rows, err := d.store.Query(
		`SELECT id, name, email, phone, company, website, country, source, business_type, score, status, gmail_thread_id, gmail_message_id, emails_sent, follow_ups_sent, reply_received, notes, created_at, updated_at FROM leads WHERE company LIKE ?`,
		"%"+strings.ToLower(company)+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []*memory.Lead
	for rows.Next() {
		l := &memory.Lead{}
		var threadID, messageID interface{}
		if err := rows.Scan(&l.ID, &l.Name, &l.Email, &l.Phone, &l.Company, &l.Website, &l.Country, &l.Source, &l.BusinessType, &l.Score, &l.Status, &threadID, &messageID, &l.EmailsSent, &l.FollowUpsSent, &l.ReplyReceived, &l.Notes, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		leads = append(leads, l)
	}
	return leads, rows.Err()
}
