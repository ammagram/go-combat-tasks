package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrTimeout = errors.New("query timeout")

func queryDB() string {
	time.Sleep(500 * time.Millisecond)
	return "result"
}

// BUG: канал result небуферизованный и после load() читать из него некому, происходит блокировка
// FIX: можно сделать канал с буфером 1, или закрыть done через defer при любом выходе
// GUARANTEE: горутина завершается при любом исходе, результат после таймаута отбрасывается
func load() (string, error) {
	result := make(chan string)
	done := make(chan struct{})
	defer close(done)

	go func() {
		value := queryDB()

		select {
		case result <- value:
		case <-done:
		}
	}()

	select {
	case value := <-result:
		return value, nil

	case <-time.After(200 * time.Millisecond):
		return "", ErrTimeout
	}
}

func main() {
	value, err := load()

	fmt.Println("value:", value)
	fmt.Println("error:", err)
}
