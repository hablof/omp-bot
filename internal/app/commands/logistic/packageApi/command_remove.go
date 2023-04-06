package packageApi

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
)

// Remove implements PackageCommander
func (pc *MypackageCommander) Remove(inputMsg *tgbotapi.Message) {

	argument := inputMsg.CommandArguments()

	id, err := strconv.ParseUint(argument, 10, 64)
	if err != nil {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, badRequestMsg)); err != nil {
			log.Debug().Err(err).Msg("MypackageCommander.Remove: error sending reply message to chat")
		}
		log.Debug().Err(err).Msgf("MypackageCommander.Remove: cannot parse int from command argument: %s", argument)

		return
	}

	isRemoved, err := pc.packageService.Remove(id)
	if err != nil {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, serviceErrMsg)); err != nil {
			log.Err(err).Msg("MypackageCommander.Describe: error sending reply message to chat")
		}
		log.Debug().Err(err).Msg("packageService.Remove failed")

		return
	}

	if isRemoved {
		log.Debug().Msgf("packageService: entity ID %d removed", id)
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("Упаковка ID=%d успешно удалена", id))); err != nil {
			log.Err(err).Msg("MypackageCommander.Remove: error sending reply message to chat")
		}
	} else {
		log.Debug().Msgf("packageService: entity ID %d NOT removed", id)
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("Упаковка ID=%d НЕ удалена", id))); err != nil {
			log.Err(err).Msg("MypackageCommander.Remove: error sending reply message to chat")
		}

	}
}
