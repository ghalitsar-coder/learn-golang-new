package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("========== DEMO: Sequential vs Parallel Select (Dramatic!) ==========\n")

	fmt.Println("Skenario: Ada 3 task yang masing-masing butuh 2 detik\n")

	fmt.Println("--- SEQUENTIAL: Loop Select ---")
	testSequentialDramatic()

	fmt.Println("\n--- PARALLEL: Select dalam Goroutine ---")
	testParallelDramatic()
}

// ===== SEQUENTIAL: Total 6 detik! =====
func testSequentialDramatic() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)

	// Task 1: Butuh 2 detik
	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "Task 1 selesai"
	}()

	// Task 2: Butuh 2 detik (dimulai bersamaan dengan Task 1)
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "Task 2 selesai"
	}()

	// Task 3: Butuh 2 detik (dimulai bersamaan dengan Task 1 & 2)
	go func() {
		time.Sleep(2 * time.Second)
		ch3 <- "Task 3 selesai"
	}()

	startTime := time.Now()
	received := 0

	fmt.Println("  Menjalankan tasks...")

	// Loop select: Process satu-satu secara sequential
	for received < 3 {
		select {
		case msg := <-ch1:
			elapsed := time.Since(startTime).Seconds()
			fmt.Printf("  [%.1fs] ✅ %s\n", elapsed, msg)
			received++
		case msg := <-ch2:
			elapsed := time.Since(startTime).Seconds()
			fmt.Printf("  [%.1fs] ✅ %s\n", elapsed, msg)
			received++
		case msg := <-ch3:
			elapsed := time.Since(startTime).Seconds()
			fmt.Printf("  [%.1fs] ✅ %s\n", elapsed, msg)
			received++
		}
	}

	totalTime := time.Since(startTime).Seconds()
	fmt.Printf("\n  ⏱️  Total waktu: %.1f detik\n", totalTime)
	fmt.Println("  📊 Penjelasan:")
	fmt.Println("     - Semua task selesai hampir bersamaan (~2 detik)")
	fmt.Println("     - Tapi loop select ambil data satu-satu dengan cepat")
	fmt.Println("     - Total tetap ~2 detik karena sudah ready semua!")
}

// ===== PARALLEL: Total 2 detik! =====
func testParallelDramatic() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)

	// Task 1: Butuh 2 detik
	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "Task 1 selesai"
	}()

	// Task 2: Butuh 2 detik
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "Task 2 selesai"
	}()

	// Task 3: Butuh 2 detik
	go func() {
		time.Sleep(2 * time.Second)
		ch3 <- "Task 3 selesai"
	}()

	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(3)

	fmt.Println("  Menjalankan tasks parallel...")

	// Goroutine 1: Langsung process hasil dari ch1
	go func() {
		defer wg.Done()
		msg := <-ch1
		elapsed := time.Since(startTime).Seconds()
		fmt.Printf("  [%.1fs] ✅ %s (processed by goroutine-1)\n", elapsed, msg)
		// Bisa langsung proses lebih lanjut di sini!
		time.Sleep(500 * time.Millisecond) // Simulasi processing
		fmt.Printf("  [%.1fs] 🔄 Task 1 fully processed!\n", time.Since(startTime).Seconds())
	}()

	// Goroutine 2: Langsung process hasil dari ch2
	go func() {
		defer wg.Done()
		msg := <-ch2
		elapsed := time.Since(startTime).Seconds()
		fmt.Printf("  [%.1fs] ✅ %s (processed by goroutine-2)\n", elapsed, msg)
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("  [%.1fs] 🔄 Task 2 fully processed!\n", time.Since(startTime).Seconds())
	}()

	// Goroutine 3: Langsung process hasil dari ch3
	go func() {
		defer wg.Done()
		msg := <-ch3
		elapsed := time.Since(startTime).Seconds()
		fmt.Printf("  [%.1fs] ✅ %s (processed by goroutine-3)\n", elapsed, msg)
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("  [%.1fs] 🔄 Task 3 fully processed!\n", time.Since(startTime).Seconds())
	}()

	wg.Wait()

	totalTime := time.Since(startTime).Seconds()
	fmt.Printf("\n  ⏱️  Total waktu: %.1f detik\n", totalTime)
	fmt.Println("  📊 Penjelasan:")
	fmt.Println("     - Tasks selesai ~2 detik (sama seperti sequential)")
	fmt.Println("     - TAPI processing bisa jalan PARALLEL!")
	fmt.Println("     - Masing-masing goroutine langsung process hasilnya")
	fmt.Println("     - Total dengan processing: ~2.5 detik (bukan 2+0.5+0.5+0.5=3.5!)")
	fmt.Println()
	fmt.Println("  ✅ KAPAN PAKAI PARALLEL SELECT?")
	fmt.Println("     - Ketika hasil perlu di-process lebih lanjut")
	fmt.Println("     - Processing bisa independent (tidak saling tunggu)")
	fmt.Println("     - Mau maximize throughput/concurrency")
	fmt.Println()
	fmt.Println("  ❌ Bukan lebih mahal, justru lebih EFISIEN untuk kasus tertentu!")
}
