package bot

import (
	"fmt"
	"log"
	cryptoFetcher "ticker-pulse-bot/internal/crypto_fetcher"
	dataFormatter "ticker-pulse-bot/internal/pkg/data_formatter"
	quotes "ticker-pulse-bot/internal/pkg/quotes"
	telegramBot "ticker-pulse-bot/internal/telegram_bot"
	workerPool "ticker-pulse-bot/internal/worker_pool"
	"time"
)

type Bot struct {
	tgBot      *telegramBot.TelegramBot
	workerPool *workerPool.WorkerPool
	quotes     []quotes.QuoteInfo // временное решение до подьъема базы
}

func NewBot(maxWorkers int) (*Bot, error) {
	tgBot, err := telegramBot.NewTelegramBot()
	if err != nil {
		log.Fatal("[TICKER-PULSE-BOT]: Ошибка инициализации Telegram бота: ", err)
		return nil, err
	}

	wp := workerPool.NewWorkerPool(maxWorkers)

	return &Bot{
		tgBot:      tgBot,
		workerPool: wp,
		quotes:     quotes.Quotes,
	}, nil
}

// Запуск вместе с WorkerPool
// todo: на следующей итерации обеспечить гибкость currency
func (b *Bot) Start() {
	log.Println("[TICKER-PULSE-BOT]: Бот запущен")
	b.workerPool.Start()
	b.SendMessageAsync("Привет 🌍 Я тут, чтобы держать руку на пульсе, если что - дам знать. 🚀")
	b.CreateKeyboardAsync()
	b.ListenKeyboardEventsAsync(map[string]func(){
		"CURRENT_QUOTES_RATE": func() {
			err := b.tgBot.SendMessage("🔄 Запрашиваем актуальные котировки..")
			if err != nil {
				log.Println("[TICKER-PULSE-BOT]: Ошибка отправки сообщения: ", err)
			}

			quotesString, err := dataFormatter.FormatQuotesToString()
			if err != nil {
				log.Println("[TICKER-PULSE-BOT]: Ошибка форматирования котировок ", err)
			}

			b.ReportCurrentQuotesRateAsync(quotesString)
		},
	})
	b.CalculateQuotesInfo()
	b.CheckQuoteLimitsByInterval(3600)
}

// Остановка WorkerPool
func (b *Bot) Stop() {
	b.workerPool.Stop()
	log.Println("[TICKER-PULSE-BOT]: Бот остановлен")
}

// Отправка сообщений асинхронно через WorkerPool
func (b *Bot) SendMessageAsync(text string) {
	b.workerPool.AddTask(func() {
		if err := b.tgBot.SendMessage(text); err != nil {
			log.Printf("[TICKER-PULSE-BOT]: Ошибка отправки сообщения: %v", err)
		} else {
			log.Println("[TICKER-PULSE-BOT]: Сообщение успешно отправлено")
		}
	})
}

func (b *Bot) CreateKeyboardAsync() {
	b.workerPool.AddTask(func() {
		if err := b.tgBot.CreateKeyboard(); err != nil {
			log.Printf("[TICKER-PULSE-BOT]: Ошибка в создании GUI: %v", err)
		} else {
			log.Println("[TICKER-PULSE-BOT]: GUI успешно создан")
		}
	})
}

func (b *Bot) ListenKeyboardEventsAsync(handlers map[string]func()) {
	b.workerPool.AddTask(func() {
		b.tgBot.ListenKeyboardEvents(handlers)
	})
}

func (b *Bot) ReportCurrentQuotesRateAsync(quoteID string) {
	b.workerPool.AddTask(func() {
		cf := cryptoFetcher.NewCryptoFetcher()
		data, err := cf.FetchCoingeckoQuotesRate(quoteID, "usd")
		if err != nil {
			log.Printf("[TICKER-PULSE-BOT]: Ошибка получения данных котировки: %v", err)
			return
		}
		b.tgBot.SendMessage(b.tgBot.ConvertQuotesRateToMsg(data))
	})
}

func (b *Bot) CalculateQuotesInfo() {
	b.workerPool.AddTask(func() {
		cf := cryptoFetcher.NewCryptoFetcher()
		var updatedQuotes []quotes.QuoteInfo

		for _, quote := range b.quotes {
			time.Sleep(time.Duration(15) * time.Second)

			historicalData, err := cf.FetchCoingeckoHistoricalData(quote.QuoteID, 14)
			if err != nil {
				log.Printf("[TICKER-PULSE-BOT]: Ошибка получения исторических данных: %v", err)
				return
			}

			minMaxPrices, err := dataFormatter.CalculateHistoricalMinMax(historicalData, "prices")
			if err != nil {
				log.Printf("[TICKER-PULSE-BOT]: Ошибка форматирования исторических данных:  %v", err)
				return
			}
			log.Printf("%v min / max:  %v", quote.Ticker, minMaxPrices)

			quote.MinPrice = minMaxPrices.MinPrice
			quote.MaxPrice = minMaxPrices.MaxPrice
			quote.UpdatedAt = time.Now()

			updatedQuotes = append(updatedQuotes, quote)
		}

		b.quotes = updatedQuotes
		log.Println("[TICKER-PULSE-BOT]: Список котировок обновлен успешно")
	})
}

func (b *Bot) CheckQuoteLimitsByInterval(interval int) {
	cf := cryptoFetcher.NewCryptoFetcher()

	b.workerPool.AddTask(func() {
		for {
			quotesString, err := dataFormatter.FormatQuotesToString()
			if err != nil {
				log.Printf("[TICKER-PULSE-BOT]: Ошибка форматирования данных: %v", err)
				return

			}

			quotesRate, err := cf.FetchCoingeckoQuotesRate(quotesString, "usd") // TODO: улучшить гибкость метода
			if err != nil {
				log.Printf("[TICKER-PULSE-BOT]: Ошибка получения данных криптовалюты: %v", err)
				return
			}

			for _, quote := range b.quotes {
				quoteRateMap, ok := quotesRate[quote.QuoteID].(map[string]interface{})
				if !ok {
					log.Printf("[TICKER-PULSE-BOT]: Неверный тип для QuoteID %s\n", quote.QuoteID)
				}

				quoteUsdPrice, exists := quoteRateMap["usd"].(float64)
				if exists {
					log.Printf("%v: %v, min: %v, max: %v\n", quote.Label, quoteUsdPrice, quote.MinPrice, quote.MaxPrice)

					if quote.MinPrice != 0 && quoteUsdPrice < quote.MinPrice {
						msg := fmt.Sprintf("⬇️ %+v %+v, Спустился ниже недельного значения: %.2f $", quote.Label, quote.Ticker, quoteUsdPrice)
						err := b.tgBot.SendMessage(msg)
						if err != nil {
							log.Println("[TICKER-PULSE-BOT]: Ошибка отправки сообщения: ", err)
						}
					}

					if quote.MinPrice != 0 && quoteUsdPrice > quote.MaxPrice {
						msg := fmt.Sprintf("⬆️ %+v %+v, Поднялся выше недельного значения: %.2f $", quote.Label, quote.Ticker, quoteUsdPrice)
						err := b.tgBot.SendMessage(msg)
						if err != nil {
							log.Println("[TICKER-PULSE-BOT]: Ошибка отправки сообщения: ", err)
						}
					}

				} else {
					log.Printf("[TICKER-PULSE-BOT]: Нет USD значения для %s\n", quote.QuoteID)
				}
			}

			time.Sleep(time.Duration(interval) * time.Second)
		}
	})

}
