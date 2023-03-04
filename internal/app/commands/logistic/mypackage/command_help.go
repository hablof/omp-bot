package mypackage

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const helpMsg = `/help__logistic__package — вывести список доступных команд 

/get__logistic__package {number} — вывести описание упаковки под номером {number}

/list__logistic__package — вывести список упаковок

/delete__logistic__package {number} — удалить упаковку под номером {number}

/new__logistic__package {title}; {material}; {volume}; {reusable}** — создать новую упаковку 

/edit__logistic__package {number}; {title}; {material}; {volume}; {reusable}** — редактировать упаковку под номером {number}


* - в миллилитрах без, указания единиц измерения
** - в формате "да" или "нет"`

// Help implements PackageCommander
func (pc *MypackageCommander) Help(inputMsg *tgbotapi.Message) {

	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, helpMsg)

	if _, err := pc.bot.Send(msg); err != nil {
		log.Printf("MypackageCommander.Help: error sending reply message to chat - %v", err)
	}
}
