package tg

import (
	"log/slog"
	"time-manager/internal/logging/sl"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)



func Start (log *slog.Logger) func(bot *telego.Bot, update telego.Update) {
	
	return func (bot *telego.Bot, update telego.Update)  {
		_, err := bot.SendMessage(tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			"Hello %s!", update.Message.From.FirstName,
		))
	
		if err != nil {
			log.Error("failed to send message", sl.Err(err))
		}
	}

}