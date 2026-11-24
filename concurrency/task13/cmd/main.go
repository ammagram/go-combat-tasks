package main

import (
	"fmt"
	"sync"
	"time"
)

type Account struct {
	ID      int
	mu      sync.Mutex
	Balance int
}

// BUG: циклическое ожидание и deadlock, из-за того, что трансфер захватывает from и to в противоположных направлениях
// FIX: единый порядок блокировки по возрастанию ID счета
// GUARANTEE: нет циклического ожидания, суммарный баланс остается 2000, программа завершается
func transfer(
	from *Account,
	to *Account,
	amount int,
) {
	if from == to {
		return
	}

	first, second := from, to
	if second.ID < first.ID {
		first, second = second, first
	}

	first.mu.Lock()
	defer first.mu.Unlock()

	time.Sleep(10 * time.Millisecond)

	second.mu.Lock()
	defer second.mu.Unlock()

	from.Balance -= amount
	to.Balance += amount
}

func main() {
	a := &Account{
		ID:      1,
		Balance: 1000,
	}

	b := &Account{
		ID:      2,
		Balance: 1000,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		transfer(a, b, 100)
	}()

	go func() {
		defer wg.Done()
		transfer(b, a, 50)
	}()

	wg.Wait()

	fmt.Println(
		a.Balance,
		b.Balance,
	)
}
