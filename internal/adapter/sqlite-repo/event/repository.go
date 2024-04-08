package event

import (
	"time"
)

type Repository struct {
	store Storage
}

type Storage interface {
	Save(name string, time time.Time) error
}

func NewRepo(s Storage) Repository {
	return Repository{
		store: s,
	}
}