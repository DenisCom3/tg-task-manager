package tg

import (
	"fmt"
	"log/slog"
	"time-manager/internal/adapter/sqlite-repo/event"
	"time-manager/internal/logging/sl"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func GetTask (log *slog.Logger, repo event.Repository) func(bot *telego.Bot, update telego.Update) {

	return func (bot *telego.Bot, update telego.Update)  {
		event, err := repo.GetByName(update.Message.Text)
		if err != nil {
			log.Error("failed to get event", sl.Err(err))
		}

		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("%s в %s", event.Title, event.Time.Format("15:04 02.01.2006")),
		))
	}
}