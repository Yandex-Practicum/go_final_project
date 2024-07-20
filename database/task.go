package database

import "github.com/jmoiron/sqlx"

var _ Task = (*TaskRepo)(nil)

type TaskRepo struct {
	Db *sqlx.DB
}
