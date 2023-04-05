package packageApi

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hablof/omp-bot/internal/model/logistic"
)

// New implements PackageCommander
func (pc *MypackageCommander) New(inputMsg *tgbotapi.Message) {
	args := strings.Split(inputMsg.CommandArguments(), ";")

	if len(args) != 4 {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "неверно количество аргументов")); err != nil {
			log.Printf("MypackageCommander.New: error sending reply message to chat - %v", err)
		}
		return
	}

	volume, err := strconv.ParseFloat(strings.TrimSpace(args[2]), 32)
	if err != nil {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "неверно указан объём")); err != nil {
			log.Printf("MypackageCommander.New: error sending reply message to chat - %v", err)
		}
		return
	}

	newPackage := logistic.Package{
		Title:         args[0],
		Material:      args[1],
		MaximumVolume: float32(volume),
		Reusable:      strings.ToLower(strings.TrimSpace(args[3])) == "да",
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
