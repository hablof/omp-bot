package mypackage

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const template = `Название: %s
Материал: %s
Максимальный объём: %3.0f %s
Переиспользование: %s`

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

	volume := pack.MaximumVolume
	units := "мл"
	if pack.MaximumVolume > 1000 {
		volume = pack.MaximumVolume / 1000
		units = "л"
	}

	reusableStr := "невозможно"
	if pack.Reusable {
		reusableStr = "возможно"
	}

	description := fmt.Sprintf(template, pack.Title, pack.Material, volume, units, reusableStr)

	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, description)

	if _, err := pc.bot.Send(msg); err != nil {
		log.Printf("MypackageCommander.Get: error sending reply message to chat - %v", err)
	}
}
