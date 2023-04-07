package packageApi

import (
	"errors"
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
		pc.sendMsgWithErrLog(inputMsg, mtdUpdate, fmt.Sprintf("идентификатор не может быть \"%s\"", args[0]))

		return
	}

	editArgMap := make(map[string]string, logistic.PackageFieldsCount)

	errBadArg := &ErrBadArgument{}
	if errors.As(pc.fillArgMap(args[1:], editArgMap), errBadArg) { // args[1:] no info about id
		arg := pc.mapArg(errBadArg.Argument)
		log.Debug().Msgf("MypackageCommander.Update: unknown argument: %s", arg)
		pc.sendMsgWithErrLog(inputMsg, mtdUpdate, fmt.Sprintf("Некорректный аргумент: \"%s\"", arg))

		return
	}

	isUpdated, err := pc.packageService.Update(id, editArgMap)
	if err != nil {
		log.Debug().Err(err).Msg("packageService.Update failed")

		switch {
		case errors.As(err, errBadArg):
			arg := pc.mapArg(errBadArg.Argument)

			log.Debug().Msgf("MypackageCommander.Update: unknown argument: %s", arg)
			pc.sendMsgWithErrLog(inputMsg, mtdUpdate, fmt.Sprintf("Некорректный аргумент: \"%s\"", arg))

			return

		case err == ErrNotFound:
			log.Debug().Msgf("MypackageCommander.Update: package [id=%d] not found", id)
			pc.sendMsgWithErrLog(inputMsg, mtdUpdate, fmt.Sprintf("Упаковка [id=%d] не найдена", id))

			return
		}

		pc.sendMsgWithErrLog(inputMsg, mtdUpdate, serviceErrMsg)

		return
	}

	if isUpdated {
		log.Debug().Msgf("MypackageCommander.Update: package id %d updated", id)
		pc.sendMsgWithErrLog(inputMsg, mtdUpdate, fmt.Sprintf("Упаковка ID=%d успешно отредактирована", id))
	} else {
		log.Debug().Msgf("MypackageCommander.Update: package id %d NOT updated", id)
		pc.sendMsgWithErrLog(inputMsg, mtdUpdate, fmt.Sprintf("Упаковка ID=%d НЕ отредактирована", id))
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
