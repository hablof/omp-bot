package main

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/hablof/omp-bot/internal/app/cache"
	grpcclient "github.com/hablof/omp-bot/internal/app/grpc-client"
	"github.com/hablof/omp-bot/internal/app/kafka"
	routerPkg "github.com/hablof/omp-bot/internal/app/router"
	"github.com/hablof/omp-bot/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("failed to load ENV")
		// return
	}

	token, found := os.LookupEnv("TOKEN")
	if !found {
		log.Error().Msg("environment variable TOKEN not found in .env")
		return
	}

	cfg, err := config.ReadConfigYML("config.yml")
	if err != nil {
		log.Error().Err(err).Msg("failed to read config")
		return
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

	updates := bot.GetUpdatesChan(u)
	// if err != nil {
	// 	log.Error().Err(err).Msg("failed to setup tg updates")
	// 	return
	// }

	cc, err := grpcclient.NewConn(cfg)
	if err != nil {
		log.Error().Err(err).Msgf("failed to dial grpc-server")
		return
	}

	kafkaProducer, err := kafka.NewKafkaProducer(cfg.Kafka)
	if err != nil {
		log.Error().Err(err).Msgf("Failed init kafka")
		return
	}
	log.Info().Msg("Kafka producer started")

	redis, err := cache.NewCache(*cfg)
	if err != nil {
		log.Error().Err(err).Msgf("Failed init redis")
		return
	}
	log.Info().Msg("Redis client started")

	routerHandler := routerPkg.NewRouter(bot, cc, cfg, kafkaProducer, redis, kafkaProducer)

	log.Info().Msg("Ready to handle updates")
	for update := range updates {
		routerHandler.HandleUpdate(update)
	}
}
