package broadcaster

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
	"time-manager/internal/logging/sl"
	"time-manager/internal/usecase/notify"
)

func Start(s notify.EventService, log *slog.Logger) error {
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

				if event.Time.Add(-1 * time.Hour).After(time.Now()) && time.Now().Add(2 * time.Hour).After(event.Time) {
					fmt.Println("event triggered", event.Title)
					notify.BeforeHour(event)
				}

				if event.Time.Add(-15 * time.Minute).After(time.Now()) && time.Now().Add(1 * time.Hour).After(event.Time) {
					fmt.Println("event triggered", event.Title)
					notify.Before15Min(event)
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