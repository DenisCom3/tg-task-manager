package usecase

type EventService interface {
	Save() error
}

type SaveEvent struct {
	EventService
}