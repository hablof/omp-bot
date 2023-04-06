package packageApi

import (
	"fmt"
	"strconv"

	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hablof/omp-bot/internal/model/logistic"
	"github.com/rs/zerolog/log"
)

// Update implements PackageCommander
func (pc *MypackageCommander) Update(inputMsg *tgbotapi.Message) {

	args := strings.Split(inputMsg.CommandArguments(), ";")

	id, err := strconv.ParseUint(strings.TrimSpace(args[0]), 10, 64)
	if err != nil {
		log.Debug().Err(err).Msgf("MypackageCommander.Update: cannot parse ID (int) from command argument: %s", args[0])
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("идентификатор не может быть \"%s\"", args[0]))); err != nil {
			log.Debug().Err(err).Msg("MypackageCommander.Update: error sending reply message to chat")
		}

		return
	}

	editArgMap := make(map[string]string, logistic.PackageFieldsCount)

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, logistic.Title):
			editArgMap[logistic.Title] = strings.TrimSpace(strings.TrimPrefix(arg, logistic.Title))

		case strings.HasPrefix(arg, logistic.Material):
			editArgMap[logistic.Material] = strings.TrimSpace(strings.TrimPrefix(arg, logistic.Material))

		case strings.HasPrefix(arg, logistic.MaximumVolume):
			editArgMap[logistic.MaximumVolume] = strings.TrimSpace(strings.TrimPrefix(arg, logistic.MaximumVolume))

		case strings.HasPrefix(arg, logistic.Reusable):
			editArgMap[logistic.Reusable] = strings.TrimSpace(strings.TrimPrefix(arg, logistic.Reusable))

		default:
			log.Debug().Msgf("MypackageCommander.Update: found argument: %s", arg)
			if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("Неизвестный аргумент: \"%s\"", arg))); err != nil {
				log.Debug().Err(err).Msg("MypackageCommander.Update: error sending reply message to chat")
			}

			return
		}

	}

	isUpdated, err := pc.packageService.Update(id, editArgMap)
	if err != nil {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "bad request")); err != nil {
			log.Printf("MypackageCommander.Update: error sending reply message to chat - %v", err)
		}

		return
	}

	if isUpdated {
		log.Debug().Msgf("MypackageCommander.Update: package id %d updated", id)
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("Упаковка ID=%d успешно отредактирована", id))); err != nil {
			log.Printf("MypackageCommander.Update: error sending reply message to chat - %v", err)
		}
	} else {
		log.Debug().Msgf("MypackageCommander.Update: package id %d NOT updated", id)
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("Упаковка ID=%d НЕ отредактирована", id))); err != nil {
			log.Printf("MypackageCommander.Update: error sending reply message to chat - %v", err)
		}
	}
}
