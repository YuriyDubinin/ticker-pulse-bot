package bot

import (
	"log"
	telegramBot "ticker-pulse-bot/internal/telegram_bot"
	workerPool "ticker-pulse-bot/internal/worker_pool"
)

type Bot struct {
	tgBot      *telegramBot.TelegramBot
	workerPool *workerPool.WorkerPool
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
	b.ListenKeyboardEventsAsync()

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

func (b *Bot) ListenKeyboardEventsAsync() {
	b.workerPool.AddTask(func() {
		b.tgBot.ListenKeyboardEvents()
	})
}
