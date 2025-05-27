package factory

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

func New(storagePath string) (*Storage, error) {

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) Close() error {
	return s.DB.Close()
}

func (s *Storage) Insert(response string, args ...any) error {
	_, err := s.DB.Exec(response, args...)
	if err != nil {
		return fmt.Errorf("insert failed: %w", err)

	}
	return nil
}
