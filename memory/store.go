package memory

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Lead struct {
	ID             int64
	Name           string
	Email          string
	Phone          string
	Company        string
	Website        string
	Country        string
	Source         string
	BusinessType   string
	Score          int
	Status         string
	GmailThreadID  string
	GmailMessageID string
	EmailsSent     int
	FollowUpsSent  int
	LastEmailedAt  *time.Time
	NextFollowUp   *time.Time
	ReplyReceived  bool
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

func (s *Store) migrate() error {
	tables := []string{
		createLeadsTable,
		createDailyLogTable,
		createCompanyKnowledgeTable,
		createCrawlHistoryTable,
		createConversationsTable,
	}
	for _, t := range tables {
		if _, err := s.db.Exec(t); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) UpsertLead(lead *Lead) error {
	_, err := s.db.Exec(`
		INSERT INTO leads (name, email, phone, company, website, country, source, business_type, score, status, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(email) DO UPDATE SET
			name=excluded.name,
			phone=excluded.phone,
			company=excluded.company,
			website=excluded.website,
			country=excluded.country,
			source=excluded.source,
			business_type=excluded.business_type,
			score=excluded.score,
			notes=excluded.notes,
			updated_at=CURRENT_TIMESTAMP`,
		lead.Name, lead.Email, lead.Phone, lead.Company, lead.Website,
		lead.Country, lead.Source, lead.BusinessType, lead.Score, lead.Status, lead.Notes)
	return err
}

func (s *Store) GetLeadByEmail(email string) (*Lead, error) {
	row := s.db.QueryRow(`SELECT id, name, email, phone, company, website, country, source, business_type, score, status, gmail_thread_id, gmail_message_id, emails_sent, follow_ups_sent, reply_received, notes, created_at, updated_at FROM leads WHERE email = ?`, email)
	l := &Lead{}
	var threadID, messageID sql.NullString
	err := row.Scan(&l.ID, &l.Name, &l.Email, &l.Phone, &l.Company, &l.Website, &l.Country, &l.Source, &l.BusinessType, &l.Score, &l.Status, &threadID, &messageID, &l.EmailsSent, &l.FollowUpsSent, &l.ReplyReceived, &l.Notes, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	l.GmailThreadID = threadID.String
	l.GmailMessageID = messageID.String
	return l, nil
}

func (s *Store) GetLeadsByStatus(status string, limit int) ([]*Lead, error) {
	rows, err := s.db.Query(`SELECT id, name, email, phone, company, website, country, source, business_type, score, status, gmail_thread_id, gmail_message_id, emails_sent, follow_ups_sent, reply_received, notes, created_at, updated_at FROM leads WHERE status = ? ORDER BY score DESC LIMIT ?`, status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var leads []*Lead
	for rows.Next() {
		l := &Lead{}
		var threadID, messageID sql.NullString
		if err := rows.Scan(&l.ID, &l.Name, &l.Email, &l.Phone, &l.Company, &l.Website, &l.Country, &l.Source, &l.BusinessType, &l.Score, &l.Status, &threadID, &messageID, &l.EmailsSent, &l.FollowUpsSent, &l.ReplyReceived, &l.Notes, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		l.GmailThreadID = threadID.String
		l.GmailMessageID = messageID.String
		leads = append(leads, l)
	}
	return leads, rows.Err()
}

func (s *Store) GetDueFollowUps() ([]*Lead, error) {
	rows, err := s.db.Query(`SELECT id, name, email, phone, company, website, country, source, business_type, score, status, gmail_thread_id, gmail_message_id, emails_sent, follow_ups_sent, reply_received, notes, created_at, updated_at FROM leads WHERE next_follow_up <= datetime('now') AND status IN ('emailed','follow_up_1','follow_up_2') AND reply_received = 0`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var leads []*Lead
	for rows.Next() {
		l := &Lead{}
		var threadID, messageID sql.NullString
		if err := rows.Scan(&l.ID, &l.Name, &l.Email, &l.Phone, &l.Company, &l.Website, &l.Country, &l.Source, &l.BusinessType, &l.Score, &l.Status, &threadID, &messageID, &l.EmailsSent, &l.FollowUpsSent, &l.ReplyReceived, &l.Notes, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		l.GmailThreadID = threadID.String
		l.GmailMessageID = messageID.String
		leads = append(leads, l)
	}
	return leads, rows.Err()
}

func (s *Store) UpdateLeadStatus(id int64, status string) error {
	_, err := s.db.Exec(`UPDATE leads SET status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, status, id)
	return err
}

func (s *Store) UpdateLeadEmail(id int64, threadID, messageID string) error {
	_, err := s.db.Exec(`UPDATE leads SET gmail_thread_id=?, gmail_message_id=?, emails_sent=emails_sent+1, last_emailed_at=CURRENT_TIMESTAMP, updated_at=CURRENT_TIMESTAMP WHERE id=?`, threadID, messageID, id)
	return err
}

func (s *Store) UpdateFollowUp(id int64, nextFollowUp time.Time, status string) error {
	_, err := s.db.Exec(`UPDATE leads SET next_follow_up=?, follow_ups_sent=follow_ups_sent+1, status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, nextFollowUp, status, id)
	return err
}

func (s *Store) MarkReplyReceived(id int64) error {
	_, err := s.db.Exec(`UPDATE leads SET reply_received=1, status='replied', updated_at=CURRENT_TIMESTAMP WHERE id=?`, id)
	return err
}

func (s *Store) CountLeads() (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM leads`).Scan(&count)
	return count, err
}

func (s *Store) GetStats() (map[string]int, error) {
	rows, err := s.db.Query(`SELECT status, COUNT(*) FROM leads GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	stats := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats[status] = count
	}
	return stats, rows.Err()
}

func (s *Store) EmailExists(email string) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM leads WHERE email=?`, email).Scan(&count)
	return count > 0, err
}

func (s *Store) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(sql, args...)
}

func (s *Store) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

func (s *Store) GetLeadByThreadID(threadID string) (*Lead, error) {
	row := s.db.QueryRow(`SELECT id, name, email, phone, company, website, country, source, business_type, score, status, gmail_thread_id, gmail_message_id, emails_sent, follow_ups_sent, reply_received, notes, created_at, updated_at FROM leads WHERE gmail_thread_id = ?`, threadID)
	l := &Lead{}
	var tID, mID sql.NullString
	err := row.Scan(&l.ID, &l.Name, &l.Email, &l.Phone, &l.Company, &l.Website, &l.Country, &l.Source, &l.BusinessType, &l.Score, &l.Status, &tID, &mID, &l.EmailsSent, &l.FollowUpsSent, &l.ReplyReceived, &l.Notes, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	l.GmailThreadID = tID.String
	l.GmailMessageID = mID.String
	return l, nil
}
