package main

import (
	"fmt"
	"sync"
)

// BUG: чтение мапы и запись в мапу без синхронизации
// FIX: добавить мьютекс
// GUARANTEE: нет data race, читатель не видит мап в процессе записи
type Repository struct {
	mu     sync.RWMutex
	values map[int]int
}

func NewRepository() *Repository {
	return &Repository{
		values: make(map[int]int),
	}
}

func (r *Repository) Set(
	key int,
	value int,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.values[key] = value
}

func (r *Repository) Get(
	key int,
) (int, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.values[key]
	return value, ok
}

func main() {
	repo := NewRepository()

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(2)

		go func(id int) {
			defer wg.Done()

			repo.Set(id, id*10)
		}(i)

		go func(id int) {
			defer wg.Done()

			repo.Get(id)
		}(i)
	}

	wg.Wait()

	fmt.Println("completed")
}
