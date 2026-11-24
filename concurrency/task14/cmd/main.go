package main

import (
	"fmt"
	"runtime"
	"sync"
)

type Stats struct {
	mu      sync.Mutex
	pending int64
	done    int64
}

func NewStats(total int64) *Stats {
	return &Stats{pending: total}
}

// BUG: нарушение инварианта, race condition
// FIX: изменение pending, done и чтение пары выполняется под мьютексом,
// GUARANTEE: любой snapshot выполняет инвариант pending >= 0, done >= 0 и pending + done = total
func (s *Stats) Complete() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pending--
	runtime.Gosched()
	s.done++
}

func (s *Stats) Snapshot() (
	pending int64,
	done int64,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.pending, s.done
}

func main() {
	const total = 1000

	stats := NewStats(total)

	var wg sync.WaitGroup
	wg.Add(total)

	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()

			stats.Complete()
		}()
	}

	for {
		pending, done := stats.Snapshot()

		if pending+done != total {
			fmt.Printf(
				"BROKEN SNAPSHOT: pending=%d done=%d\n",
				pending,
				done,
			)

			break
		}

		if done == total {
			break
		}
	}

	wg.Wait()
}
