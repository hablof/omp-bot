package mypackage

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Get implements PackageCommander
func (pc *MypackageCommander) Get(inputMsg *tgbotapi.Message) {

	argument := inputMsg.CommandArguments()

	idx, err := strconv.Atoi(argument)
	if err != nil {
		log.Printf("MypackageCommander.Get: cannot parse int from command argument: %s", argument)
		return
	}

	pack, err := pc.packageService.Describe(uint64(idx))
	if err != nil {
		log.Printf("packageService.Describe: error: %v", err)
		return
	}

	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, pack.Title)

	if _, err := pc.bot.Send(msg); err != nil {
		log.Printf("MypackageCommander.Get: error sending reply message to chat - %v", err)
	}
}
