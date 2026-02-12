package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sqlx.DB, error) {
	const op = "database.Connect"

	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("%s: cannot connect to db: %w", op, err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)

	return db, nil
}
