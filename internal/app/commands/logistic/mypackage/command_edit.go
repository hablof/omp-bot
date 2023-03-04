package mypackage

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

	if len(args) != 5 {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "неверно количество аргументов")); err != nil {
			log.Printf("MypackageCommander.New: error sending reply message to chat - %v", err)
		}
		return
	}

	idx, err := strconv.Atoi(args[0])
	if err != nil {
		log.Printf("MypackageCommander.Edit: cannot parse int from command argument: %s", args[0])
		return
	}

	volume, err := strconv.ParseFloat(strings.TrimSpace(args[3]), 32)
	if err != nil {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "неверно указан объём")); err != nil {
			log.Printf("MypackageCommander.New: error sending reply message to chat - %v", err)
		}
		return
	}

	newPackage := logistic.Package{
		Title:         args[1],
		Material:      args[2],
		MaximumVolume: float32(volume),
		Reusable:      strings.ToLower(strings.TrimSpace(args[4])) == "да",
	}

	if err := pc.packageService.Update(uint64(idx), newPackage); err != nil {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, "bad request")); err != nil {
			log.Printf("MypackageCommander.Edit: error sending reply message to chat - %v", err)
		}
		return
	}

	if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, fmt.Sprintf("package #%d updated", idx))); err != nil {
		log.Printf("MypackageCommander.Edit: error sending reply message to chat - %v", err)
	}
}
