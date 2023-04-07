package packageApi

import (
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"

	"github.com/hablof/omp-bot/internal/model/logistic"
)

// Create implements PackageCommander
func (pc *MypackageCommander) Create(inputMsg *tgbotapi.Message) {
	args := strings.Split(inputMsg.CommandArguments(), ";")

	// количестпо полей, не считая поле ID
	if len(args) != logistic.PackageFieldsCount-1 {
		log.Debug().Msg("MypackageCommander.Create: wrong args count")
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "неверное количество аргументов")); err != nil {
			log.Debug().Err(err).Msg("MypackageCommander.Create: error sending reply message to chat")
		}

		return
	}

	createArgMap := make(map[string]string, logistic.PackageFieldsCount-1)

	if err, ok := pc.fillArgMap(args, createArgMap).(*ErrBadArgument); ok { // дурно пахнет

		log.Debug().Msgf("MypackageCommander.Update: unknown argument: %s", err.argument)
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("Некорректный аргумент: \"%s\"", err.argument))); err != nil {
			log.Debug().Err(err).Msg("MypackageCommander.Update: error sending reply message to chat")
		}

		return
	}

	id, err := pc.packageService.Create(createArgMap)
	switch {
	case errors.Is(err, ErrBadRequest):
		log.Debug().Err(err).Msg("packageService.Create failed")
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, badRequestMsg)); err != nil {
			log.Debug().Err(err).Msg("MypackageCommander.Create: error sending reply message to chat")
		}
		return

	case err != nil:
		log.Debug().Err(err).Msg("packageService.Create failed")
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, serviceErrMsg)); err != nil {
			log.Debug().Err(err).Msg("MypackageCommander.Create: error sending reply message to chat")
		}
		return
	}

	log.Debug().Msgf("MypackageCommander.Create: package id %d created", id)

	if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("New package id: %d", id))); err != nil {
		log.Debug().Err(err).Msg("MypackageCommander.Create: error sending reply message to chat")
	}
}
