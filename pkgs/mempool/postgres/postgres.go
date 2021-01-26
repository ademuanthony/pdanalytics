package postgres

import (
	"database/sql"
)

type PgDb struct {
	db *sql.DB
}
