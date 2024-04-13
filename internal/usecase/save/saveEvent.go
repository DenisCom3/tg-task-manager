package save

type EventService interface {
	Save() error
}

type SaveEvent struct {
	EventService
}