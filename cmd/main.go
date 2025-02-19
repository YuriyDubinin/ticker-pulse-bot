package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	bot "ticker-pulse-bot/internal/bot"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Загрузка количества воркеров из .env
	maxWorkers, err := strconv.Atoi(os.Getenv("MAX_WORKERS"))
	if err != nil {
		log.Fatal("[TICKER-PULSE-BOT]: MAX_WORKERS должно быть числом")
	}

	tickerPulseBot, err := bot.NewBot(maxWorkers)
	if err != nil {
		log.Fatalf("[TICKER-PULSE-BOT]: Ошибка инициализации бота: %v", err)
	}

	tickerPulseBot.Start()

	// Ожидание SIGINT/SIGTERM для корректного завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("[TICKER-PULSE-BOT]: Завершение работы...")
	tickerPulseBot.Stop()
}
