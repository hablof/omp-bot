package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	routerPkg "github.com/hablof/omp-bot/internal/app/router"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
		return
	}

	token, found := os.LookupEnv("TOKEN")
	if !found {
		log.Println("environment variable TOKEN not found in .env")
		return
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println(err)
		return
	}

	// Uncomment if you want debugging
	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.UpdateConfig{
		Timeout: 60,
	}

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
		return
	}

	routerHandler := routerPkg.NewRouter(bot)

	for update := range updates {
		routerHandler.HandleUpdate(update)
	}
}
