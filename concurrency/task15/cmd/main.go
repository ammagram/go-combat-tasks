package main

import (
	"context"
	"fmt"
	"time"
)

type Request struct {
	ID int
}

// BUG: создание нового тикера на каждой итерации цикла, игнорирование контекста
// FIX: один тикер на цикл
// GUARANTEE: не более rate обработок в секунду, остановка тикера после завершения run, выход при отмене контекста
func run(
	ctx context.Context,
	requests <-chan Request,
	rate int,
) {
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case request, ok := <-requests:
			if !ok {
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}

			fmt.Println(
				"processed:",
				request.ID,
			)
		}
	}
}

func main() {
	ctx := context.Background()

	requests := make(chan Request)

	go run(
		ctx,
		requests,
		5,
	)

	for i := 0; i < 20; i++ {
		requests <- Request{ID: i}
	}

	close(requests)

	time.Sleep(time.Second)
}
