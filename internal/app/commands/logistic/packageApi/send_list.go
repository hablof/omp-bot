package packageApi

import (
	"encoding/json"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hablof/omp-bot/internal/app/path"
	"github.com/rs/zerolog/log"
)

// const defaultLimit = 5

func (pc *MypackageCommander) sendList(offset int, inputMsg *tgbotapi.Message) {

	if offset < 0 {
		offset = 0
	}

	packages, err := pc.packageService.List(uint64(offset), uint64(pc.paginationStep))
	if err != nil {
		if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, badRequestMsg)); err != nil {
			log.Debug().Err(err).Msg("MypackageCommander.Create: error sending reply message to chat")
		}
		log.Debug().Err(err).Msg("packageService.List failed")

		return
	}

	var sb strings.Builder
	sb.WriteString("Вот список упаковок: \n")
	for i, p := range packages {
		sb.WriteString(fmt.Sprintf("\n%d. ", i+1+offset))
		sb.WriteString(p.String())
	}

	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, sb.String())

	serializedDataNextPage, _ := json.Marshal(CallbackListData{
		Offset: offset + pc.paginationStep,
	})

	serializedDataPervPage, _ := json.Marshal(CallbackListData{
		Offset: offset - pc.paginationStep,
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

	var row []tgbotapi.InlineKeyboardButton
	switch {
	case len(packages) < pc.paginationStep:
		row = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev page", callbackPathPrevPage.String()))

	case offset == 0:
		row = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Next page ➡️", callbackPathNextPage.String()))

	default:
		row = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev page", callbackPathPrevPage.String()),
			tgbotapi.NewInlineKeyboardButtonData("Next page ➡️", callbackPathNextPage.String()),
		)
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(row)

	if _, err := pc.bot.Send(msg); err != nil {
		log.Debug().Err(err).Msg("MypackageCommander.sendList: error sending reply message to chat")
		return
	}

	log.Debug().Msgf("MypackageCommander.sendList: list with offset %d sent", offset)
}
