package db

type Storage interface {
	Save(shortened, original string) error
}
