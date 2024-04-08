package service

import (
	"time-manager/internal/adapter/sqlite-repo/event"
	"time-manager/internal/entity"
)

type EventService struct {
	event entity.Event
	repo  event.Repository
}

func NewEventService(event entity.Event, repo event.Repository) EventService {

	return EventService{
		event: event,
		repo:  repo,
	}
}

func (e EventService) Save() error {
	err := e.repo.Save(e.event)

	if err != nil {
		return err
	}

	return nil
}