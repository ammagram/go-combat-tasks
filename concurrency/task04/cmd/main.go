package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

func worker(
	jobs <-chan Job,
) {
	for job := range jobs {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("processed:", job.ID)
	}
}

// BUG:
// нет блокировки продюсера, все 100 000 джоб отправляются разом и создают 100 000 горутин
//
// FIX:
// синхронно отправлять джобы воркерам. когда буфер канала полон - отправка блокируется и тормозит до скрорости воркеров
//
// GUARANTEE:
// кол-во горутин = кол-во воркеров, продюсер не обгоняет потребителей, все задачи отправлены до закрытия канала, wg.Wait() после обработки всех джоб

func main() {
	const workers = 4
	const queueSize = 8
	const jobsCount = 100_000

	jobs := make(chan Job, queueSize)

	var wg sync.WaitGroup

	for range workers {
		wg.Go(func() {
			worker(jobs)
		})
	}

	for i := range jobsCount {
		job := Job{ID: i}

		jobs <- job
	}

	close(jobs)
	wg.Wait()
}
