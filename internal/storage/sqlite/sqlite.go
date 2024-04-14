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
		time TIMESTAMP NOT NULL,
		status TEXT NOT NULL,
		chat_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
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

func (s *Storage) Save(name string, time time.Time, status string, chatId int) error {
	const op = "storage.sqlite.Create"

	_, err := s.db.Exec("INSERT INTO events(name, time, status, chat_id) VALUES (?, ?, ?, ?)", name, time, status, chatId)
	if condition, ok := err.(*sqlite3.Error); ok && condition.Code == sqlite3.ErrConstraint {
		return fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
	}

	return nil

}

func (s *Storage) Update(event entity.Event) error {

	const op = "storage.sqlite.Update"

	if event.ID == nil {
		return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	res, err := s.db.Exec("UPDATE events SET name = ?, time = ?, status = ?, chat_id = ? WHERE id = ?", event.Title, event.Time, event.Status, event.ChatId, *event.ID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
    if err != nil {
        return  err
    }

    if rowsAffected == 0 {
        return storage.ErrUpdateFailed
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
	var status string
	var chatId int
	stmt.QueryRow(name).Scan(&id, &title, &time, &status, &chatId)
	
	return entity.Event{ID: &id, Title: title, Time: time, Status: status, ChatId: chatId}, nil

}

func (s *Storage) GetAllWithStatus( status string) ([]entity.Event, error) {
	const op = "storage.sqlite.GetAllWithStatus"
	stmt, err := s.db.Prepare("SELECT id, name, time, status, chat_id FROM events WHERE status = ? ORDER BY time ASC")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var entities []entity.Event = []entity.Event{}

	rows, err := stmt.Query(status)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var event entity.Event
		err = rows.Scan(&event.ID, &event.Title, &event.Time, &event.Status, &event.ChatId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		entities = append(entities, event)
	}

	return entities, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}