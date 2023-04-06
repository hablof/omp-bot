package packageApi

import (
	"encoding/json"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hablof/omp-bot/internal/app/path"
)

type CallbackListData struct {
	Offset int `json:"offset"`
}

func (pc *MypackageCommander) CallbackList(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {

	parsedData := CallbackListData{}
	if err := json.Unmarshal([]byte(callbackPath.CallbackData), &parsedData); err != nil {
		log.Printf("MypackageCommander.CallbackList: "+
			"error reading json data for type CallbackListData from "+
			"input string %v - %v", callbackPath.CallbackData, err)
		return
	}

	pc.sendList(parsedData.Offset, callback.Message)
}
