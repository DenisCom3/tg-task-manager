package broadcaster

import (
	"context"
	"fmt"
	"log/slog"
	// "sync"
	"time"
	"time-manager/internal/entity"
	"time-manager/internal/logging/sl"
	"time-manager/internal/usecase/notify"

	"github.com/mymmrac/telego"
)

func Start(ctx context.Context, s notify.EventService, bot *telego.Bot, log *slog.Logger) error {
	ticker := time.NewTicker(1 * time.Hour)
	errChan := make(chan error)
	notify := notify.New(s)
	go func() {
		CheckEvents(s, bot, log, notify, errChan)
		tickerLoop:
		for {
			select {
				case <-ticker.C:
					CheckEvents(s, bot, log, notify, errChan)
				case <-ctx.Done():
					break tickerLoop
				default:
					continue
			}
		}
	}()


	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
	
}

func CheckEvents( s notify.EventService, bot *telego.Bot, log *slog.Logger, notify notify.NotifyEvent, errChan chan error) {
	
	events, err := notify.EventService.GetAllPendingTasks()
	if err != nil {
		log.Error("failed to get all pending tasks", sl.Err(err))
		errChan <- err
	}

	for _, event := range events {			
		CheckEventTimeAndHandle(event, bot, notify)
	}
}


func CheckEventTimeAndHandle(event entity.Event, bot *telego.Bot, notify notify.NotifyEvent) {
	isBeforeHour := event.Time.Add(-1 * time.Hour).After(time.Now()) && time.Now().Add(2 * time.Hour).After(event.Time)
	isBefore15Min := event.Time.Add(-15 * time.Minute).After(time.Now()) && time.Now().Add(1 * time.Hour).After(event.Time)
	if isBeforeHour  {
		fmt.Println("event triggered", event.Title)
		go notify.BeforeHour(event, bot)
	}
	if isBefore15Min {
		fmt.Println("event triggered", event.Title)
		go notify.Before15Min(event, bot)
	}
}