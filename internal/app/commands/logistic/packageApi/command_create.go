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
		pc.sendMsgWithErrLog(inputMsg, mtdCreate, "неверное количество аргументов")

		return
	}

	createArgMap := make(map[string]string, logistic.PackageFieldsCount-1)

	errBadArg := &ErrBadArgument{}
	if errors.As(pc.fillArgMap(args, createArgMap), errBadArg) {
		arg := pc.mapArg(errBadArg.Argument)
		log.Debug().Msgf("MypackageCommander.Create: unknown argument: %s", arg)
		pc.sendMsgWithErrLog(inputMsg, mtdCreate, fmt.Sprintf("Некорректный аргумент: \"%s\"", arg))

		return
	}

	id, err := pc.packageService.Create(createArgMap)
	if err != nil {
		log.Debug().Err(err).Msg("packageService.Create failed")

		if errors.As(err, errBadArg) {
			arg := pc.mapArg(errBadArg.Argument)
			log.Debug().Msgf("MypackageCommander.Update: unknown argument: %s", arg)
			pc.sendMsgWithErrLog(inputMsg, mtdCreate, fmt.Sprintf("Некорректный аргумент: \"%s\"", arg))

			return
		}

		pc.sendMsgWithErrLog(inputMsg, mtdCreate, serviceErrMsg)

		return
	}

	log.Debug().Msgf("MypackageCommander.Create: package id %d created", id)
	pc.sendMsgWithErrLog(inputMsg, mtdCreate, fmt.Sprintf("Добавлена новая упаковка [id: %d]", id))
}
