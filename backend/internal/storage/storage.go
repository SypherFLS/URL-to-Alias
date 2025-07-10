package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

type Storage interface {
	SaveUrl(urlToSave string, alias string) (int64, error)
	GetUrl(alias string) (string, error)
	DeleteUrl(alias string) error
	AliasExists(alias string) (bool, error)
}
