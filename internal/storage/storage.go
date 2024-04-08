package storage

import (
	"database/sql"
	"errors"
)

type Storage struct {
	db *sql.DB
}

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)