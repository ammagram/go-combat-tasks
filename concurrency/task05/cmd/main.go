package main

import "fmt"

func generate(n int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()

	return out
}

func even(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for value := range in {
			if value%2 == 0 {
				out <- value
			}
		}
	}()

	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for value := range in {
			out <- value * value
		}
	}()

	return out
}

// BUG:
// владелец канала не закрывает его
//
// FIX:
// закрыть каналы, после завершения отправки
//
// GUARANTEE:
// все каналы завершаются, порядок сохраняется

func main() {
	numbers := generate(10)
	evens := even(numbers)
	squares := square(evens)

	for value := range squares {
		fmt.Println(value)
	}

	fmt.Println("pipeline completed")
}
