package router

import (
	"runtime/debug"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	// "github.com/hablof/omp-bot/internal/app/commands/demo"
	"github.com/hablof/omp-bot/internal/app/commands/logistic"
	"github.com/hablof/omp-bot/internal/app/path"
	"github.com/hablof/omp-bot/internal/config"
	"github.com/hablof/omp-bot/internal/service/logistic/mypackage"
)

const showCommandFormat = `Формат команд помощи: /help__{domain}__{subdomain}

Доступные команды:
 - /help__logistic__package`

type Commander interface {
	HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath)
	HandleCommand(callback *tgbotapi.Message, commandPath path.CommandPath)
}

type CommandSender interface {
	Send(update tgbotapi.Update) error
}

type Router struct {
	// bot
	bot           *tgbotapi.BotAPI
	commandSender CommandSender // kafka

	// demoCommander
	// demoCommander Commander
	// user
	// access
	// buy
	// delivery
	// recommendation
	// travel
	// loyalty
	// bank
	// subscription
	// license
	// insurance
	// payment
	// storage
	// streaming
	// business
	// work
	// service
	// exchange
	// estate
	// rating
	// security
	// cinema
	logistic Commander
	// product
	// education
}

func NewRouter(
	bot *tgbotapi.BotAPI,
	cc grpc.ClientConnInterface,
	cfg *config.Config,
	commandSender CommandSender,
	cache mypackage.CacheDict,
	ces mypackage.CacheEventSender,
) *Router {
	return &Router{
		// bot
		bot:           bot,
		commandSender: commandSender,
		// demoCommander
		// demoCommander: demo.NewDemoCommander(bot),
		// user
		// access
		// buy
		// delivery
		// recommendation
		// travel
		// loyalty
		// bank
		// subscription
		// license
		// insurance
		// payment
		// storage
		// streaming
		// business
		// work
		// service
		// exchange
		// estate
		// rating
		// security
		// cinema
		logistic: logistic.NewLogisticCommander(bot, cc, cache, ces, cfg),
		// product
		// education
	}
}

func (c *Router) HandleUpdate(update tgbotapi.Update) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			log.Debug().Msgf("recovered from panic: %v\n%v", panicValue, string(debug.Stack()))
		}
	}()

	go func(update tgbotapi.Update) {
		if err := c.commandSender.Send(update); err != nil {
			log.Debug().Err(err).Msg("failed send message to kafka")
		}
	}(update)

	switch {
	case update.CallbackQuery != nil:
		c.handleCallback(update.CallbackQuery)
	case update.Message != nil:
		c.handleMessage(update.Message)
	}
}

func (c *Router) handleCallback(callback *tgbotapi.CallbackQuery) {
	callbackPath, err := path.ParseCallback(callback.Data)
	if err != nil {
		log.Debug().Msgf("Router.handleCallback: error parsing callback data `%s` - %v", callback.Data, err)
		return
	}

	switch callbackPath.Domain {
	case "demo":
		// c.demoCommander.HandleCallback(callback, callbackPath)
	case "user":
		break
	case "access":
		break
	case "buy":
		break
	case "delivery":
		break
	case "recommendation":
		break
	case "travel":
		break
	case "loyalty":
		break
	case "bank":
		break
	case "subscription":
		break
	case "license":
		break
	case "insurance":
		break
	case "payment":
		break
	case "storage":
		break
	case "streaming":
		break
	case "business":
		break
	case "work":
		break
	case "service":
		break
	case "exchange":
		break
	case "estate":
		break
	case "rating":
		break
	case "security":
		break
	case "cinema":
		break
	case "logistic":
		c.logistic.HandleCallback(callback, callbackPath)
	case "product":
		break
	case "education":
		break
	default:
		log.Debug().Msgf("Router.handleCallback: unknown domain - %s", callbackPath.Domain)
	}
}

func (c *Router) handleMessage(msg *tgbotapi.Message) {
	switch {
	case !msg.IsCommand():
		fallthrough

	case path.IsGeneralCommand(msg.Command()):
		c.showCommandFormat(msg)
		return
	}

	commandPath, err := path.ParseCommand(msg.Command())
	if err != nil {
		log.Debug().Msgf("Router.handleCallback: error parsing callback data `%s` - %v", msg.Command(), err)
		return
	}

	switch commandPath.Domain {
	case "demo":
		// c.demoCommander.HandleCommand(msg, commandPath)
	case "user":
		break
	case "access":
		break
	case "buy":
		break
	case "delivery":
		break
	case "recommendation":
		break
	case "travel":
		break
	case "loyalty":
		break
	case "bank":
		break
	case "subscription":
		break
	case "license":
		break
	case "insurance":
		break
	case "payment":
		break
	case "storage":
		break
	case "streaming":
		break
	case "business":
		break
	case "work":
		break
	case "service":
		break
	case "exchange":
		break
	case "estate":
		break
	case "rating":
		break
	case "security":
		break
	case "cinema":
		break
	case "logistic":
		c.logistic.HandleCommand(msg, commandPath)
	case "product":
		break
	case "education":
		break
	default:
		log.Debug().Msgf("Router.handleCallback: unknown domain - %s", commandPath.Domain)
	}
}

func (c *Router) showCommandFormat(inputMessage *tgbotapi.Message) {
	outputMsg := tgbotapi.NewMessage(inputMessage.Chat.ID, showCommandFormat)

	_, err := c.bot.Send(outputMsg)
	if err != nil {
		log.Debug().Err(err).Msg("Router.showCommandFormat: error sending reply message to chat")
	}
}
