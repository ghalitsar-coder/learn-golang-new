package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// APPROACH 1: One Goroutine Per Job (TIDAK DIREKOMENDASIKAN)
func oneGoroutinePerJob(numJobs int) {
	fmt.Printf("\n🔥 APPROACH 1: One Goroutine Per Job (%d jobs = %d goroutines)\n", numJobs, numJobs)
	fmt.Println("=" + fmt.Sprintf("%60s", "="))

	start := time.Now()
	var wg sync.WaitGroup

	// Catat goroutine sebelum
	goroutinesBefore := runtime.NumGoroutine()

	for i := 0; i < numJobs; i++ {
		wg.Add(1)
		go func(jobID int) {
			defer wg.Done()

			// Simulasi pekerjaan
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Job %d selesai\n", jobID)
		}(i)
	}

	// Catat puncak goroutine
	maxGoroutines := runtime.NumGoroutine()

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("\n📊 RESULTS:\n")
	fmt.Printf("   ⏱️  Duration: %v\n", duration)
	fmt.Printf("   🧵 Goroutines before: %d\n", goroutinesBefore)
	fmt.Printf("   🧵 Peak goroutines: %d\n", maxGoroutines)
	fmt.Printf("   💾 Memory overhead: ~%d KB (estimasi)\n", (maxGoroutines-goroutinesBefore)*8)
}

// APPROACH 2: Worker Pool Pattern (DIREKOMENDASIKAN)
func workerPoolPattern(numJobs int, numWorkers int) {
	fmt.Printf("\n✅ APPROACH 2: Worker Pool (%d jobs dengan %d workers)\n", numJobs, numWorkers)
	fmt.Println("=" + fmt.Sprintf("%60s", "="))

	start := time.Now()

	// Channel untuk jobs
	jobs := make(chan int, numJobs)
	results := make(chan string, numJobs)

	var wg sync.WaitGroup

	// Catat goroutine sebelum
	goroutinesBefore := runtime.NumGoroutine()

	// Start workers (fixed number)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for job := range jobs {
				// Simulasi pekerjaan
				time.Sleep(100 * time.Millisecond)
				result := fmt.Sprintf("Worker-%d completed job-%d", workerID, job)
				results <- result
			}
		}(w)
	}

	// Send jobs
	for i := 0; i < numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	// Catat puncak goroutine
	maxGoroutines := runtime.NumGoroutine()

	// Wait for workers to complete
	wg.Wait()
	close(results)

	// Collect results
	var resultCount int
	for result := range results {
		_ = result // Process result (simplified)
		resultCount++
	}

	duration := time.Since(start)

	fmt.Printf("\n📊 RESULTS:\n")
	fmt.Printf("   ⏱️  Duration: %v\n", duration)
	fmt.Printf("   🧵 Goroutines before: %d\n", goroutinesBefore)
	fmt.Printf("   🧵 Peak goroutines: %d\n", maxGoroutines)
	fmt.Printf("   💾 Memory overhead: ~%d KB (estimasi)\n", (maxGoroutines-goroutinesBefore)*8)
	fmt.Printf("   ✅ Jobs completed: %d\n", resultCount)
}

