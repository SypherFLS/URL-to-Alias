package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // init sqlite driver
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL
	)`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare(`INSERT INTO url(url, alias) VALUES(?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias = ?`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var resUrl string

	err = stmt.QueryRow(alias).Scan(&resUrl)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resUrl, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.sqlite.DeleteUrl"

	stmt, err := s.db.Prepare(`DELETE FROM url WHERE alias = ?`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AliasExists(alias string) (bool, error) {
	const op = "storage.sqlite.AliasExists"

	stmt, err := s.db.Prepare(`SELECT COUNT(*) FROM url WHERE alias = ?`)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(alias).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}
