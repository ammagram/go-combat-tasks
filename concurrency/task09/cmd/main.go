package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

func worker(
	ctx context.Context,
	jobs <-chan Job,
	results chan<- int,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(200 * time.Millisecond):
			}

			select {
			case <-ctx.Done():
				return
			case results <- job.ID:
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	jobs := make(chan Job)
	results := make(chan int)

	var wg sync.WaitGroup
	wg.Add(1)

	go worker(
		ctx,
		jobs,
		results,
		&wg,
	)

	go func() {
		jobs <- Job{ID: 1}
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	wg.Wait()

	fmt.Println("worker stopped")
}
