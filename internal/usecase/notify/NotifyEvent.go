package notify

import (
	"fmt"
	"sync"
	"time"
	"time-manager/internal/entity"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/mymmrac/telego"
)

type EventService interface {
	GetAllPendingTasks() ([]entity.Event, error)
	Update(entity.Event) error
}

type NotifyEvent struct {
	EventService
}

func New (eventService EventService) NotifyEvent {
	return NotifyEvent{
		EventService: eventService,
	}
}

func (n NotifyEvent) BeforeHour(e entity.Event, bot *telego.Bot) {
	
	var wg sync.WaitGroup
	wg.Add(1)
	notifyTime := e.Time.Add(-1 * time.Hour)

	fmt.Println("before hour", notifyTime)
	for {
		fmt.Println("now",time.Now(), "event", e.Title, "|||", notifyTime, "|||", time.Now().After(notifyTime))
		if time.Now().After(notifyTime) {
			
			_, _ = bot.SendMessage(
				tu.Message(
					tu.ID(int64(e.ChatId)),
					fmt.Sprintf("До события '%s' остался 1 час", e.Title),
				),
			)
			wg.Done()
			break
		}
		fmt.Println("wait 1 min...")
		time.Sleep(1 * time.Minute)
	}

	wg.Wait()
}

func (n NotifyEvent) Before15Min(e entity.Event, bot *telego.Bot) {

	var wg sync.WaitGroup
	wg.Add(1)
	notifyTime := e.Time.Add(-15 * time.Minute)

	for {
		if time.Now().After(notifyTime) {
			fmt.Println("event triggered 15 min", e.Title)
			_, _ = bot.SendMessage(
				tu.Message(
					tu.ID(int64(e.ChatId)),
					fmt.Sprintf("До события '%s' осталось 15 минут", e.Title),
				),
			)
			e.Status = "sent"
			n.Update(e)
			wg.Done()
			break
		}
		fmt.Println("wait 1 min...")
		time.Sleep(1 * time.Minute)
	}

	wg.Wait()
}