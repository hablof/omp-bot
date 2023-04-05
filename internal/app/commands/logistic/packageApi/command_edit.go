package packageApi

import (
	"fmt"
	"log"
	"strconv"

	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hablof/omp-bot/internal/model/logistic"
)

// Edit implements PackageCommander
func (pc *MypackageCommander) Edit(inputMsg *tgbotapi.Message) {
	args := strings.Split(inputMsg.CommandArguments(), ";")

	// if len(args) != 5 {
	// 	if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "неверно количество аргументов")); err != nil {
	// 		log.Printf("MypackageCommander.New: error sending reply message to chat - %v", err)
	// 	}
	// 	return
	// }

	id, err := strconv.ParseUint(strings.TrimSpace(args[0]), 10, 64)
	if err != nil {
		log.Printf("MypackageCommander.Edit: cannot parse ID (int) from command argument: %s", args[0])
		pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("идентификатор не может быть \"%s\"", args[0])))
		return
	}

	// volume, err := strconv.ParseFloat(strings.TrimSpace(args[3]), 32)
	// if err != nil {
	// 	if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "неверно указан объём")); err != nil {
	// 		log.Printf("MypackageCommander.New: error sending reply message to chat - %v", err)
	// 	}
	// 	return
	// }
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
			log.Printf("MypackageCommander.Edit: found argument: %s", arg)
			pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("Неизвестный аргумент: \"%s\"", arg)))
			return
		}

	}

	// newPackage := logistic.Package{
	// 	Title:         args[1],
	// 	Material:      args[2],
	// 	MaximumVolume: float32(volume),
	// 	Reusable:      strings.ToLower(strings.TrimSpace(args[4])) == "да",
	// }

	if err := pc.packageService.Update(id, editArgMap); err != nil {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "bad request")); err != nil {
			log.Printf("MypackageCommander.Edit: error sending reply message to chat - %v", err)
		}
		return
	}

	if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("Упаковка ID=%d успешно отредактирована", id))); err != nil {
		log.Printf("MypackageCommander.Edit: error sending reply message to chat - %v", err)
	}
}
