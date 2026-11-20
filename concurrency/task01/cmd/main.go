package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// BUG:
// 1000 горутин одновременно делают counter++
// это не одна операция, а чтение, увеличение и запись
// две горутины могут читать одно и то же значение.
// go run -race показывает ошибку
//
// FIX:
// синхронизировать доступ к counter
//
// GUARANTEE:
// после wg.Wait() counter будет равен goroutines*increments
// Wait происходит после всех Done
// go run -race не находят гонок

func solutionWithAtomic() {
	const goroutines = 1000
	const increments = 1000

	var wg sync.WaitGroup
	var counter atomic.Int64

	for i := 0; i < goroutines; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < increments; j++ {
				counter.Add(1)
			}
		}()
	}

	wg.Wait()

	fmt.Println("counter:", counter.Load())
}

func solutionWithMutex() {
	const goroutines = 1000
	const increments = 1000

	var wg sync.WaitGroup
	var mu sync.Mutex
	var counter int64

	for range goroutines {

		wg.Go(func() {
			mu.Lock()
			for range increments {
				counter++
			}
			mu.Unlock()
		})
	}

	wg.Wait()

	fmt.Println("counter:", counter)
}

func main() {
	solutionWithAtomic()
	solutionWithMutex()
}
