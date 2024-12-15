package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // needed driver for go to be able to work with PostgresDb
)

type PostgresDb struct {
	db *sql.DB
}

func NewPostgresDb(connectionString string) (Storage, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %v", err)
	}

	return &PostgresDb{db: db}, nil
}

func (db *PostgresDb) Save(shortened, original string) error {
	_, err := db.db.Exec("INSERT INTO urls (shortened, original) VALUES ($1, $2)", shortened, original)
	return err
}

func (db *PostgresDb) Get(shortened string) (string, bool) {
	var original string
	err := db.db.QueryRow("SELECT original FROM urls WHERE shortened = $1", shortened).Scan(&original)
	if err == sql.ErrNoRows {
		return "", false
	}
	if err != nil {
		return "", false
	}
	return original, true
}