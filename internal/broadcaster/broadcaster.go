package broadcaster

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
	"time-manager/internal/logging/sl"
	"time-manager/internal/usecase/notify"

	"github.com/mymmrac/telego"
)

func Start(s notify.EventService, bot *telego.Bot, log *slog.Logger) error {
	var wg sync.WaitGroup
	errChan := make(chan error)
	notify := notify.New(s)
	wg.Add(1)
	go func() {
		for {
			events, err := notify.EventService.GetAllPendingTasks()

			if err != nil {
				log.Error("failed to get all pending tasks", sl.Err(err))
				errChan <- err
				break
			}

			for _, event := range events {
				fmt.Println("now",time.Now(), "event", event.Title, event.Time)
				
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
			time.Sleep(1 * time.Hour)
		}
		defer wg.Done()
	}()

	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
	
}