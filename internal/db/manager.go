package db

import (
	"errors"
	"fmt"

	"github.com/khaled4vokalz/gourl_shortener/internal/config"
)

var (
	supportedStorageOptions = [...]string{"in-memory", "postgres"}
	InvalidStorageOption    = errors.New(fmt.Sprintf("Invalid storage option provided, supported options are %v", supportedStorageOptions))
)

func GetDb(conf config.StorageConfig) (Storage, error) {
	if conf.Type == "in-memory" {
		return NewInMemoryDb(), nil
	} else if conf.Type == "postgres" {
		return NewPostgresDb(conf.Db_Conn_String)
	} else {
		return nil, InvalidStorageOption
	}
}
