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
	commandArgs := inputMsg.CommandArguments()
	args := strings.Split(commandArgs, " ")
	if len(args) < 2 {
		log.Printf("MypackageCommander.Edit: wrong command arguments: %s", commandArgs)
		return
	}

	idx, err := strconv.Atoi(args[0])
	if err != nil {
		log.Printf("MypackageCommander.Edit: cannot parse int from command argument: %s", commandArgs)
		return
	}

	newPackage := logistic.Package{
		Title: args[1],
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
