package main

import (
	"fmt"
	"time"
)

// Simple batching example: collect up to batchSize or flush after timeout

func processBatch(batch []int) {
	fmt.Printf("Processing batch: %v\n", batch)
	// simulate work
	time.Sleep(100 * time.Millisecond)
}

func main() {
	jobs := make(chan int)
	batchSize := 5
	flushAfter := 300 * time.Millisecond

	// producer
	go func() {
		for i := 1; i <= 17; i++ {
			jobs <- i
			time.Sleep(50 * time.Millisecond)
		}
		close(jobs)
	}()

	var batch []int
	flushTimer := time.NewTimer(flushAfter)
	defer flushTimer.Stop()

	for {
		if len(batch) == 0 {
			// reset timer when batch empty
			flushTimer.Reset(flushAfter)
		}
		select {
		case j, ok := <-jobs:
			if !ok {
				if len(batch) > 0 {
					processBatch(batch)
				}
				fmt.Println("done")
				return
			}
			batch = append(batch, j)
			if len(batch) >= batchSize {
				processBatch(batch)
				batch = nil
			}
		case <-flushTimer.C:
			if len(batch) > 0 {
				processBatch(batch)
				batch = nil
			}
		}
	}
}
