package event

import (
	"time"
	"time-manager/internal/entity"
)

type Repository struct {
	store Storage
}

type Storage interface {
	Save(name string, time time.Time, status string, chatId int) error
	GetByName(name string) (entity.Event, error)
	GetAllWithStatus(status string) ([]entity.Event, error)
	Update(event entity.Event) error
}

func NewRepo(s Storage) Repository {
	return Repository{
		store: s,
	}
}