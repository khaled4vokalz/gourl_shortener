package db

import (
	"fmt"
	"sync"

	"github.com/khaled4vokalz/gourl_shortener/internal/common"
	logger "github.com/khaled4vokalz/gourl_shortener/internal/logging"
)

type InMemoryDb struct {
	log   logger.Logger
	store map[string]string
	mu    sync.Mutex // a MUTEX is needed to keep the store updates threadsafe
}

func NewInMemoryDb() (*InMemoryDb, error) {
	return &InMemoryDb{
		store: make(map[string]string),
		log:   logger.GetLogger(),
	}, nil
}

func (db *InMemoryDb) Save(shortened, original string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.store[shortened] = original
	return nil
}

func (db *InMemoryDb) Get(shortened string) (string, error) {
	db.mu.Lock()
	var err error
	defer db.mu.Unlock()
	original, exists := db.store[shortened]
	if !exists {
		logger.GetLogger().Debug(fmt.Sprintf("No url found for key '%s'", shortened))
		err = common.NotFound
	}
	return original, err
}
