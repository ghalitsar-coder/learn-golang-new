package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Job struct {
	Index int
	Value string
}

type Result struct {
	Index int
	Value string
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
	for j := range jobs {
		// simulate variable work
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		results <- Result{Index: j.Index, Value: fmt.Sprintf("worker-%d processed %s", id, j.Value)}
	}
}

// collector menerima results yang bisa datang out-of-order dan mencetaknya
// dalam urutan Index ascending.
func collector(results <-chan Result, total int) {
	pending := make(map[int]Result)
	next := 0
	for received := 0; received < total; {
		r := <-results
		pending[r.Index] = r
		// print semua yang berurutan
		for {
			if res, ok := pending[next]; ok {
				fmt.Printf("OUT→IN ORDER: index=%d result=%s\n", res.Index, res.Value)
				delete(pending, next)
				next++
				received++
			} else {
				break
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	jobs := make(chan Job)
	results := make(chan Result)

	numJobs := 10
	// start workers
	for w := 0; w < 3; w++ {
		go worker(w, jobs, results)
	}

	// start collector
	go collector(results, numJobs)

	// send jobs
	for i := 0; i < numJobs; i++ {
		jobs <- Job{Index: i, Value: fmt.Sprintf("job-%d", i)}
	}
	close(jobs)

	// wait enough time
	time.Sleep(2 * time.Second)
}
