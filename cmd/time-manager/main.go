package main

import (
	"fmt"


	"time-manager/internal/adapter/sqlite-repo/event"
	"time-manager/internal/broadcaster"
	"time-manager/internal/service"

	"time-manager/internal/config"
	"time-manager/internal/entity"

	"time-manager/internal/handler/tg"
	"time-manager/internal/logging"
	"time-manager/internal/storage/sqlite"

	"github.com/joho/godotenv"
	th "github.com/mymmrac/telego/telegohandler"


)

func main() {

	err := run()

	if err != nil {
		panic(err)
	}
}


func run() error {
	
	err := godotenv.Load()

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	cfg := config.MustLoad()

	
	log, err := logging.Setup(cfg.Env)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	storage, err := sqlite.New(cfg.StoragePath)
	defer func(storage *sqlite.Storage) error {
		err := storage.Close()
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		return nil
	}(storage)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	repo := event.NewRepo(storage)

	bot, bh, err := tg.Init(cfg.TelegramToken, repo, log)
	defer bot.StopLongPolling()


	if err != nil {
		return fmt.Errorf("%w", err)
	}
	service := service.NewEventService(entity.Event{}, repo)

	go broadcaster.Start(service, bot, log)

	bh.Handle(tg.Start(log), th.CommandEqual("start"))
	bh.Handle(tg.CreateTaskDescription(log), th.CommandEqual("create_task"))
	bh.Handle(tg.CreateTask(log, repo), th.AnyMessage())
	fmt.Println("started")
	bh.Start()
	defer bh.Stop()


	
	return nil
}