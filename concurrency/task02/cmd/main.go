package main

import (
	"fmt"
	"sync"
	"time"
)

type Result struct {
	URL    string
	Status int
}

// BUG:
// wg передается по значению, wg.Add добавляется все сразу, fetch знает о wg
//
// FIX:
// использовать wg.Go или функцию обертку для вызова fetch
//
// GUARANTEE:
// wg ждет все результаты вне зависимости от кол-ва запросов
func fetch(
	url string,
	results chan<- Result,
) {
	time.Sleep(100 * time.Millisecond)

	results <- Result{
		URL:    url,
		Status: 200,
	}
}

func main() {
	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://golang.org",
	}

	results := make(chan Result, len(urls))

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Go(func() {
			fetch(url, results)
		})
	}

	wg.Wait()
	close(results)

	i := 0
	for result := range results {
		i++
		fmt.Printf(
			"%d - %s -> %d\n",
			i,
			result.URL,
			result.Status,
		)
	}

	fmt.Println("done")
}
