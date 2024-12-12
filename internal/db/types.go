package db

type Storage interface {
	Save(shortened, original string) error
	Get(shortened string) (string, bool)
}
