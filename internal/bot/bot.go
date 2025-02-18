package bot

import (
	"log"
	telegramBot "ticker-pulse-bot/internal/telegram_bot"
	workerpool "ticker-pulse-bot/internal/worker_pool"
)

type Bot struct {
	tgBot      *telegramBot.TelegramBot
	workerPool *workerpool.WorkerPool
}

func NewBot(maxWorkers int) (*Bot, error) {
	tgBot, err := telegramBot.NewTelegramBot()
	if err != nil {
		log.Fatal("[TICKER-PULSE-BOT]: Ошибка инициализации Telegram бота: ", err)
		return nil, err
	}

	wp := workerpool.New(maxWorkers)

	return &Bot{
		tgBot:      tgBot,
		workerPool: wp,
	}, nil
}

// Start запускает WorkerPool
func (b *Bot) Start() {
	b.workerPool.Start()
	log.Println("[TICKER-PULSE-BOT]: Бот запущен")
}

// Stop останавливает WorkerPool
func (b *Bot) Stop() {
	b.workerPool.Stop()
	log.Println("[TICKER-PULSE-BOT]: Бот остановлен")
}

// SendAsyncMessage отправляет сообщение асинхронно через WorkerPool
func (b *Bot) SendAsyncMessage(text string) {
	b.workerPool.AddTask(func() {
		if err := b.tgBot.SendMessage(text); err != nil {
			log.Printf("[TICKER-PULSE-BOT]: Ошибка отправки сообщения: %v", err)
		} else {
			log.Println("[TICKER-PULSE-BOT]: Сообщение успешно отправлено")
		}
	})
}
