package tg

import (
	"fmt"
	"log/slog"
	"time-manager/internal/adapter/sqlite-repo/event"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

)

func Init(token string, repo event.Repository, log *slog.Logger) (*telego.Bot, *th.BotHandler, error) {
	bot, err := telego.NewBot(token, telego.WithDefaultDebugLogger())

	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}
	bh, err := th.NewBotHandler(bot, updates)

	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	return bot, bh, nil
}