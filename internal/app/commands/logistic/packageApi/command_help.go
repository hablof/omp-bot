package packageApi

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const helpMsg = "🔹/help\\_\\_logistic\\_\\_package — Вывести список доступных команд\\.\n\n🔹 /list\\_\\_logistic\\_\\_package — Вывести список упаковок\\.\n\n🔹 /get\\_\\_logistic\\_\\_package `\\[id\\]` — Вывести описание упаковки c идентификатором `\\[id\\]`\\.\n\n🔹 /delete\\_\\_logistic\\_\\_package `\\[id\\]` — Удалить упаковку под номером `\\[id\\]`\\.\n\n🔹 /new\\_\\_logistic\\_\\_package `название\\=\\[title\\]; материал\\=\\[material\\]; объём\\=\\[volume\\*\\]; многоразовая\\=\\[reusable\\*\\*\\]` — Создать новую упаковку\\.\nПример: \n\"/new\\_\\_logistic\\_\\_package `название\\=Ведёрко; материал\\=Пластик; объём\\=700; многоразовая\\=да`\"\\.\n\"`;`\" — Обязательна в качестве разделителя\\.\n\n🔹 /edit\\_\\_logistic\\_\\_package `\\[id\\]; название\\=\\[title\\]; материал\\=\\[material\\]; объём\\=\\[volume\\*\\]; многоразовая\\=\\[reusable\\*\\*\\]` — Редактировать упаковку под номером `\\[id\\]`\\.\nУказывать все поля необязательно\\. \nПример: \n\"/edit\\_\\_logistic\\_\\_package `15; объём\\=750; многоразовая\\=нет`\"\n\n\n\\* \\- в миллилитрах без, указания единиц измерения\\.\n\\*\\* \\- в формате \"да\" или \"нет\"\\.\n\"`;`\" — Обязательна в качестве разделителя\\."

// Help implements PackageCommander
func (pc *MypackageCommander) Help(inputMsg *tgbotapi.Message) {

	// s := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, helpMsg)
	// log.Debug().Msg(s)
	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, helpMsg)
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	if _, err := pc.bot.Send(msg); err != nil {
		log.Debug().Err(err).Msg("MypackageCommander.Help: error sending reply message to chat")
	}
}
