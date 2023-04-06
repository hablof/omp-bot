package main

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hablof/omp-bot/internal/config"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	grpcclient "github.com/hablof/omp-bot/internal/app/grpc-client"
	routerPkg "github.com/hablof/omp-bot/internal/app/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("failed to load ENV")
		return
	}

	token, found := os.LookupEnv("TOKEN")
	if !found {
		log.Error().Msg("environment variable TOKEN not found in .env")
		return
	}

	cfg, err := config.ReadConfigYML("config.yaml")
	if err != nil {
		log.Error().Err(err).Msg("failed to read config")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Error().Err(err).Msg("failed to create telegram BotAPI instance")
		return
	}

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	bot.Debug = cfg.Tgbot.Debug
	u := tgbotapi.UpdateConfig{
		Timeout: 60,
	}

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Error().Err(err).Msg("failed to setup tg updates")
		return
	}

	cc, err := grpcclient.NewConn(cfg)
	if err != nil {
		log.Error().Err(err).Msgf("failed to dial grpc-server")
		return
	}

	routerHandler := routerPkg.NewRouter(bot, cc, cfg)

	for update := range updates {
		routerHandler.HandleUpdate(update)
	}
}
