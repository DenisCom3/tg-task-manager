package tg

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
	"time-manager/internal/adapter/sqlite-repo/event"
	"time-manager/internal/broadcaster"
	"time-manager/internal/entity"
	"time-manager/internal/logging/sl"
	"time-manager/internal/service"
	"time-manager/internal/usecase/notify"
	"time-manager/internal/usecase/save"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func CreateTask (log *slog.Logger, repo event.Repository) func(bot *telego.Bot, update telego.Update) {


	var task entity.Event
	
	
	return func (bot *telego.Bot, update telego.Update)  {
		messageText := update.Message.Text
		arr := strings.Split(messageText, "время")
		
		taskName := strings.Trim(arr[0], " ")
		taskTime, err := time.ParseInLocation("15:04 02.01.2006", strings.Trim(arr[1], 	" "), time.Local)
		if err != nil {
			log.Info("failed to parse time", slog.String("time", strings.Trim(arr[1], " ")))
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
			Status: "pending",
			ChatId: int(update.Message.Chat.ID),
		}

		eService := service.NewEventService(task, repo)

		useCase := save.SaveEvent{
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

		broadcaster.CheckEventTimeAndHandle(task, bot, notify.New(eService))
	}
}

