package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// =============================================================================
// PATTERN: FAN-OUT & FAN-IN (LOAD BALANCER)
// =============================================================================
// Skenario: "Image Processing Service"
// Kita punya 100 Gambar yang harus di-resize.
// Jika dikerjakan sendiri (Sequential) = Lama.
// Kita sebar ke 5 Worker (Fan-Out) -> Gabungkan hasilnya (Fan-In).
// =============================================================================

type Job struct {
	ID        int
	ImageName string
}

type Result struct {
	JobID       int
	ImageName   string
	OutputSize  string
	ProcessedBy string
	Duration    time.Duration
}

func main() {
	fmt.Println("🏭 STARTING IMAGE PROCESSING SERVICE...")

	const totalImages = 15
	const totalWorkers = 3 // Kita pakai 3 worker untuk bantu proses

	start := time.Now()

	// 1. INPUT CHANNEL (Pipeline)
	jobs := make(chan Job, totalImages)
	results := make(chan Result, totalImages)

	// 2. FAN-OUT (Spawn Workers)
	// Kita buat 3 goroutine yang mendengarkan channel 'jobs' yang SAMA.
	// Mereka akan balapan mengambil job (Competing Consumers).
	var wg sync.WaitGroup
	for w := 1; w <= totalWorkers; w++ {
		wg.Add(1)
		go imageWorker(w, jobs, results, &wg)
	}

	// 3. GENERATE JOBS (Producer)
	// Kirim pekerjaan ke channel
	for i := 1; i <= totalImages; i++ {
		jobs <- Job{ID: i, ImageName: fmt.Sprintf("image_%d.jpg", i)}
	}
	close(jobs) // Tanda pekerjaan habis

	// 4. FAN-IN (Wait & Close)
	// Kita butuh goroutine terpisah untuk menutup channel results setelah semua worker selesai
	go func() {
		wg.Wait()      // Tunggu semua worker pulang
		close(results) // Tutup pintu hasil
	}()

	// 5. CONSUMER (Process Results)
	// Membaca laporan hasil kerjaan
	totalSuccess := 0
	for res := range results {
		totalSuccess++
		fmt.Printf("✅ [%s] %s -> %s (%v)\n",
			res.ProcessedBy, res.ImageName, res.OutputSize, res.Duration)
	}

	fmt.Printf("\n⏱️  Total Waktu: %v (Harusnya ~%d detik / %d workers)\n",
		time.Since(start), totalImages, totalWorkers)
}

// WORKER (Konsumen yang rajin)
func imageWorker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	workerName := fmt.Sprintf("Worker-%d", id)

	for job := range jobs {
		// Simulasi proses berat (CPU Bound atau IO Bound)
		// Anggap ini resize gambar yang butuh waktu random 500ms - 1s
		start := time.Now()
		sleepTime := time.Duration(rand.Intn(500)+500) * time.Millisecond
		time.Sleep(sleepTime)

		// Kirim hasil ke channel output
		results <- Result{
			JobID:       job.ID,
			ImageName:   job.ImageName,
			OutputSize:  "1920x1080",
			ProcessedBy: workerName,
			Duration:    time.Since(start),
		}
	}
}
