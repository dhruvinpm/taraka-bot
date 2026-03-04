package memory

import (
	"database/sql"
	"fmt"
)

type CompanyKB struct {
	store *Store
}

func NewCompanyKB(store *Store) *CompanyKB {
	return &CompanyKB{store: store}
}

func (kb *CompanyKB) Set(key, value string) error {
	_, err := kb.store.db.Exec(`INSERT INTO company_knowledge (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value`, key, value)
	return err
}

func (kb *CompanyKB) Get(key string) (string, error) {
	var value string
	err := kb.store.db.QueryRow(`SELECT value FROM company_knowledge WHERE key=?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (kb *CompanyKB) GetAll() (map[string]string, error) {
	rows, err := kb.store.db.Query(`SELECT key, value FROM company_knowledge`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		result[k] = v
	}
	return result, rows.Err()
}

func (kb *CompanyKB) Seed(data map[string]string) error {
	for k, v := range data {
		if err := kb.Set(k, v); err != nil {
			return fmt.Errorf("seed key %s: %w", k, err)
		}
	}
	return nil
}
