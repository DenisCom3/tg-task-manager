package tg

import (
	"log/slog"
	"time-manager/internal/logging/sl"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func CreateTaskDescription (log *slog.Logger) func(bot *telego.Bot, update telego.Update) {
	
	return func (bot *telego.Bot, update telego.Update)  {
		_, err := bot.SendMessage(tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			"Создайте задачу, написав её ниже. Пример: 'Помыть посуду в 21:00 02.01.2006'",
		))
	
		if err != nil {
			log.Error("failed to send message", sl.Err(err))
		}
	}
}