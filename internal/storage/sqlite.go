package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db              *sql.DB
	userRepository  *userRepository
	questRepository *questRepository
	manager         *Manager
}

func New(storagePath string) (*Storage, error) {
	op := "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return &Storage{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return &Storage{}, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) User() *userRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &userRepository{
		storage: s,
	}

	return s.userRepository
}

func (s *Storage) Quest() *questRepository {
	if s.questRepository != nil {
		return s.questRepository
	}

	s.questRepository = &questRepository{storage: s}

	return s.questRepository
}

func (s *Storage) Manager() *Manager {
	if s.manager != nil {
		return s.manager
	}

	s.manager = &Manager{
		storage: s,
	}

	return s.manager
}
