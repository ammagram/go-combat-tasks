package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrStopped = errors.New("service stopped")

type Job struct {
	ID int
}

// BUG: submit и shutdown панакуют после закрытия канала, тк воркеры еще разбирают очередь. нет ожидания
// FIX: добавить флаг остановки и мьютекс для состояния закрытия. добавить ожидание воркеров
// GUARANTEE: все принятые задачи обработаны, воркеры завершены, остановка сервиса после воркеров. Сабмит после начала остановки корректно возвращает ошибку
type Service struct {
	mu      sync.RWMutex
	stopped bool
	jobs    chan Job
	wg      sync.WaitGroup
}

func NewService(workers int) *Service {
	s := &Service{
		jobs: make(chan Job, 16),
	}

	s.wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(id int) {
			defer s.wg.Done()

			for job := range s.jobs {
				time.Sleep(100 * time.Millisecond)

				fmt.Printf(
					"worker %d processed %d\n",
					id,
					job.ID,
				)
			}
		}(i)
	}

	return s
}

func (s *Service) Submit(job Job) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.stopped {
		return ErrStopped
	}

	s.jobs <- job

	return nil
}

func (s *Service) Shutdown() {
	s.mu.Lock()
	if !s.stopped {
		s.stopped = true
		close(s.jobs)
	}
	s.mu.Unlock()

	s.wg.Wait()

	fmt.Println("service stopped")
}

func main() {
	service := NewService(4)

	for i := 0; i < 20; i++ {
		if err := service.Submit(Job{ID: i}); err != nil {
			fmt.Println("submit failed:", err)
		}
	}

	service.Shutdown()
	fmt.Println("submit after shutdown:", service.Submit(Job{ID: 100}))
	fmt.Println("submit after shutdown:", service.Submit(Job{ID: 101}))
	fmt.Println("submit after shutdown:", service.Submit(Job{ID: 102}))
	service.Shutdown()
}
