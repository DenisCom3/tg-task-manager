package service

import (
	"time-manager/internal/adapter/sqlite-repo/event"
	"time-manager/internal/entity"
)

type EventService struct {
	Event *entity.Event
	repo  event.Repository
}

func NewEventService(event entity.Event, repo event.Repository) EventService {

	return EventService{
		Event: &event,
		repo:  repo,
	}
}

func (e EventService) Save() error {
	err := e.repo.Save(*e.Event)

	if err != nil {
		return err
	}

	return nil
}

func (e EventService) GetByName() (entity.Event, error) {

	event, err := e.repo.GetByName(e.Event.Title)
	if err != nil {
		return entity.Event{}, err
	}

	return event, nil
}

func (e EventService) Update(event entity.Event) error {

	err := e.repo.Update(event)
	if err != nil {
		return err
	}
	return nil
	
}

func (e EventService) GetAllPendingTasks() ([]entity.Event, error) {

	events, err := e.repo.GetAllWithStatus("pending")
	if err != nil {
		return []entity.Event{}, err
	}	
	return events, nil
}
