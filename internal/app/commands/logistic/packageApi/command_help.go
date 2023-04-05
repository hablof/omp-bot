package packageApi

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const helpMsg = "`/help__logistic__package` — Вывести список доступных команд.\n\n🔹 `/list__logistic__package` — Вывести список упаковок.\n\n🔹 `/get__logistic__package {id}` — Вывести описание упаковки c идентификатором `{id}`.\n\n🔹 `/delete__logistic__package {id}` — Удалить упаковку под номером `{id}`.\n\n🔹 `/new__logistic__package название={title}; материал={material}; объём={volume*}; многоразовая={reusable**}` — Создать новую упаковку.\nПример: \"`/new__logistic__package название=Ведёрко; материал=Пластик; объём=700; многоразовая=да`\".\n\";\" — Обязательна в качестве разделителя.\n\n🔹 `/edit__logistic__package {id}; название={title}; материал={material}; объём={volume*}; многоразовая={reusable**}` — Редактировать упаковку под номером `{id}`.\nУказывать все поля необязательно. Пример: \"`/edit__logistic__package 15; объём=750; многоразовая=нет`\"\n\n\n* - в миллилитрах без, указания единиц измерения.\n** - в формате \"да\" или \"нет\".\n\";\" — Обязательна в качестве разделителя."

// Help implements PackageCommander
func (pc *MypackageCommander) Help(inputMsg *tgbotapi.Message) {

	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, helpMsg)

	if _, err := pc.bot.Send(msg); err != nil {
		log.Printf("MypackageCommander.Help: error sending reply message to chat - %v", err)
	}
}
