package mypackage

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hablof/omp-bot/internal/app/path"
	"github.com/hablof/omp-bot/internal/service/logistic/mypackage"
)

type PackageCommander interface {
	Help(inputMsg *tgbotapi.Message)
	Get(inputMsg *tgbotapi.Message)
	List(inputMsg *tgbotapi.Message)
	Delete(inputMsg *tgbotapi.Message)
	New(inputMsg *tgbotapi.Message)
	Edit(inputMsg *tgbotapi.Message)
}

var _ PackageCommander = &MypackageCommander{}

type MypackageCommander struct {
	bot            *tgbotapi.BotAPI
	packageService mypackage.PackageService
}

func NewMypackageCommander(bot *tgbotapi.BotAPI) *MypackageCommander {
	packageService := mypackage.NewService()

	return &MypackageCommander{
		bot:            bot,
		packageService: packageService,
	}
}

func (c *MypackageCommander) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.CallbackName {
	case "list":
		c.CallbackList(callback, callbackPath)
	default:
		log.Printf("MypackageCommander.HandleCallback: unknown callback name: %s", callbackPath.CallbackName)
	}
}

func (pc *MypackageCommander) HandleCommand(msg *tgbotapi.Message, commandPath path.CommandPath) {
	switch commandPath.CommandName {
	case "help":
		pc.Help(msg)
	case "get":
		pc.Get(msg)
	case "list":
		pc.List(msg)
	case "new":
		pc.New(msg)
	case "delete":
		pc.Delete(msg)
	case "edit":
		pc.Edit(msg)
	}
}
