package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(
	id int,
	stop <-chan struct{},
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-stop:
			fmt.Println("worker stopped:", id)
			return

		case <-time.After(100 * time.Millisecond):
			fmt.Println("worker working:", id)
		}
	}
}

// BUG: становка делается отправкой пустой структуры одному воркеры, остальные продолжают крутиться в селект
// FIX: заменить отправку на закрытие канала
// GUARANTEE: сигнал наблюдают все воркеры, каждый воркер выходит и вызывает wg.Done(), main дожидается всех воркеров

func main() {
	const workers = 5

	stop := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go worker(i, stop, &wg)
	}

	time.Sleep(300 * time.Millisecond)

	close(stop)

	wg.Wait()

	fmt.Println("all workers stopped")
}
