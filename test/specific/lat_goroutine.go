package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// APPROACH 1: Worker Pool Pattern (YANG ANDA GUNAKAN - RECOMMENDED)
func workerPoolDemo() {
	fmt.Println("✅ WORKER POOL PATTERN (3 workers untuk 11 jobs)")
	fmt.Println("==============================================")

	start := time.Now()
	jobs := make(chan int, 10)
	results := make(chan string, 11)

	var wg sync.WaitGroup

	// 3 workers yang akan mengerjakan 11 jobs
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			fmt.Printf("🏃 Worker-%d: SIAP bekerja\n", workerId)

			for job := range jobs {
				fmt.Printf("🔄 Worker-%d: Mulai job-%d\n", workerId, job)

				// simulasi pekerjaan
				workTime := time.Duration(rand.Intn(500)) * time.Millisecond
				time.Sleep(workTime)

				result := fmt.Sprintf("Worker-%d completed job-%d in %v", workerId, job, workTime)
				results <- result
				fmt.Printf("✅ Worker-%d: Selesai job-%d\n", workerId, job)
			}
			fmt.Printf("🏁 Worker-%d: Tidak ada job lagi, berhenti\n", workerId)
		}(i)
	}

	// Kirim 11 jobs
	fmt.Println("\n📤 Mengirim 11 jobs...")
	for i := 0; i < 11; i++ {
		jobs <- i
		fmt.Printf("📤 Job-%d dikirim\n", i)
	}
	close(jobs)

	// Tunggu semua worker selesai
	wg.Wait()
	close(results)

	// Kumpulkan hasil
	fmt.Println("\n📊 HASIL:")
	for result := range results {
		fmt.Println("   " + result)
	}

	fmt.Printf("\n⏱️ Total waktu: %v\n", time.Since(start))
}

// APPROACH 2: One Goroutine Per Job (TIDAK DIREKOMENDASIKAN)
func oneGoroutinePerJobDemo() {
	fmt.Println("\n🔥 ONE GOROUTINE PER JOB (11 goroutines untuk 11 jobs)")
	fmt.Println("=====================================================")

	start := time.Now()
	results := make(chan string, 11)
	var wg sync.WaitGroup

	// 11 goroutines untuk 11 jobs
	for i := 0; i < 11; i++ {
		wg.Add(1)
		go func(jobId int) {
			defer wg.Done()
			fmt.Printf("🔄 Goroutine-%d: Mulai job-%d\n", jobId, jobId)

			// simulasi pekerjaan
			workTime := time.Duration(rand.Intn(500)) * time.Millisecond
			time.Sleep(workTime)

			result := fmt.Sprintf("Goroutine-%d completed job-%d in %v", jobId, jobId, workTime)
			results <- result
			fmt.Printf("✅ Goroutine-%d: Selesai job-%d\n", jobId, jobId)
		}(i)
	}

	// Tunggu semua goroutine selesai
	wg.Wait()
	close(results)

	// Kumpulkan hasil
	fmt.Println("\n📊 HASIL:")
	for result := range results {
		fmt.Println("   " + result)
	}

	fmt.Printf("\n⏱️ Total waktu: %v\n", time.Since(start))
}

func main() {
	fmt.Println("🏭 PERBANDINGAN: WORKER POOL vs ONE GOROUTINE PER JOB")
	fmt.Println("====================================================")

	// Demo kode Anda (Worker Pool)
	workerPoolDemo()

	// Demo alternatif (One Goroutine Per Job)
	oneGoroutinePerJobDemo()

	fmt.Println(`
🤔 MENGAPA TIDAK N GOROUTINES UNTUK N JOBS?

1️⃣ RESOURCE OVERHEAD:
   • Setiap goroutine: ~8KB stack memory
   • 11 jobs = 11 goroutines = ~88KB
   • 1000 jobs = 1000 goroutines = ~8MB
   • 100k jobs = 100k goroutines = ~800MB!

2️⃣ SCHEDULER OVERHEAD:
   • Go scheduler harus manage semua goroutines
   • Context switching cost meningkat
   • CPU cache misses lebih sering

3️⃣ TIDAK SCALABLE:
   • Bagaimana jika 1 juta jobs?
   • System bisa kehabisan memory
   • Performance degradation

4️⃣ UNCONTROLLED CONCURRENCY:
   • Tidak ada batasan berapa goroutine aktif
   • Bisa overload system resources
   • Sulit untuk monitoring dan debugging

✅ KEUNGGULAN WORKER POOL:
   • Fixed number of workers (predictable resources)
   • Better resource management
   • Natural backpressure
   • Scalable untuk any number of jobs
   • Better performance untuk large workloads

🎯 RULE OF THUMB:
   • CPU-bound tasks: workers ≈ CPU cores
   • I/O-bound tasks: workers ≈ 2-4x CPU cores
   • Your case: 3 workers untuk any number of jobs ✅
`)
}
