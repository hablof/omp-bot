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

	// args[1:] no info about id
	if err, ok := pc.fillArgMap(args[1:], editArgMap).(*ErrBadArgument); ok { // дурно пахнет

		log.Debug().Msgf("MypackageCommander.Update: unknown argument: %s", err.argument)
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("Некорректный аргумент: \"%s\"", err.argument))); err != nil {
			log.Debug().Err(err).Msg("MypackageCommander.Update: error sending reply message to chat")
		}

		return
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

func (pc *MypackageCommander) fillArgMap(args []string, argMap map[string]string) error {

	for _, arg := range args {

		kv := strings.SplitN(arg, "=", 2)
		if len(kv) < 2 {
			return &ErrBadArgument{arg}
		}

		switch key := strings.ToLower(strings.TrimSpace(kv[0])); key {
		case logistic.Title:
			argMap[logistic.Title] = strings.TrimSpace(kv[1])

		case logistic.Material:
			argMap[logistic.Material] = strings.TrimSpace(kv[1])

		case logistic.MaximumVolume:
			argMap[logistic.MaximumVolume] = strings.TrimSpace(kv[1])

		case logistic.Reusable:
			argMap[logistic.Reusable] = strings.TrimSpace(kv[1])

		default:
			return &ErrBadArgument{arg}
		}
	}

	return nil
}
