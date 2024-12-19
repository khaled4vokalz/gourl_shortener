package db

import "time"

type Storage interface {
	Save(shortened, original string, expiredAt time.Time) error
	Get(shortened string) (string, error)
	IsAlive() bool
}
