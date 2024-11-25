package repository

import "database/sql"

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
