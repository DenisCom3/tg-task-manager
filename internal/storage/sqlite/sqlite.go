package sqlite

import (
	"database/sql"
	"fmt"
	"time"
	"time-manager/internal/entity"
	"time-manager/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, op)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS events(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		time TIMESTAMP NOT NULL
        );
	CREATE INDEX IF NOT EXISTS idx_name ON events(name);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(name string, time time.Time) error {
	const op = "storage.sqlite.Create"

	_, err := s.db.Exec("INSERT INTO events(name, time) VALUES (?, ?)", name, time)
	if condition, ok := err.(*sqlite3.Error); ok && condition.Code == sqlite3.ErrConstraint {
		return fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
	}

	return nil

}

func (s *Storage) GetByName(name string) (entity.Event, error) {
	const op = "storage.sqlite.GetByName"

	stmt, err := s.db.Prepare("SELECT id, name, time FROM events WHERE name = ?")
	if err != nil {
		return entity.Event{}, fmt.Errorf("%s: %w", op, err)
	}
	var id int
	var title string
	var time time.Time
	stmt.QueryRow(name).Scan(&id, &title, &time)
	
	return entity.Event{ID: id, Title: title, Time: time}, nil

}