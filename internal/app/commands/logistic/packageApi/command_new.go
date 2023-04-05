package packageApi

import (
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"

	"github.com/hablof/omp-bot/internal/model/logistic"
)

// New implements PackageCommander
func (pc *MypackageCommander) New(inputMsg *tgbotapi.Message) {
	args := strings.Split(inputMsg.CommandArguments(), ";")

	// количестпо полей, не считая поле ID
	if len(args) != logistic.PackageFieldsCount-1 {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "неверное количество аргументов")); err != nil {
			log.Debug().Err(err).Msg("MypackageCommander.New: error sending reply message to chat")
		}
		log.Debug().Msg("MypackageCommander.New: wrong args count")

		return
	}

	createArgMap := make(map[string]string, logistic.PackageFieldsCount-1)

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, logistic.Title):
			createArgMap[logistic.Title] = strings.TrimSpace(strings.TrimPrefix(arg, logistic.Title))

		case strings.HasPrefix(arg, logistic.Material):
			createArgMap[logistic.Material] = strings.TrimSpace(strings.TrimPrefix(arg, logistic.Material))

		case strings.HasPrefix(arg, logistic.MaximumVolume):
			createArgMap[logistic.MaximumVolume] = strings.TrimSpace(strings.TrimPrefix(arg, logistic.MaximumVolume))

		case strings.HasPrefix(arg, logistic.Reusable):
			createArgMap[logistic.Reusable] = strings.TrimSpace(strings.TrimPrefix(arg, logistic.Reusable))

		default:
			log.Debug().Msgf("MypackageCommander.Edit: unknown argument: %s", arg)
			pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("Неизвестный аргумент: \"%s\"", arg)))

			return
		}

	}

	id, err := pc.packageService.Create(createArgMap)
	switch {
	case errors.Is(err, ErrBadRequest):
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "некорректный запрос")); err != nil {
			log.Debug().Err(err).Msg("MypackageCommander.New: error sending reply message to chat")
		}
		return

	case err != nil:
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "🤡🤡🤡 Ошибка сервиса 🤡🤡🤡")); err != nil {
			log.Debug().Err(err).Msg("MypackageCommander.New: error sending reply message to chat")
		}
		return
	}

	log.Debug().Msgf("MypackageCommander.New: package id %d created", id)

	if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("New package id: %d", id))); err != nil {
		log.Debug().Err(err).Msg("MypackageCommander.New: error sending reply message to chat")
	}
}
