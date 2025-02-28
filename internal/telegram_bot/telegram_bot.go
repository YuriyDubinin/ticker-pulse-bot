package telegram_bot

import (
	"fmt"
	"log"
	"os"
	"strings"

	quotes "ticker-pulse-bot/internal/pkg/quotes"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	godotenv "github.com/joho/godotenv"
)

type TelegramBot struct {
	api    *tgbotapi.BotAPI
	chatID string
}

func NewTelegramBot() (*TelegramBot, error) {
	envPath := os.Getenv("ENV_FILE")
	if envPath == "" {
		envPath = ".env" // значение по умолчанию
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("[TICKER-PULSE-BOT]: Ошибка загрузки .env")
	}

	apiKey := os.Getenv("TELEGRAM_BOT_API_KEY")
	chatID := os.Getenv("TELEGRAM_GROUP_ID")

	if apiKey == "" || chatID == "" {
		log.Fatal("[TICKER-PULSE-BOT]: TELEGRAM_BOT_API_KEY / TELEGRAM_GROUP_ID не найдено в .env")
	}

	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, err
	}

	log.Printf("[TICKER-PULSE-BOT]: Авторизован как %s\n", bot.Self.UserName)

	return &TelegramBot{
		api:    bot,
		chatID: chatID,
	}, nil
}

func (tb *TelegramBot) SendMessageToChannel(text string) error {
	msg := tgbotapi.NewMessageToChannel(tb.chatID, text)
	_, err := tb.api.Send(msg)
	return err
}

func (tb *TelegramBot) ConvertQuotesRateToMsg(data map[string]any) string {
	var result strings.Builder

	for _, quote := range quotes.Quotes {
		quoteRateMap, ok := data[quote.QuoteID].(map[string]any)
		if !ok {
			log.Printf("[TICKER-PULSE-BOT]: [ConvertQuotesRateToMsg]: Неверный тип для QuoteID %s\n", quote.QuoteID)
		}

		quoteUsdPrice, exists := quoteRateMap["usd"].(float64)
		if exists {
			result.WriteString(fmt.Sprintf("%-15s %-5s %15.2f $\n", quote.Label, quote.Ticker, quoteUsdPrice))
		} else {
			log.Printf("[TICKER-PULSE-BOT]: [ConvertQuotesRateToMsg]: Нет USD значения для %s\n", quote.QuoteID)
		}
	}

	return result.String()
}

// GUI
func (tb *TelegramBot) CreateKeyboard() error {
	inlineButtons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("📊 Курс актуальных котировок", "CURRENT_QUOTES_RATE"),
		},
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)

	msg := tgbotapi.NewMessageToChannel(tb.chatID, "📋 Основное меню: ")
	msg.ReplyMarkup = inlineKeyboard

	_, err := tb.api.Send(msg)

	return err
}

func (tb *TelegramBot) ListenKeyboardEvents(handlers map[string]func()) {
	updates, err := tb.api.GetUpdatesChan(tgbotapi.NewUpdate(0))
	if err != nil {
		log.Fatal("[TICKER-PULSE-BOT]: Ошибка получения обновлений: ", err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			if handler, exists := handlers[callbackData]; exists {
				handler()
			} else {
				log.Println("[TICKER-PULSE-BOT]: Неизвестная команда: ", callbackData)
			}
		}
	}
}
