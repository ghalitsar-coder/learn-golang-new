package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("========== PERBANDINGAN: Sequential vs Parallel Select ==========\n")

	fmt.Println("--- Test 1: SEQUENTIAL (Loop Select) ---")
	testSequentialSelect()

	fmt.Println("\n--- Test 2: PARALLEL (Select dalam Goroutine) ---")
	testParallelSelect()

	fmt.Println("\n--- Test 3: PARALLEL dengan Race Condition ---")
	testParallelRaceCondition()
}

// ===== TEST 1: SEQUENTIAL - Loop Select =====
func testSequentialSelect() {
	ch1 := make(chan string)
	ch2 := make(chan int)
	ch3 := make(chan bool)

	// Kirim data dengan delay berbeda
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "Data dari ch1"
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- 99
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ch3 <- true
	}()

	startTime := time.Now()
	received := 0

	// SEQUENTIAL: Satu select dalam loop
	for received < 3 {
		select {
		case msg := <-ch1:
			elapsed := time.Since(startTime).Milliseconds()
			fmt.Printf("  [%dms] 📥 Dapat dari ch1: %s\n", elapsed, msg)
			received++
		case num := <-ch2:
			elapsed := time.Since(startTime).Milliseconds()
			fmt.Printf("  [%dms] 📥 Dapat dari ch2: %d\n", elapsed, num)
			received++
		case flag := <-ch3:
			elapsed := time.Since(startTime).Milliseconds()
			fmt.Printf("  [%dms] 📥 Dapat dari ch3: %v\n", elapsed, flag)
			received++
		}
	}

	totalTime := time.Since(startTime).Milliseconds()
	fmt.Printf("\n  ⏱️  Total waktu SEQUENTIAL: %dms\n", totalTime)
	fmt.Println("  💡 Kenapa ~1 detik? Karena semua channel ready hampir bersamaan,")
	fmt.Println("     select bisa ambil langsung tanpa delay lagi.")
}

// ===== TEST 2: PARALLEL - Select dalam Goroutine Terpisah =====
func testParallelSelect() {
	ch1 := make(chan string)
	ch2 := make(chan int)
	ch3 := make(chan bool)

	// Kirim data dengan delay berbeda
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "Data dari ch1"
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- 99
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ch3 <- true
	}()

	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(3)

	// PARALLEL: Setiap select dalam goroutine sendiri!

	// Goroutine 1: Monitor ch1
	go func() {
		defer wg.Done()
		select {
		case msg := <-ch1:
			elapsed := time.Since(startTime).Milliseconds()
			fmt.Printf("  [%dms] 📥 Goroutine-1 dapat dari ch1: %s\n", elapsed, msg)
		}
	}()

	// Goroutine 2: Monitor ch2
	go func() {
		defer wg.Done()
		select {
		case num := <-ch2:
			elapsed := time.Since(startTime).Milliseconds()
			fmt.Printf("  [%dms] 📥 Goroutine-2 dapat dari ch2: %d\n", elapsed, num)
		}
	}()

	// Goroutine 3: Monitor ch3
	go func() {
		defer wg.Done()
		select {
		case flag := <-ch3:
			elapsed := time.Since(startTime).Milliseconds()
			fmt.Printf("  [%dms] 📥 Goroutine-3 dapat dari ch3: %v\n", elapsed, flag)
		}
	}()

	wg.Wait()

	totalTime := time.Since(startTime).Milliseconds()
	fmt.Printf("\n  ⏱️  Total waktu PARALLEL: %dms\n", totalTime)
	fmt.Println("  💡 Kenapa masih ~1 detik? Karena semua goroutine jalan parallel,")
	fmt.Println("     selesai hampir bersamaan setelah 1 detik.")
	fmt.Println("  ✅ Kelebihannya: Bisa process data independently/concurrent!")
}

// ===== TEST 3: PARALLEL dengan Race Condition =====
func testParallelRaceCondition() {
	ch := make(chan string)

	// Kirim 3 data
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- "Data 1"
		ch <- "Data 2"
		ch <- "Data 3"
		close(ch)
	}()

	var wg sync.WaitGroup
	wg.Add(3)

	fmt.Println("  ⚠️  Spawning 3 goroutine untuk COMPETE mendapat data dari 1 channel...")
	fmt.Println()

	// 3 goroutine COMPETE untuk data yang sama!
	for i := 1; i <= 3; i++ {
		workerID := i
		go func() {
			defer wg.Done()

			// Setiap worker coba ambil data
			if msg, ok := <-ch; ok {
				fmt.Printf("  🏆 Worker-%d dapat: %s\n", workerID, msg)
			} else {
				fmt.Printf("  ❌ Worker-%d channel sudah closed\n", workerID)
			}
		}()
	}

	wg.Wait()

	fmt.Println()
	fmt.Println("  💡 Yang terjadi:")
	fmt.Println("     - 3 goroutine compete untuk 3 data")
	fmt.Println("     - Masing-masing dapat 1 data (atau mungkin ada yang ga dapat)")
	fmt.Println("     - Non-deterministic siapa yang dapat apa")
	fmt.Println("     - Ini WORKER POOL pattern!")
}
