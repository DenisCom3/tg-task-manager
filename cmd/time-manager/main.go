package main

import (
	"fmt"
	"time-manager/internal/adapter/sqlite-repo/event"
	"time-manager/internal/config"
	"time-manager/internal/handler/tg"
	"time-manager/internal/logging"
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

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	repo := event.NewRepo(storage)

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
	defer bot.StopLongPolling()
	defer bh.Stop()


	bh.Handle(tg.Start(log), th.CommandEqual("start"))


	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Send message
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Создайте задачу, написав её ниже. Пример: 'Помыть посуду в 21:00 02.01.2006'",
		))
	}, th.CommandEqual("create_task"))

	bh.Handle(tg.CreateTask(log, repo), th.AnyMessage())

	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		// Send message
		_, _ = bot.SendMessage(tu.Message(tu.ID(query.Message.GetChat().ID), "GO GO GO"))

		// Answer callback query
		_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID).WithText("Done"))
	}, th.AnyCallbackQueryWithMessage(), th.CallbackDataEqual("go"))

	bh.Start()

	return nil
}
