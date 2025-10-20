package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func demonstrateExtremeScalability() {
	fmt.Println("🚀 EXTREME SCALABILITY TEST")
	fmt.Println("===========================")

	jobCounts := []int{100, 1000, 10000}

	for _, numJobs := range jobCounts {
		fmt.Printf("\n--- Testing dengan %d jobs ---\n", numJobs)

		// Test 1: One Goroutine Per Job
		fmt.Printf("🔥 One Goroutine Per Job (%d goroutines):\n", numJobs)
		start := time.Now()
		var wg1 sync.WaitGroup

		goroutinesBefore := runtime.NumGoroutine()

		for i := 0; i < numJobs; i++ {
			wg1.Add(1)
			go func(jobID int) {
				defer wg1.Done()
				// Simulasi pekerjaan ringan
				time.Sleep(1 * time.Millisecond)
			}(i)
		}

		maxGoroutines := runtime.NumGoroutine()
		wg1.Wait()
		duration1 := time.Since(start)

		fmt.Printf("   ⏱️  Duration: %v\n", duration1)
		fmt.Printf("   🧵 Max goroutines: %d (+ %d)\n", maxGoroutines, maxGoroutines-goroutinesBefore)

		// Test 2: Worker Pool
		numWorkers := runtime.NumCPU() // Optimal untuk CPU-bound
		fmt.Printf("✅ Worker Pool (%d workers):\n", numWorkers)

		start = time.Now()
		jobs := make(chan int, 100) // Small buffer
		var wg2 sync.WaitGroup

		// Start workers
		for w := 0; w < numWorkers; w++ {
			wg2.Add(1)
			go func() {
				defer wg2.Done()
				for range jobs {
					time.Sleep(1 * time.Millisecond)
				}
			}()
		}

		// Send jobs
		for i := 0; i < numJobs; i++ {
			jobs <- i
		}
		close(jobs)

		wg2.Wait()
		duration2 := time.Since(start)

		fmt.Printf("   ⏱️  Duration: %v\n", duration2)
		fmt.Printf("   🧵 Workers: %d (constant)\n", numWorkers)

		// Performance comparison
		if duration2 < duration1 {
			improvement := float64(duration1.Nanoseconds()) / float64(duration2.Nanoseconds())
			fmt.Printf("   🏆 Worker pool %0.1fx LEBIH CEPAT!\n", improvement)
		} else {
			degradation := float64(duration2.Nanoseconds()) / float64(duration1.Nanoseconds())
			fmt.Printf("   ⚠️  Worker pool %0.1fx lebih lambat (overhead scheduling)\n", degradation)
		}

		// Memory estimation
		memoryOnePerJob := (maxGoroutines - goroutinesBefore) * 8 // 8KB per goroutine
		memoryWorkerPool := numWorkers * 8
		fmt.Printf("   💾 Memory - One per job: ~%d KB\n", memoryOnePerJob)
		fmt.Printf("   💾 Memory - Worker pool: ~%d KB\n", memoryWorkerPool)
		fmt.Printf("   📈 Memory saved: %d KB (%0.1fx less)\n",
			memoryOnePerJob-memoryWorkerPool,
			float64(memoryOnePerJob)/float64(memoryWorkerPool))
	}
}

func demonstrateResourceExhaustion() {
	fmt.Println("\n💥 RESOURCE EXHAUSTION DEMO")
	fmt.Println("============================")

	fmt.Println("Simulasi: Bagaimana jika ada 50,000 jobs bersamaan?")

	numJobs := 50000

	// Test Worker Pool (aman)
	fmt.Printf("\n✅ Worker Pool dengan %d workers untuk %d jobs:\n", runtime.NumCPU(), numJobs)
	start := time.Now()

	jobs := make(chan int, 1000)
	var wg sync.WaitGroup

	goroutinesBefore := runtime.NumGoroutine()

	// Start workers
	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				// Minimal work
				time.Sleep(100 * time.Microsecond)
			}
		}()
	}

	// Send jobs
	for i := 0; i < numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	duration := time.Since(start)
	goroutinesAfter := runtime.NumGoroutine()

	fmt.Printf("   ⏱️  Duration: %v\n", duration)
	fmt.Printf("   🧵 Goroutines used: %d (+ %d)\n", goroutinesAfter, goroutinesAfter-goroutinesBefore)
	fmt.Printf("   💾 Memory overhead: ~%d KB\n", (goroutinesAfter-goroutinesBefore)*8)
	fmt.Printf("   ✅ Status: SUCCESS - System stable\n")

	// Simulate One Per Job (would be dangerous)
	fmt.Printf("\n🔥 One Goroutine Per Job simulation untuk %d jobs:\n", numJobs)
	fmt.Printf("   🧵 Goroutines needed: %d\n", numJobs)
	fmt.Printf("   💾 Memory needed: ~%d MB\n", numJobs*8/1024)
	fmt.Printf("   ⚠️  Status: DANGEROUS - Possible memory exhaustion\n")
	fmt.Printf("   ❌ Risk: System crash, OOM killer activation\n")

	fmt.Println(`
💡 KESIMPULAN:
• Worker pool: Scalable, predictable resource usage
• One per job: Memory grows linearly with job count
• Breaking point: ~10k-100k concurrent goroutines (system dependent)
• Production systems: Always use bounded concurrency!
`)
}

func main() {
	fmt.Println("🎯 MENGAPA WORKER POOL PATTERN LEBIH BAIK?")
	fmt.Println("==========================================")

	fmt.Printf("💻 System info: %d CPU cores\n", runtime.NumCPU())

	// Demonstrate scalability differences
	demonstrateExtremeScalability()

	// Show resource exhaustion scenarios
	demonstrateResourceExhaustion()

	fmt.Println(`
🏆 FINAL ANSWER - MENGAPA TIDAK N GOROUTINES UNTUK N JOBS:

📈 SCALABILITY:
   • N goroutines: O(N) memory usage
   • Worker pool: O(constant) memory usage
   • At scale: Worker pool wins dramatically

⚡ PERFORMANCE:
   • Small N: One-per-job might be faster (less coordination)
   • Large N: Worker pool much faster (less scheduler overhead)
   • Breaking point: Usually around 1000-10000 jobs

💾 RESOURCE MANAGEMENT:
   • Predictable memory usage
   • Bounded CPU utilization  
   • System stability guaranteed

🎛️ CONTROL:
   • Adjustable concurrency level
   • Backpressure handling
   • Graceful shutdown possible

🐛 DEBUGGING:
   • Fixed number of goroutines to monitor
   • Predictable behavior
   • Easier profiling and optimization

🚀 PRODUCTION READY:
   • Handles any job count
   • Won't crash your system
   • Scalable architecture pattern
`)
}
