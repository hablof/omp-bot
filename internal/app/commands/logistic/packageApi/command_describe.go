package packageApi

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
)

const template = `Название: %s
Материал: %s
Максимальный объём: %3.0f %s
Многоразовая: %s`

// Describe implements PackageCommander
func (pc *MypackageCommander) Describe(inputMsg *tgbotapi.Message) {

	argument := inputMsg.CommandArguments()

	id, err := strconv.ParseUint(argument, 10, 64)
	if err != nil {
		log.Err(err).Msgf("MypackageCommander.Describe: cannot parse int from command argument: %s", argument)
		return
	}

	pack, err := pc.packageService.Describe(id)
	if err != nil {
		log.Debug().Err(err).Msg("packageService.Describe failed")
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, serviceErrMsg)); err != nil {
			log.Err(err).Msg("MypackageCommander.Describe: error sending reply message to chat")
		}

		return
	}

	volume := pack.MaximumVolume
	units := "мл"
	if pack.MaximumVolume > 1000 {
		volume = pack.MaximumVolume / 1000
		units = "л"
	}

	reusableStr := "нет"
	if pack.Reusable {
		reusableStr = "да"
	}

	description := fmt.Sprintf(template, pack.Title, pack.Material, volume, units, reusableStr)

	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, description)

	if _, err := pc.bot.Send(msg); err != nil {
		log.Debug().Err(err).Msgf("MypackageCommander.Describe: error sending reply message to chat")
		return
	}

	log.Debug().Msgf("MypackageCommander.Describe: package id %d described", id)
}
