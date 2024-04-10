package event

import (
	"time"
	"time-manager/internal/entity"
)

type Repository struct {
	store Storage
}

type Storage interface {
	Save(name string, time time.Time) error
	GetByName(name string) (entity.Event, error)
}

func NewRepo(s Storage) Repository {
	return Repository{
		store: s,
	}
}