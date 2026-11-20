package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID int
}

func worker(
	id int,
	tasks <-chan Task,
	results chan<- string,
) {
	for task := range tasks {
		fmt.Printf(
			"worker %d started task %d\n",
			id,
			task.ID,
		)

		time.Sleep(50 * time.Millisecond)

		results <- fmt.Sprintf(
			"worker %d completed task %d",
			id,
			task.ID,
		)
	}
}

// BUG:
// FIX:
// GUARANTEE:

func main() {
	const numTasks = 100
	const numWorkers = 5

	tasks := make(chan Task)
	results := make(chan string)

	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Go(func() {
			worker(i, tasks, results)
		})
	}

	go func() {
		defer close(tasks)
		for i := 1; i <= numTasks; i++ {
			tasks <- Task{ID: i}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}

	fmt.Println("all tasks completed")
}
