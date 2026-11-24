package main

import (
	"fmt"
	"sync"
	"time"
)

type Client struct {
	ID int
}

// BUG: проверка на nil и запись  без синхронизации. горутины одновременно видят nil и создание клиента
// FIX: sync.Once первый вызов выполняет инициализацию, остальные блокируются до ее завершения и не запускаются повторно
// GUARANTEE: инициализация один раз, получение одного экземпляра клиента

var (
	client *Client
	once   sync.Once
)

func getClient() *Client {
	once.Do(func() {
		time.Sleep(100 * time.Millisecond)

		client = &Client{
			ID: 42,
		}
	})
	return client
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			c := getClient()
			fmt.Printf("%p\n", c)
		}()
	}

	wg.Wait()
}
