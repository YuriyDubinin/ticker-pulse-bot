package bot

import (
	"log"
	"math/rand"
	cryptoFetcher "ticker-pulse-bot/internal/crypto_fetcher"
	dataFormatter "ticker-pulse-bot/internal/pkg/data_formatter"
	telegramMsgFmt "ticker-pulse-bot/internal/pkg/telegram_message_formatter"
	telegramBot "ticker-pulse-bot/internal/telegram_bot"
	workerPool "ticker-pulse-bot/internal/worker_pool"
	"time"
)

type QuoteInfo struct {
	CryptoID  string
	Ticker    string
	Label     string
	MinPrice  float64
	MaxPrice  float64
	UpdatedAt time.Time
}

type Bot struct {
	tgBot      *telegramBot.TelegramBot
	workerPool *workerPool.WorkerPool
	quotes     *[]QuoteInfo
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
	}, nil
}

// Запуск вместе с WorkerPool
func (b *Bot) Start() {
	b.workerPool.Start()

	b.SendMessageAsync("Привет 🌍 Я тут, чтобы держать руку на пульсе, если что - дам знать. 🚀")
	b.CreateKeyboardAsync()
	b.ListenKeyboardEventsAsync(map[string]func(){
		"CURRENT_QUOTES_RATE": func() {
			err := b.tgBot.SendMessage("🔄 Запрашиваем актуальные котировки..")
			if err != nil {
				log.Println("[TICKER-PULSE-BOT]: Ошибка отправки сообщения: ", err)
			}
			b.ReportCurrentQuotesRateAsync("bitcoin,ethereum,tether,binancecoin,usd-coin,ripple,cardano,dogecoin,solana,the-open-network")
		},
	})
	b.CalculateQuotesInfo()

	log.Println("[TICKER-PULSE-BOT]: Бот запущен")
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

// Report current quotes rate to chat
func (b *Bot) ReportCurrentQuotesRateAsync(cryptoID string) {
	b.workerPool.AddTask(func() {
		cf := cryptoFetcher.NewCryptoFetcher()
		data, err := cf.FetchCoingeckoQuotesRate(cryptoID, "usd")
		if err != nil {
			log.Printf("[TICKER-PULSE-BOT]: Ошибка получения данных криптовалюты: %v", err)
			return
		}
		b.tgBot.SendMessage(telegramMsgFmt.ConvertCryptoDataToString(data))
	})
}

func (b *Bot) CalculateQuotesInfo() {
	quotes := []QuoteInfo{
		{CryptoID: "bitcoin", Ticker: "BTC", Label: "Bitcoin"},
		{CryptoID: "ethereum", Ticker: "ETH", Label: "Ethereum"},
		{CryptoID: "tether", Ticker: "USDT", Label: "Tether"},
		{CryptoID: "binancecoin", Ticker: "BNB", Label: "Binance Coin"},
		{CryptoID: "usd-coin", Ticker: "USDC", Label: "USD Coin"},
		{CryptoID: "ripple", Ticker: "XRP", Label: "Ripple"},
		{CryptoID: "cardano", Ticker: "ADA", Label: "Cardano"},
		{CryptoID: "dogecoin", Ticker: "DOGE", Label: "Dogecoin"},
		{CryptoID: "solana", Ticker: "SOL", Label: "Solana"},
		{CryptoID: "the-open-network", Ticker: "TON", Label: "The Open Network"},
	}

	rand.Seed(time.Now().UnixNano()) // Инициализируем seed для случайных чисел

	b.workerPool.AddTask(func() {
		var updatedQuotes []QuoteInfo

		for _, quote := range quotes {
			cf := cryptoFetcher.NewCryptoFetcher()
			historicalData, err := cf.FetchCoingeckoHistoricalData(quote.CryptoID, 14)
			if err != nil {
				log.Printf("[TICKER-PULSE-BOT]:  %v", err)
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

			time.Sleep(time.Duration(15) * time.Second)
		}

		b.quotes = &updatedQuotes
		log.Println("[TICKER-PULSE-BOT]: Список котировок обновлен успешно")
	})
}
