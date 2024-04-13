package notify

import (
	"fmt"
	"sync"
	"time"
	"time-manager/internal/entity"
)

type EventService interface {
	GetAllPendingTasks() ([]entity.Event, error)
}

type NotifyEvent struct {
	EventService
}

func New (eventService EventService) NotifyEvent {
	return NotifyEvent{
		EventService: eventService,
	}
}

func (n NotifyEvent) BeforeHour(e entity.Event) {
	
	var wg sync.WaitGroup
	wg.Add(1)
	notifyTime := e.Time.Add(-1 * time.Hour)

	fmt.Println("before hour", notifyTime)
	for {
		fmt.Println("now",time.Now(), "event", e.Title, "|||", notifyTime, "|||", time.Now().After(notifyTime))
		if time.Now().After(notifyTime) {
			fmt.Println("event triggered 1 hour", e.Title)
			wg.Done()
			break
		}
		fmt.Println("wait 1 min...")
		time.Sleep(1 * time.Minute)
	}

	wg.Wait()
}

func (n NotifyEvent) Before15Min(e entity.Event) {

	var wg sync.WaitGroup
	wg.Add(1)
	notifyTime := e.Time.Add(-15 * time.Minute)

	for {
		if time.Now().After(notifyTime) {
			fmt.Println("event triggered 15 min", e.Title)
			wg.Done()
			break
		}
		fmt.Println("wait 1 min...")
		time.Sleep(1 * time.Minute)
	}

	wg.Wait()
}