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
	GetAllUrls() ([]URLRecord, error)
}

type URLRecord struct {
	ID    int64  `json:"id"`
	URL   string `json:"url"`
	Alias string `json:"alias"`
}
