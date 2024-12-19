package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/khaled4vokalz/gourl_shortener/internal/common"
	"github.com/khaled4vokalz/gourl_shortener/internal/config"
	logger "github.com/khaled4vokalz/gourl_shortener/internal/logging"
)

type InMemoryDb struct {
	log   logger.Logger
	store map[string]*ShortenedUrl
	mu    sync.Mutex // a MUTEX is needed to keep the store updates threadsafe
}

type ShortenedUrl struct {
	URL       string
	ExpiresAt time.Time
}

func NewInMemoryDb() (*InMemoryDb, error) {
	return &InMemoryDb{
		store: make(map[string]*ShortenedUrl),
		log:   logger.GetLogger(),
	}, nil
}

func (db *InMemoryDb) Save(shortened, original string, expiresAt time.Time) error {
	expires := config.GetConfig().UrlsExpiresAt
	if !expiresAt.IsZero() {
		expires = expiresAt
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	db.store[shortened] = &ShortenedUrl{
		URL:       original,
		ExpiresAt: expires,
	}
	return nil
}

func (db *InMemoryDb) Get(shortened string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	shortenedObj, exists := db.store[shortened]
	if !exists {
		logger.GetLogger().Debug(fmt.Sprintf("No url found for key '%s'", shortened))
		return "", common.NotFound
	}
	if !shortenedObj.ExpiresAt.IsZero() && time.Now().After(shortenedObj.ExpiresAt) {
		logger.GetLogger().Debug(fmt.Sprintf("URL '%s' for key '%s' has expired", shortenedObj.URL, shortened))
		return "", common.Expired
	}
	return shortenedObj.URL, nil
}

func (db *InMemoryDb) IsAlive() bool {
	return true
}
