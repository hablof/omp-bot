package packageApi

import (
	"errors"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

// Remove implements PackageCommander
func (pc *MypackageCommander) Remove(inputMsg *tgbotapi.Message) {

	argument := inputMsg.CommandArguments()

	id, err := strconv.ParseUint(argument, 10, 64)
	if err != nil {
		pc.sendMsgWithErrLog(inputMsg, mtdRemove, badRequestMsg)
		log.Debug().Err(err).Msgf("MypackageCommander.Remove: cannot parse int from command argument: %s", argument)

		return
	}

	errBadArg := &ErrBadArgument{}
	isRemoved, err := pc.packageService.Remove(id)
	if err != nil {
		log.Debug().Err(err).Msg("packageService.Remove failed")

		switch {
		case errors.As(err, errBadArg):
			arg := pc.mapArg(errBadArg.Argument)

			log.Debug().Msgf("MypackageCommander.Remove: unknown argument: %s", arg)
			pc.sendMsgWithErrLog(inputMsg, mtdRemove, fmt.Sprintf("Некорректный аргумент: \"%s\"", arg))

			return

		case err == ErrNotFound:
			log.Debug().Msgf("MypackageCommander.Remove: package [id=%d] not found", id)
			pc.sendMsgWithErrLog(inputMsg, mtdRemove, fmt.Sprintf("Упаковка [id=%d] не найдена", id))

			return
		}

		pc.sendMsgWithErrLog(inputMsg, mtdRemove, serviceErrMsg)

		return
	}

	if isRemoved {
		log.Debug().Msgf("packageService: entity ID %d removed", id)
		pc.sendMsgWithErrLog(inputMsg, mtdRemove, fmt.Sprintf("Упаковка ID=%d успешно удалена", id))
	} else {
		log.Debug().Msgf("packageService: entity ID %d NOT removed", id)
		pc.sendMsgWithErrLog(inputMsg, mtdRemove, fmt.Sprintf("Упаковка ID=%d НЕ удалена", id))
	}
}
