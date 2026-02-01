// Package main mendemonstrasikan KAPAN "1 Goroutine Per Job" BERBAHAYA.
// Kita akan mensimulasikan job yang membutuhkan memory cukup besar (misal: proses gambar/file).
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Ukuran memory per job (misal: load file 10MB ke memory)
const MemoryPerJob = 10 * 1024 * 1024 // 10 MB

// Dummy struct untuk memakan memory
type HeavyJob struct {
	data []byte
}

func doHeavyWork(id int) {
	// Simulasi alokasi memory
	_ = &HeavyJob{
		data: make([]byte, MemoryPerJob), // Alloc 10MB
	}
	// Simulasi processing time
	time.Sleep(100 * time.Millisecond)
}

func main() {
	fmt.Println("💥 DEMO: MEMORY EXHAUSTION (Crash Test)")
	fmt.Println("======================================")

	// Kita coba dengan 100 jobs dulu (1GB usage).
	// Ubah ke 500 atau 1000 untuk melihat crash/swap!
	numJobs := 100

	fmt.Printf("Skenario: %d jobs, setiap job butuh 10MB RAM.\n", numJobs)
	fmt.Printf("Total kebutuhan jika concurrent: %d MB RAM\n", numJobs*10)
	fmt.Println("-------------------------------------------")

	// 1. WORKER POOL (Safe)
	fmt.Println("\n✅ UJI COBA 1: Worker Pool (Limited to 10 workers)")
	fmt.Println("   Hanya perlu: 10 * 10MB = 100MB RAM (Aman)")
	runWorkerPool(numJobs, 10)

	// 2. UNLIMITED GOROUTINES (Dangerous)
	fmt.Println("\n🔥 Uji COBA 2: Unlimited Goroutines")
	fmt.Println("   Akan mencoba alokasi 5GB RAM sekaligus...")
	fmt.Println("   (Ini mungkin akan membuat laptop lagg atau crash program)")

	// Uncomment baris ini jika berani mencoba (resiko OOM)
	runUnlimitedGoroutines(numJobs)
}

func runWorkerPool(numJobs, numWorkers int) {
	start := time.Now()
	jobs := make(chan int, numJobs)
	var wg sync.WaitGroup

	// Monitor Memory Awal
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	fmt.Printf("   [Start] Sys Mem: %d MB\n", m1.Sys/1024/1024)

	// Start Workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for jobID := range jobs {
				doHeavyWork(jobID)
			}
		}()
	}

	for i := 0; i < numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	// Monitor Memory Peak (approximate)
	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			if m.Sys/1024/1024 > 500 { // If > 500MB
				// fmt.Printf("   [Alert] High usage: %d MB\n", m.Sys/1024/1024)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	wg.Wait()
	duration := time.Since(start)

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	fmt.Printf("   [End]   Sys Mem: %d MB\n", m2.Sys/1024/1024)
	fmt.Printf("   ⏱️  Selesai dalam: %v\n", duration)
}

func runUnlimitedGoroutines(numJobs int) {
	start := time.Now()
	var wg sync.WaitGroup

	// Monitor Memory Awal
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	fmt.Printf("   [Start] Sys Mem: %d MB\n", m1.Sys/1024/1024)

	for i := 0; i < numJobs; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			doHeavyWork(id)
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	fmt.Printf("   [End]   Sys Mem: %d MB\n", m2.Sys/1024/1024)
	fmt.Printf("   ⏱️  Selesai dalam: %v\n", duration)
}