// DEMO: Memory dan Performance Impact
func demonstrateScalabilityIssues() {
	fmt.Println("\n🚨 DEMO: SCALABILITY ISSUES dengan Many Goroutines")
	fmt.Println("=" + fmt.Sprintf("%60s", "="))

	jobs := []int{10, 100, 1000, 5000}

	for _, numJobs := range jobs {
		fmt.Printf("\n--- Testing dengan %d jobs ---\n", numJobs)

		// Test 1: One Goroutine Per Job
		start := time.Now()
		var wg sync.WaitGroup

		goroutinesBefore := runtime.NumGoroutine()

		for i := 0; i < numJobs; i++ {
			wg.Add(1)
			go func(jobID int) {
				defer wg.Done()
				time.Sleep(10 * time.Millisecond) // Pekerjaan ringan
			}(i)
		}

		maxGoroutines := runtime.NumGoroutine()
		wg.Wait()
		duration1 := time.Since(start)

		// Test 2: Worker Pool (optimal workers = CPU cores)
		numWorkers := runtime.NumCPU()
		start = time.Now()

		jobs := make(chan int, numJobs)
		var wg2 sync.WaitGroup

		// Start workers
		for w := 0; w < numWorkers; w++ {
			wg2.Add(1)
			go func() {
				defer wg2.Done()
				for range jobs {
					time.Sleep(10 * time.Millisecond)
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

		fmt.Printf("🔥 One-per-job: %v, %d goroutines\n", duration1, maxGoroutines-goroutinesBefore)
		fmt.Printf("✅ Worker pool: %v, %d workers\n", duration2, numWorkers)

		// Calculate efficiency
		efficiency := float64(duration1.Nanoseconds()) / float64(duration2.Nanoseconds())
		if duration2 < duration1 {
			fmt.Printf("💡 Worker pool %0.1fx lebih efisien!\n", efficiency)
		} else if duration1 < duration2 {
			fmt.Printf("⚠️  One-per-job %0.1fx lebih cepat (overhead kecil)\n", 1/efficiency)
		}
	}
}

// DEMO: Resource Management
func demonstrateResourceManagement() {
	fmt.Println("\n💾 DEMO: RESOURCE MANAGEMENT")
	fmt.Println("=" + fmt.Sprintf("%50s", "="))

	numJobs := 1000

	fmt.Println("\n1️⃣ Sebelum membuat goroutines:")
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	goroutines1 := runtime.NumGoroutine()
	fmt.Printf("   Memory: %d KB, Goroutines: %d\n", m1.Alloc/1024, goroutines1)

	// Create many goroutines
	var wg sync.WaitGroup
	for i := 0; i < numJobs; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
		}(i)
	}

	fmt.Println("\n2️⃣ Setelah membuat 1000 goroutines:")
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	goroutines2 := runtime.NumGoroutine()
	fmt.Printf("   Memory: %d KB, Goroutines: %d\n", m2.Alloc/1024, goroutines2)
	fmt.Printf("   📈 Memory increase: %d KB\n", (m2.Alloc-m1.Alloc)/1024)
	fmt.Printf("   📈 Goroutines increase: %d\n", goroutines2-goroutines1)

	wg.Wait()

	fmt.Println("\n3️⃣ Setelah semua goroutines selesai:")
	var m3 runtime.MemStats
	runtime.ReadMemStats(&m3)
	goroutines3 := runtime.NumGoroutine()
	fmt.Printf("   Memory: %d KB, Goroutines: %d\n", m3.Alloc/1024, goroutines3)

	fmt.Println(`
💡 OBSERVATIONS:
• Setiap goroutine menggunakan ~8KB stack memory
• 1000 goroutines = ~8MB overhead minimal
• Scheduler overhead meningkat dengan jumlah goroutines
• Context switching cost meningkat drastis
`)
}

// DEMO: Optimal Worker Count
func demonstrateOptimalWorkerCount() {
	fmt.Println("\n🎯 DEMO: MENCARI OPTIMAL WORKER COUNT")
	fmt.Println("=" + fmt.Sprintf("%50s", "="))

	numJobs := 1000
	workersToTest := []int{1, 2, 4, 8, 16, 32, 64, 128}

	fmt.Printf("Testing dengan %d jobs, berbagai jumlah workers:\n", numJobs)
	fmt.Println("Workers | Duration | Efficiency")
	fmt.Println("--------|----------|----------")

	var bestDuration time.Duration = time.Hour
	var bestWorkers int

	for _, numWorkers := range workersToTest {
		start := time.Now()

		jobs := make(chan int, numJobs)
		var wg sync.WaitGroup

		// Start workers
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for range jobs {
					// Simulasi CPU-bound work
					time.Sleep(5 * time.Millisecond)
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

		efficiency := "Normal"
		if duration < bestDuration {
			bestDuration = duration
			bestWorkers = numWorkers
			efficiency = "🏆 BEST"
		} else if duration > bestDuration*2 {
			efficiency = "❌ Poor"
		}

		fmt.Printf("%7d | %8v | %s\n", numWorkers, duration.Round(time.Millisecond), efficiency)
	}

	fmt.Printf("\n🎯 OPTIMAL: %d workers dengan duration %v\n", bestWorkers, bestDuration.Round(time.Millisecond))
	fmt.Printf("💡 CPU Cores: %d (reference)\n", runtime.NumCPU())
}

func main() {
	fmt.Println("🏭 WORKER POOL vs ONE GOROUTINE PER JOB")
	fmt.Println("=======================================")

	fmt.Println(`
🤔 PERTANYAAN: 
Mengapa tidak menggunakan N goroutines untuk N jobs?
Mengapa perlu Worker Pool pattern?

Mari kita buktikan dengan benchmark!
`)

	// Demo 1: Perbandingan langsung
	numJobs := 50
	oneGoroutinePerJob(numJobs)
	workerPoolPattern(numJobs, 5)

	// Demo 2: Scalability issues
	demonstrateScalabilityIssues()

	// Demo 3: Resource management
	demonstrateResourceManagement()

	// Demo 4: Optimal worker count
	demonstrateOptimalWorkerCount()

	fmt.Println(`
🎓 KESIMPULAN - MENGAPA WORKER POOL LEBIH BAIK:

1️⃣ MEMORY EFFICIENCY:
   • Goroutine overhead: ~8KB per goroutine
   • 1000 jobs = 1000 goroutines = ~8MB minimum
   • Worker pool: Fixed overhead, scalable

2️⃣ SCHEDULER OVERHEAD:
   • Go scheduler harus manage banyak goroutines
   • Context switching cost meningkat
   • Thrashing pada high concurrency

3️⃣ RESOURCE LIMITS:
   • OS memiliki limit untuk threads/processes
   • Memory fragmentation
   • CPU cache misses

4️⃣ CONTROL & MONITORING:
   • Bounded concurrency (tidak unlimited)
   • Predictable resource usage
   • Better debugging dan profiling

5️⃣ BACKPRESSURE:
   • Natural flow control
   • Prevents system overload
   • Graceful degradation

🚀 BEST PRACTICE:
Workers ≈ CPU cores untuk CPU-bound tasks
Workers ≈ 2-4x CPU cores untuk I/O-bound tasks
`)
}
