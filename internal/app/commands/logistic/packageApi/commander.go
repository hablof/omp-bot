package packageApi

import (
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hablof/omp-bot/internal/app/path"
	"github.com/hablof/omp-bot/internal/model/logistic"
)

// type PackageCommander interface {
// 	Help(inputMsg *tgbotapi.Message)
// 	Get(inputMsg *tgbotapi.Message)
// 	List(inputMsg *tgbotapi.Message)
// 	Delete(inputMsg *tgbotapi.Message)
// 	New(inputMsg *tgbotapi.Message)
// 	Edit(inputMsg *tgbotapi.Message)
// }

// var _ PackageCommander = &MypackageCommander{}
const (
	serviceErrMsg = "🤡🤡🤡 Ошибка сервиса 🤡🤡🤡"
	badRequestMsg = "некорректный запрос"
)

var ErrBadRequest = errors.New("bad request")

type PackageService interface {
	Describe(packageID uint64) (logistic.Package, error)
	List(cursor uint64, limit uint64) ([]logistic.Package, error)
	Create(createMap map[string]string) (uint64, error)
	Update(packageID uint64, editMap map[string]string) error
	Remove(packageID uint64) (bool, error)
}

type MypackageCommander struct {
	bot            *tgbotapi.BotAPI
	packageService PackageService
}

func NewMypackageCommander(bot *tgbotapi.BotAPI, packageService PackageService) *MypackageCommander {
	// packageService := mypackage.NewService()

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
		pc.Describe(msg)
	case "list":
		pc.List(msg)
	case "new":
		pc.Create(msg)
	case "delete":
		pc.Remove(msg)
	case "edit":
		pc.Update(msg)
	}
}
