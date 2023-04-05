package packageApi

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Delete implements PackageCommander
func (pc *MypackageCommander) Delete(inputMsg *tgbotapi.Message) {

	argument := inputMsg.CommandArguments()

	idx, err := strconv.Atoi(argument)
	if err != nil {
		log.Printf("MypackageCommander.Delete: cannot parse int from command argument: %s", argument)
		return
	}

	isRemoved, err := pc.packageService.Remove(uint64(idx))
	if err != nil {
		log.Printf("packageService.Remove: cannot remove: %v", err)
	}

	if isRemoved {
		log.Printf("packageService: entity #%d removed", idx)
	}
}
