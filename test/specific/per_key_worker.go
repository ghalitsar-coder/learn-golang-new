package main

import (
	"fmt"
	"sync"
	"time"
)

// Simple per-key worker pattern: each key has its own goroutine to preserve
// order per key while allowing parallelism across different keys.

type KeyJob struct {
	Key string
	Val int
}

func perKeyWorker(key string, ch <-chan KeyJob, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range ch {
		// process sequentially per key
		fmt.Printf("[Key=%s] processing %d\n", key, job.Val)
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Printf("[Key=%s] worker stopped\n", key)
}

func main() {
	jobs := []KeyJob{
		{"user:1", 1}, {"user:2", 1}, {"user:1", 2}, {"user:3", 1},
		{"user:2", 2}, {"user:1", 3}, {"user:3", 2},
	}

	// dispatcher
	var mu sync.Mutex
	workers := make(map[string]chan KeyJob)
	var wg sync.WaitGroup

	for _, j := range jobs {
		mu.Lock()
		ch, ok := workers[j.Key]
		if !ok {
			ch = make(chan KeyJob, 10)
			workers[j.Key] = ch
			wg.Add(1)
			go perKeyWorker(j.Key, ch, &wg)
		}
		mu.Unlock()

		ch <- j
	}

	// close all
	for k, ch := range workers {
		close(ch)
		fmt.Printf("closed worker channel for %s\n", k)
	}

	wg.Wait()
	fmt.Println("all per-key workers done")
}
