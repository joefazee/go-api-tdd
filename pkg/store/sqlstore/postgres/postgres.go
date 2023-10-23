package postgres

import "database/sql"

type postgresStore struct {
	db *sql.DB
}

// NewPostgresStore create new instance of postgresStore
func NewPostgresStore(db *sql.DB) *postgresStore {
	return &postgresStore{
		db: db,
	}
}
