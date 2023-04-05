package packageApi

import (
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hablof/omp-bot/internal/app/path"
)

const defaultLimit = 5

func (pc *MypackageCommander) sendList(offset int, inputMsg *tgbotapi.Message) {
	products, err := pc.packageService.List(uint64(offset)+1, defaultLimit)
	if err != nil {
		log.Printf("packageService.List: error: %v", err)
		return
	}

	outputMsgText := "Here all the packages: \n"
	for i, p := range products {
		outputMsgText += fmt.Sprintf("\n%d. ", i+1+offset)
		outputMsgText += p.Title
	}

	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, outputMsgText)

	serializedDataNextPage, _ := json.Marshal(CallbackListData{
		Offset: offset + defaultLimit,
	})

	serializedDataPervPage, _ := json.Marshal(CallbackListData{
		Offset: offset - defaultLimit,
	})

	callbackPathNextPage := path.CallbackPath{
		Domain:       "logistic",
		Subdomain:    "package",
		CallbackName: "list",
		CallbackData: string(serializedDataNextPage),
	}

	callbackPathPrevPage := path.CallbackPath{
		Domain:       "logistic",
		Subdomain:    "package",
		CallbackName: "list",
		CallbackData: string(serializedDataPervPage),
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev page", callbackPathPrevPage.String()),
			tgbotapi.NewInlineKeyboardButtonData("Next page ➡️", callbackPathNextPage.String()),
		),
	)

	if _, err := pc.bot.Send(msg); err != nil {
		log.Printf("MypackageCommander.sendList: error sending reply message to chat - %v", err)
	}
}
