package mypackage

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hablof/omp-bot/internal/model/logistic"
)

// New implements PackageCommander
func (pc *MypackageCommander) New(inputMsg *tgbotapi.Message) {
	newPackage := logistic.Package{
		Title: inputMsg.CommandArguments(),
	}

	u, err := pc.packageService.Create(newPackage)
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
