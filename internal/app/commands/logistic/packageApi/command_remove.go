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

	removed, err := pc.packageService.Remove(id)
	if err != nil {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, serviceErrMsg)); err != nil {
			log.Err(err).Msg("MypackageCommander.Describe: error sending reply message to chat")
		}
		log.Debug().Err(err).Msg("packageService.Remove failed")

		return
	}

	if removed {
		log.Debug().Msgf("packageService: entity ID %d removed", id)
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("packageService: entity ID %d removed", id))); err != nil {
			log.Err(err).Msg("MypackageCommander.Remove: error sending reply message to chat")
		}
	} else {
		log.Debug().Msgf("packageService: entity ID %d not removed", id)
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("packageService: entity ID %d not removed", id))); err != nil {
			log.Err(err).Msg("MypackageCommander.Remove: error sending reply message to chat")
		}

	}
}
