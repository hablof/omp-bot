package packageApi

import (
	"errors"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const descriptionTemplate = `Название: %s
Материал: %s
Максимальный объём: %3.0f %s
Многоразовая: %s`

// Describe implements PackageCommander
func (pc *MypackageCommander) Describe(inputMsg *tgbotapi.Message) {

	argument := inputMsg.CommandArguments()

	id, err := strconv.ParseUint(argument, 10, 64)
	if err != nil {
		pc.sendMsgWithErrLog(inputMsg, mtdDescribe, badRequestMsg)
		log.Err(err).Msgf("MypackageCommander.Describe: cannot parse int from command argument: %s", argument)

		return
	}

	errBadArg := &ErrBadArgument{}
	pack, err := pc.packageService.Describe(id)
	if err != nil {
		log.Debug().Err(err).Msg("packageService.Describe failed")

		switch {
		case errors.As(err, errBadArg):
			arg := pc.mapArg(errBadArg.Argument)

			log.Debug().Msgf("MypackageCommander.Describe: unknown argument: %s", arg)
			pc.sendMsgWithErrLog(inputMsg, mtdDescribe, fmt.Sprintf("Некорректный аргумент: \"%s\"", arg))

			return

		case err == ErrNotFound:
			log.Debug().Msgf("MypackageCommander.Describe: package [id: %d] not found", id)
			pc.sendMsgWithErrLog(inputMsg, mtdDescribe, fmt.Sprintf("Упаковка [id: %d] не найдена", id))

			return
		}

		pc.sendMsgWithErrLog(inputMsg, mtdDescribe, serviceErrMsg)

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

	msgText := fmt.Sprintf(descriptionTemplate, pack.Title, pack.Material, volume, units, reusableStr)

	if pc.sendMsgWithErrLog(inputMsg, mtdDescribe, msgText) {
		log.Debug().Msgf("MypackageCommander.Describe: package [id: %d] described", id)
		return
	}

	log.Debug().Msgf("MypackageCommander.Describe: package [id: %d] not described", id)
}
