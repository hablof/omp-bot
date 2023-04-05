package packageApi

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hablof/omp-bot/internal/model/logistic"
)

// New implements PackageCommander
func (pc *MypackageCommander) New(inputMsg *tgbotapi.Message) {
	args := strings.Split(inputMsg.CommandArguments(), ";")

	// количестпо полей, не считая поле ID
	if len(args) != logistic.PackageFieldsCount-1 {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "неверное количество аргументов")); err != nil {
			log.Printf("MypackageCommander.New: error sending reply message to chat - %v", err)
		}
		log.Printf("MypackageCommander.New: wrong args count")

		return
	}

	createArgMap := make(map[string]string, logistic.PackageFieldsCount-1)

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, logistic.Title):
			createArgMap[logistic.Title] = strings.TrimSpace(strings.TrimPrefix(arg, logistic.Title))

		case strings.HasPrefix(arg, logistic.Material):
			createArgMap[logistic.Material] = strings.TrimSpace(strings.TrimPrefix(arg, logistic.Material))

		case strings.HasPrefix(arg, logistic.MaximumVolume):
			createArgMap[logistic.MaximumVolume] = strings.TrimSpace(strings.TrimPrefix(arg, logistic.MaximumVolume))

		case strings.HasPrefix(arg, logistic.Reusable):
			createArgMap[logistic.Reusable] = strings.TrimSpace(strings.TrimPrefix(arg, logistic.Reusable))

		default:
			log.Printf("MypackageCommander.Edit: found argument: %s", arg)
			pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("Неизвестный аргумент: \"%s\"", arg)))
			return
		}

	}

	u, err := pc.packageService.Create(createArgMap)
	if err != nil {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "bad request")); err != nil {
			log.Printf("MypackageCommander.New: error sending reply message to chat - %v", err)
		}
		return
	}

	if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("New package id: %d", u))); err != nil {
		log.Printf("MypackageCommander.New: error sending reply message to chat - %v", err)
	}
}
