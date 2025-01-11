package files

import (
	"database/sql"
	_ "embed"
)

// PSQLProjectStore implements ProjectStore interface using an SQL database
type PSQLFilestStore struct {
	db *sql.DB
}

// NewSQLProjectStore creates a new SQLProjectStore
func NewPSQLFilestStore(db *sql.DB) *PSQLFilestStore {
	return &PSQLFilestStore{db: db}
}

func (store *PSQLFilestStore) Upload(project string, file *File, data []byte) error {
	return nil
}
