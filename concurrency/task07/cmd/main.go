package main

import (
	"fmt"
	"sync"
)

func source(start, end int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for i := start; i <= end; i++ {
			out <- i
		}
	}()

	return out
}

// BUG: писатель закрывает канал. горитина, которая закончит первой закрывает out
// FIX: писатели не закрывают канал, только вызывают wg.Done(). горутина ожидающая wg.Wait() знает о завершении всех писателей и закрывает канал.
// GUARANTEE: канал закрывается один раз после того, как все источники исчерпаны. паник и потерь нет. приходят все значения (1-5, 100-105). порядок не детерменирован
func merge(
	a <-chan int,
	b <-chan int,
) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup

	copyValues := func(in <-chan int) {
		defer wg.Done()

		for value := range in {
			out <- value
		}
	}

	wg.Add(2)
	go copyValues(a)
	go copyValues(b)

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	a := source(1, 5)
	b := source(100, 105)

	for value := range merge(a, b) {
		fmt.Println(value)
	}
}
