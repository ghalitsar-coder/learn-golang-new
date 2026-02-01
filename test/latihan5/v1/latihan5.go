package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	//Output:
	context_test()

}

func context_test() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	maxJob := 10
	jobs := make(chan int)

	go func() {
		defer close(jobs)
		for i := range maxJob {
			select {
			case jobs <- i:
				fmt.Printf("Job-%d dikirim!\n", i)
				time.Sleep(500 * time.Millisecond)
			case <-ctx.Done():
				fmt.Println("Proses di hentikan\n")
				return
			}

		}
	}()

	// fmt.Printf("program berhenti 2 detik\n ")
	// time.Sleep(2 * time.Second)
	// fmt.Printf("program lanjut dan cancel\n")
	// cancel()

	 for {
        select {
        
        // 1. Cek Data & Status Channel
        case data, ok := <-jobs:
            if !ok {
                // Channel ditutup producer = Kerja selesai normal
                fmt.Println("Consumer: Semua pekerjaan selesai (Channel closed)")
                return 
            }
            fmt.Printf("Consumer: data-%d berhasil dikerjakan\n", data)
        // 2. Cek Cancel/Timeout
        case <-ctx.Done():
            fmt.Println("Consumer: Program dibatalkan paksa (Timeout/Cancel)")
            return
        }
    }

	//Output:

}
