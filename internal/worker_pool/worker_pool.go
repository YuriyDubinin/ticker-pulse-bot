package worker_pool

import (
	"log"
	"sync"
)

type WorkerPool struct {
	tasksCh    chan func()
	maxWorkers int
	wg         sync.WaitGroup
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	return &WorkerPool{
		tasksCh:    make(chan func(), 100),
		maxWorkers: maxWorkers,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.maxWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	log.Printf("[TICKER-PULSE_BOT]: Запущено %d воркеров\n", wp.maxWorkers)
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for task := range wp.tasksCh {
		task()
	}

	log.Printf("[TICKER-PULSE_BOT]: Воркер %d завершает работу: все задачи выполнены\n", id)
}

func (wp *WorkerPool) AddTask(task func()) {
	select {
	case wp.tasksCh <- task:
	default:
		log.Println("[TICKER-PULSE_BOT]: Очередь задач переполнена, задача отклонена")
	}
}

func (wp *WorkerPool) Stop() {
	log.Println("[TICKER-PULSE_BOT]: Остановка WorkerPool..")

	close(wp.tasksCh) // Закрытие канала задачь, чтобы воркеры завершили обработку
	// wp.wg.Wait()      // Ожидание завершения всех воркеров

	log.Println("[TICKER-PULSE_BOT]: WorkerPool остановлен")
}
