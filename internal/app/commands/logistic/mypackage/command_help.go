package mypackage

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const helpMsg = `/help__logistic__package — print list of commands
/get__logistic__package {number} — get a entity
/list__logistic__package — get a list of your entity 
/delete__logistic__package {number} — delete an existing entity
/new__logistic__package {title} — create a new entity 
/edit__logistic__package {number} {title} — edit a entity`

// Help implements PackageCommander
func (pc *MypackageCommander) Help(inputMsg *tgbotapi.Message) {

	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, helpMsg)

	if _, err := pc.bot.Send(msg); err != nil {
		log.Printf("MypackageCommander.Help: error sending reply message to chat - %v", err)
	}
}
