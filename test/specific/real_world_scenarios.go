package main

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// =============================================================================
// REAL WORLD SCENARIOS (REAL CODE EDITION)
// =============================================================================
// Uncomment fungsi di main() satu per satu untuk melihat perbedaannya.
// ⚠️ PERIGATAN: Kode ini melakukan request HTTP beneran. Pastikan ada internet.

func main() {
	fmt.Println("🚀 START REAL WORLD SIMULATION...")

	// Skenario 1: IO Bound (HTTP Request ke Server Public)
	Case1_IOBound_HttpRequests()

	// Skenario 2: CPU Bound (Kalkulasi Fibonacci Berat)
	// Case2_CPUBound_Fibonacci()
}

// -----------------------------------------------------------------------------
// CASE 1: IO BOUND (HTTP REQUEST)
// Strategi: UNLIMITED GOROUTINES
// -----------------------------------------------------------------------------
func Case1_IOBound_HttpRequests() {
	intro("CASE 1: 5 Concurrent HTTP Requests (Real IO Bound)")

	// Kita batasi 5 saja agar tidak membebani server orang lain (tapi logikanya bisa 10.000)
	count := 5
	var wg sync.WaitGroup
	wg.Add(count)

	start := time.Now()

	client := &http.Client{Timeout: 10 * time.Second}

	for i := 1; i <= count; i++ {
		go func(id int) {
			defer wg.Done()

			// [REAL CODE]
			// Ini adalah IO Bound murni.
			// Saat http.Get dipanggil, Goroutine ini "TIDUR" menunggu response dari server.
			// CPU komputer Anda 0% usage untuk goroutine ini.

			// url := "https://api.github.com" // (Bisa diganti URL apapun)
			// KITA REQUEST KE MENU YANG BERBEDA (Todo 1, Todo 2, dst)
			url := fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", id)

			resp, err := client.Get(url)
			if err != nil {
				fmt.Printf("❌ Request %d Gagal: %v\n", id, err)
				return
			}

			// Wajib close body agar koneksi bisa dipakai ulang (Reuse)
			defer resp.Body.Close()

			// Baca body (Sedikit CPU work untuk copy data dari network card ke memory)
			_, _ = io.Copy(io.Discard, resp.Body)

			fmt.Printf("✅ Request %d Selesai (Status: %s)\n", id, resp.Status)
		}(i)
	}

	wg.Wait()
	fmt.Printf("\n⏱️  Total Waktu: %v\n", time.Since(start))
	fmt.Println("👉 Perhatikan: Waktu total hampir sama dengan waktu 1 request (Concurrent).")
}

// -----------------------------------------------------------------------------
// CASE 2: CPU BOUND (FIBONACCI)
// Strategi: WORKER POOL
// -----------------------------------------------------------------------------
func Case2_CPUBound_Fibonacci() {
	intro("CASE 2: Kalkulasi Fibonacci (Real CPU Bound)")

	// Fibonacci angka 40 itu BERAT sekali. Butuh jutaan rekursi.
	// Jika kita jalankan ini tanpa worker pool (misal 50 goroutine hitung Fib(40)),
	// Laptop akan hang.

	jobs := []int{35, 36, 37, 38, 39, 40, 41, 42}
	totalJobs := len(jobs)

	// Limit Worker sesuai jumlah Otak (Core)
	numWorkers := runtime.NumCPU()
	fmt.Printf("🖥️  CPU Cores: %d. Workers: %d\n", numWorkers, numWorkers)

	jobChan := make(chan int, totalJobs)
	results := make(chan int, totalJobs)

	// 1. Spawn Workers
	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			for n := range jobChan {
				fmt.Printf("👷 Worker %d mulai hitung Fib(%d)...\n", id, n)
				start := time.Now()

				res := recursiveFib(n) // <--- CPU NGEBUT DI SINI

				fmt.Printf("✅ Worker %d selesai Fib(%d) = %d (%v)\n", id, n, res, time.Since(start))
				results <- res
			}
		}(w)
	}

	// 2. Kirim Job
	for _, n := range jobs {
		jobChan <- n
	}
	close(jobChan)

	// 3. Tunggu Hasil
	for i := 0; i < totalJobs; i++ {
		<-results
	}

	fmt.Println("\n✅ Semua Selesai.")
}

// Fungsi Rekursif (Sangat boros CPU & Stack)
func recursiveFib(n int) int {
	if n <= 1 {
		return n
	}
	return recursiveFib(n-1) + recursiveFib(n-2)
}

func intro(title string) {
	fmt.Println("\n------------------------------------------------")
	fmt.Println(title)
	fmt.Println("------------------------------------------------")
}
