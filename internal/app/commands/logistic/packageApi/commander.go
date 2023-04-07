package packageApi

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hablof/omp-bot/internal/app/path"
	"github.com/hablof/omp-bot/internal/config"
	"github.com/hablof/omp-bot/internal/model/logistic"
	"github.com/rs/zerolog/log"
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
	badRequestMsg = "Некорректный запрос"
)

type methodName string

const (
	mtdCreate   methodName = "MypackageCommander.Create"
	mtdDescribe methodName = "MypackageCommander.Describe"
	mtdList     methodName = "MypackageCommander.List"
	mtdRemove   methodName = "MypackageCommander.Remove"
	mtdUpdate   methodName = "MypackageCommander.Update"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrInternal   = errors.New("internal error")
	ErrNotFound   = errors.New("not found")
)

type ErrBadArgument struct {
	Argument string
}

func (e ErrBadArgument) Error() string {
	return "bad argumrnt: " + e.Argument
}

type PackageService interface {
	Create(createMap map[string]string) (uint64, error)
	Describe(packageID uint64) (logistic.Package, error)
	List(offset uint64, limit uint64) ([]logistic.Package, error)
	Remove(packageID uint64) (bool, error)
	Update(packageID uint64, editMap map[string]string) (bool, error)
}

type MypackageCommander struct {
	bot            *tgbotapi.BotAPI
	packageService PackageService
	paginationStep int
}

func NewMypackageCommander(bot *tgbotapi.BotAPI, packageService PackageService, cfg *config.Config) *MypackageCommander {
	// packageService := mypackage.NewService()

	return &MypackageCommander{
		bot:            bot,
		packageService: packageService,
		paginationStep: cfg.Tgbot.PaginationStep,
	}
}

func (c *MypackageCommander) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.CallbackName {
	case "list":
		c.CallbackList(callback, callbackPath)
	default:
		log.Debug().Msgf("MypackageCommander.HandleCallback: unknown callback name: %s", callbackPath.CallbackName)
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

// If sent succesfully returns true
func (pc *MypackageCommander) sendMsgWithErrLog(inputMsg *tgbotapi.Message, method methodName, text string) (sentWithoutErr bool) {
	if _, err := pc.bot.Send(tgbotapi.NewMessage(inputMsg.Chat.ID, text)); err != nil {
		log.Debug().Err(err).Msg(string(method) + ": error sending reply message to chat")
		return false
	}

	return true
}

func (pc *MypackageCommander) mapArg(s string) string {
	switch s {
	case "Material":
		return "Материал"

	case "Title":
		return "Название"

	case "MaximumVolume":
		return "Объём"

	case "Reusable":
		return "Многоразовая"
	}

	return s
}
