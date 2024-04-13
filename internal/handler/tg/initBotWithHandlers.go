package tg

import (
	"fmt"
	"log/slog"
	"time-manager/internal/adapter/sqlite-repo/event"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

)

func InitBotWithHandlers (token string, repo event.Repository, log *slog.Logger) (*telego.Bot, error) {
	bot, err := telego.NewBot(token, telego.WithDefaultDebugLogger())

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	updates, err := bot.UpdatesViaLongPolling(nil)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// Create bot handler and specify from where to get updates
	bh, _ := th.NewBotHandler(bot, updates)
	defer bot.StopLongPolling()
	defer bh.Stop()


	bh.Handle(Start(log), th.CommandEqual("start"))
	bh.Handle(CreateTaskDescription(log), th.CommandEqual("create_task"))
	bh.Handle(CreateTask(log, repo), th.AnyMessage())
	// bh.Handle(tg.GetTask(log, repo), th.AnyMessage())

	bh.Start()

	return bot, nil
}