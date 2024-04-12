package tg

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
	"time-manager/internal/adapter/sqlite-repo/event"
	"time-manager/internal/entity"
	"time-manager/internal/logging/sl"
	"time-manager/internal/service"
	"time-manager/internal/usecase"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func CreateTask (log *slog.Logger, repo event.Repository) func(bot *telego.Bot, update telego.Update) {


	var task entity.Event
	
	
	return func (bot *telego.Bot, update telego.Update)  {
		messageText := update.Message.Text
		arr := strings.Split(messageText, "в")
		
		taskName := strings.Trim(arr[0], " ")
		taskTime, err := time.Parse("15:04 02.01.2006", strings.Trim(arr[1], " "))
		if err != nil {
			log.Error("failed to parse time", sl.Err(err))
			
			_, _ = bot.SendMessage(tu.Message(
				tu.ID(update.Message.Chat.ID),
				"Не смог распознать время",
			))
			return
		}

		if taskTime.Before(time.Now()) {

			if taskTime.Year() == 2007 {
				_, _ = bot.SendMessage(tu.Message(
					tu.ID(update.Message.Chat.ID),
					"Никто, никогда не вернёт 2007!",
				))
				return
			}

			_, _ = bot.SendMessage(tu.Message(
				tu.ID(update.Message.Chat.ID),
				"В прошлое не вернуться!",
			))
			return
		}

		task = entity.Event{
			Title: taskName,
			Time:  taskTime,
		}

		eService := service.NewEventService(task, repo)

		useCase := usecase.SaveEvent{
			EventService: eService,
		}

		err = useCase.Save()

		if err != nil {
			log.Error("failed to save event", sl.Err(err))
			
			_, _ = bot.SendMessage(tu.Message(
				tu.ID(update.Message.Chat.ID),
				"Не смог сохранить задачу",
			))
			return
		}

		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("Создал задачу %s в %s", taskName, taskTime.Format("15:04 02.01.2006")),
		))
	}
}