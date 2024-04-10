package main

import (
	"fmt"
	"log/slog"
	"time"
	"time-manager/internal/adapter/sqlite-repo/event"
	"time-manager/internal/config"
	"time-manager/internal/entity"
	"time-manager/internal/logging"
	"time-manager/internal/logging/sl"
	"time-manager/internal/service"
	"time-manager/internal/storage/sqlite"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
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

	log.Info("Starting time-manager", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	eventEntity := entity.Event{
		ID: 1,
		Title: "demo1",
		Time:  time.Now(),
		
	}


	eService := service.NewEventService(eventEntity, event.NewRepo(storage))

	e, err := eService.GetByName()

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println(e.Title)


	bot, err := telego.NewBot(cfg.TelegramToken, telego.WithDefaultDebugLogger())

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	updates, err := bot.UpdatesViaLongPolling(nil)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Create bot handler and specify from where to get updates
	bh, _ := th.NewBotHandler(bot, updates)
	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Send message
		_, err = bot.SendMessage(tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			"Hello %s!", update.Message.From.FirstName,
		))

		if err != nil {
			log.Error("failed to send message", sl.Err(err))
		}

	}, th.CommandEqual("start"))


	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Send message
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			"create_task, good morning!",
		))
	}, th.CommandEqual("create_task"))

	bh.Start()


	return nil
}
