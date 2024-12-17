package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/khaled4vokalz/gourl_shortener/internal/common"
	"github.com/khaled4vokalz/gourl_shortener/internal/config"
	logger "github.com/khaled4vokalz/gourl_shortener/internal/logging"
	_ "github.com/lib/pq" // needed driver for go to be able to work with PostgresDb
)

type PostgresDb struct {
	db *sql.DB
}

func NewPostgresDb(connectionString string) (*PostgresDb, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %v", err)
	}

	return &PostgresDb{db: db}, nil
}

func (db *PostgresDb) Save(shortened, original string, expiresAt time.Time) error {
	expires := config.GetConfig().UrlsExpiresAt
	if !expiresAt.IsZero() {
		expires = expiresAt
	}
	_, err := db.db.Exec("INSERT INTO urls (shortened, original, expires_at) VALUES ($1, $2, $3)", shortened, original, expires)
	return err
}

func (db *PostgresDb) Get(shortened string) (string, error) {
	var original string
	var expiresAt *time.Time
	err := db.db.QueryRow("SELECT original, expires_at FROM urls WHERE shortened = $1", shortened).Scan(&original, &expiresAt)
	if err == sql.ErrNoRows {
		logger.GetLogger().Debug(fmt.Sprintf("No url found for key '%s'", shortened))
		return "", common.NotFound
	}
	if expiresAt != nil && time.Now().After(*expiresAt) {
		logger.GetLogger().Debug(fmt.Sprintf("URL '%s' for key '%s' has expired", original, shortened))
		return "", common.Expired
	}
	return original, err
}

func (db *PostgresDb) IsAlive() bool {
	if err := db.db.Ping(); err != nil {
		return false
	}
	return true
}
