package logistic

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/grpc"

	"github.com/hablof/omp-bot/internal/app/commands/logistic/packageApi"
	"github.com/hablof/omp-bot/internal/app/path"
	"github.com/hablof/omp-bot/internal/config"
	"github.com/hablof/omp-bot/internal/service/logistic/mypackage"
)

type Commander interface {
	HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath)
	HandleCommand(message *tgbotapi.Message, commandPath path.CommandPath)
}

var _ Commander = &LogisticCommander{}

type LogisticCommander struct {
	bot                *tgbotapi.BotAPI
	mypackageCommander Commander
}

func NewLogisticCommander(bot *tgbotapi.BotAPI, cc grpc.ClientConnInterface, cfg *config.Config) *LogisticCommander {
	return &LogisticCommander{
		bot: bot,
		// subdomainCommander
		mypackageCommander: packageApi.NewMypackageCommander(bot, mypackage.NewService(cc), cfg),
	}
}

func (c *LogisticCommander) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.Subdomain {
	case "package":
		c.mypackageCommander.HandleCallback(callback, callbackPath)
	default:
		log.Printf("DemoCommander.HandleCallback: unknown subdomain - %s", callbackPath.Subdomain)
	}
}

func (c *LogisticCommander) HandleCommand(msg *tgbotapi.Message, commandPath path.CommandPath) {
	switch commandPath.Subdomain {
	case "package":
		c.mypackageCommander.HandleCommand(msg, commandPath)
	default:
		log.Printf("DemoCommander.HandleCommand: unknown subdomain - %s", commandPath.Subdomain)
	}
}
