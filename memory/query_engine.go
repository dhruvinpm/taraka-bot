package memory

import (
	"database/sql"
	"fmt"
)

type QueryEngine struct {
	store *Store
}

func NewQueryEngine(store *Store) *QueryEngine {
	return &QueryEngine{store: store}
}

func (qe *QueryEngine) Execute(query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := qe.store.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make(map[string]interface{})
		for i, col := range cols {
			val := vals[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

func (qe *QueryEngine) ExecuteOne(query string, args ...interface{}) (map[string]interface{}, error) {
	results, err := qe.Execute(query, args...)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}
	return results[0], nil
}
