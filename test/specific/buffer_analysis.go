package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func demonstrateJobDelivery() {
	fmt.Println("=== DEMO: APAKAH SEMUA JOB SUDAH TERKIRIM? ===")

	jobs := make(chan int, 10) // Buffer sama dengan jumlah job
	var wg sync.WaitGroup

	// Timestamp untuk tracking
	startTime := time.Now()

	fmt.Printf("[%v] 🚀 MULAI: Membuat 3 workers\n", time.Since(startTime).Round(time.Millisecond))

	// 3 Workers
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			for job := range jobs {
				elapsed := time.Since(startTime).Round(time.Millisecond)
				fmt.Printf("[%v] 🏃 Worker-%d: MULAI job-%d (sleep 2 detik)\n",
					elapsed, workerId, job)

				time.Sleep(2 * time.Second) // Simulasi kerja berat

				elapsed = time.Since(startTime).Round(time.Millisecond)
				fmt.Printf("[%v] ✅ Worker-%d: SELESAI job-%d\n",
					elapsed, workerId, job)
			}
		}(i)
	}

	// PENGIRIMAN JOBS
	fmt.Printf("[%v] 📤 MULAI mengirim 10 jobs...\n", time.Since(startTime).Round(time.Millisecond))

	for i := 0; i < 10; i++ {
		jobs <- i
		elapsed := time.Since(startTime).Round(time.Millisecond)
		fmt.Printf("[%v] ✉️  Job-%d TERKIRIM ke buffer\n", elapsed, i)
	}

	elapsed := time.Since(startTime).Round(time.Millisecond)
	fmt.Printf("[%v] 🎯 SEMUA 10 JOBS SUDAH TERKIRIM KE BUFFER!\n", elapsed)
	fmt.Printf("[%v] 🔒 Menutup channel...\n", elapsed)

	close(jobs)

	fmt.Printf("[%v] ⏳ Menunggu semua workers selesai...\n", elapsed)
	wg.Wait()

	finalTime := time.Since(startTime).Round(time.Millisecond)
	fmt.Printf("[%v] 🏁 SEMUA PEKERJAAN SELESAI!\n", finalTime)
}

func demonstrateWithoutBuffer() {
	fmt.Println("\n=== PERBANDINGAN: TANPA BUFFER ===")

	jobs := make(chan int) // UNBUFFERED
	var wg sync.WaitGroup
	startTime := time.Now()

	fmt.Printf("[%v] 🚀 Membuat 1 worker (unbuffered demo)\n", time.Since(startTime).Round(time.Millisecond))

	// 1 Worker saja untuk demo
	wg.Add(1)
	go func() {
		defer wg.Done()

		for job := range jobs {
			elapsed := time.Since(startTime).Round(time.Millisecond)
			fmt.Printf("[%v] 🏃 Worker: TERIMA job-%d, sleep 1 detik\n", elapsed, job)

			time.Sleep(1 * time.Second)

			elapsed = time.Since(startTime).Round(time.Millisecond)
			fmt.Printf("[%v] ✅ Worker: SELESAI job-%d\n", elapsed, job)
		}
	}()

	// Beri waktu worker untuk siap
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("[%v] 📤 MULAI mengirim 3 jobs (unbuffered)...\n", time.Since(startTime).Round(time.Millisecond))

	for i := 0; i < 3; i++ {
		elapsed := time.Since(startTime).Round(time.Millisecond)
		fmt.Printf("[%v] 📤 Mengirim job-%d (AKAN BLOCKING sampai diterima)\n", elapsed, i)

		jobs <- i // INI AKAN BLOCKING sampai worker menerima

		elapsed = time.Since(startTime).Round(time.Millisecond)
		fmt.Printf("[%v] ✉️  Job-%d BERHASIL diterima worker\n", elapsed, i)
	}

	elapsed := time.Since(startTime).Round(time.Millisecond)
	fmt.Printf("[%v] 🔒 Menutup channel...\n", elapsed)

	close(jobs)
	wg.Wait()

	finalTime := time.Since(startTime).Round(time.Millisecond)
	fmt.Printf("[%v] 🏁 SELESAI!\n", finalTime)
}

func demonstrateBufferVisualization() {
	fmt.Println("\n=== VISUALISASI BUFFER STATE ===")

	jobs := make(chan int, 5) // Buffer 5
	startTime := time.Now()

	// Helper function untuk print buffer state
	printBufferState := func(action string) {
		elapsed := time.Since(startTime).Round(time.Millisecond)
		fmt.Printf("[%v] %s | Buffer: %d/%d\n",
			elapsed, action, len(jobs), cap(jobs))
	}

	printBufferState("INIT")

	// Isi buffer
	for i := 1; i <= 5; i++ {
		jobs <- i
		printBufferState(fmt.Sprintf("SEND job-%d", i))
	}

	fmt.Println("\n🚨 BUFFER PENUH! Job berikutnya akan BLOCKING...")

	// Worker untuk mengambil data
	go func() {
		time.Sleep(2 * time.Second) // Delay untuk demo

		for i := 0; i < 3; i++ {
			job := <-jobs
			elapsed := time.Since(startTime).Round(time.Millisecond)
			fmt.Printf("[%v] 🏃 Worker: AMBIL job-%d | Buffer: %d/%d\n",
				elapsed, job, len(jobs), cap(jobs))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Coba kirim lagi (akan blocking)
	elapsed := time.Since(startTime).Round(time.Millisecond)
	fmt.Printf("[%v] 📤 Mencoba kirim job-6 (AKAN BLOCKING)...\n", elapsed)

	jobs <- 6 // INI AKAN BLOCKING sampai ada space

	printBufferState("SEND job-6 (after blocking)")

	close(jobs)
}

func main() {
	fmt.Println("🔍 ANALISIS: PENGIRIMAN JOB DENGAN BUFFER")
	fmt.Println("==========================================")

	demonstrateJobDelivery()

	fmt.Println("\n" + strings.Repeat("=", 50))
	demonstrateWithoutBuffer()

	fmt.Println("\n" + strings.Repeat("=", 50))
	demonstrateBufferVisualization()

	fmt.Println("\n🎯 KESIMPULAN:")
	fmt.Println("1. Dengan buffer = jumlah job → SEMUA job langsung terkirim")
	fmt.Println("2. Tanpa buffer → Pengiriman BLOCKING sampai ada penerima")
	fmt.Println("3. Buffer penuh → Job berikutnya akan BLOCKING")
	fmt.Println("4. len(channel) menunjukkan jumlah data dalam buffer")
}
