package packageApi

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// List implements PackageCommander
func (pc *MypackageCommander) List(inputMsg *tgbotapi.Message) {

	offset := 0
	pc.sendList(offset, inputMsg)
}
