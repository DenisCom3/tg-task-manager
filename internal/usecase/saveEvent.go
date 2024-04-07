package usecase

type EventService interface {
	save()
}

type SaveEvent struct {
	EventService
}