package db

import (
	"sync"
)

type InMemoryDb struct {
	store map[string]string
	mu    sync.Mutex // a MUTEX is needed to keep the store updates threadsafe
}

func NewInMemoryDb() Storage {
	return &InMemoryDb{
		store: make(map[string]string),
	}
}

func (db *InMemoryDb) Save(shortened, original string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.store[shortened] = original
	return nil
}

func (db *InMemoryDb) Get(shortened string) (string, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	original, exists := db.store[shortened]
	return original, exists
}
