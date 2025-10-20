package main

import (
	"fmt"
	"strings"
	"time"
)

// demonstrateUnbufferedBlocking menunjukkan bagaimana unbuffered channel melakukan blocking
func demonstrateUnbufferedBlocking() {
	fmt.Println("=== UNBUFFERED CHANNEL BLOCKING DEMO ===")

	ch := make(chan string)

	fmt.Println("1. Main goroutine: Memulai goroutine pengirim...")

	// Goroutine pengirim
	go func() {
		fmt.Println("2. Sender goroutine: Akan mengirim data...")
		fmt.Println("3. Sender goroutine: BLOCKING - menunggu penerima...")

		ch <- "Data penting"

		fmt.Println("6. Sender goroutine: Data berhasil dikirim!")
	}()

	// Simulasi main goroutine melakukan pekerjaan lain
	fmt.Println("4. Main goroutine: Melakukan pekerjaan lain selama 2 detik...")
	time.Sleep(2 * time.Second)

	fmt.Println("5. Main goroutine: Siap menerima data...")
	data := <-ch
	fmt.Println("7. Main goroutine: Menerima:", data)

	time.Sleep(100 * time.Millisecond) // Memberi waktu untuk melihat output
}

// demonstrateBufferedBlocking menunjukkan blocking pada buffered channel
func demonstrateBufferedBlocking() {
	fmt.Println("\n=== BUFFERED CHANNEL BLOCKING DEMO ===")

	// Buffer dengan kapasitas 2
	ch := make(chan int, 2)

	fmt.Println("1. Mengirim data ke buffer (tidak blocking)...")

	// Pengiriman pertama - tidak blocking
	ch <- 1
	fmt.Println("2. Data 1 dikirim - buffer: 1/2")

	// Pengiriman kedua - tidak blocking
	ch <- 2
	fmt.Println("3. Data 2 dikirim - buffer: 2/2 (PENUH)")

	// Pengiriman ketiga - akan blocking
	fmt.Println("4. Memulai goroutine untuk pengiriman ketiga...")

	go func() {
		fmt.Println("5. Sender goroutine: BLOCKING - buffer penuh, menunggu penerima...")
		ch <- 3
		fmt.Println("8. Sender goroutine: Data 3 berhasil dikirim!")
	}()

	// Memberi waktu untuk menunjukkan blocking
	time.Sleep(1 * time.Second)

	// Menerima satu data untuk membuat ruang di buffer
	fmt.Println("6. Main goroutine: Menerima data untuk membuat ruang...")
	data := <-ch
	fmt.Println("7. Main goroutine: Menerima:", data, "- buffer sekarang ada ruang")

	// Menerima data sisanya
	fmt.Println("9. Menerima data sisanya...")
	fmt.Println("Data:", <-ch)
	fmt.Println("Data:", <-ch)

	time.Sleep(100 * time.Millisecond)
}

// demonstrateReceiveBlocking menunjukkan blocking pada operasi receive
func demonstrateReceiveBlocking() {
	fmt.Println("\n=== RECEIVE BLOCKING DEMO ===")

	ch := make(chan string)

	fmt.Println("1. Memulai goroutine penerima...")

	go func() {
		fmt.Println("2. Receiver goroutine: BLOCKING - menunggu data...")
		data := <-ch
		fmt.Println("5. Receiver goroutine: Menerima:", data)
	}()

	// Memberi waktu untuk menunjukkan receiver blocking
	fmt.Println("3. Main goroutine: Receiver sedang blocking, tunggu 2 detik...")
	time.Sleep(2 * time.Second)

	// Mengirim data
	fmt.Println("4. Main goroutine: Mengirim data...")
	ch <- "Data terlambat"

	time.Sleep(100 * time.Millisecond)
}

// demonstrateDeadlock menjelaskan skenario deadlock tanpa menjalankannya
func demonstrateDeadlock() {
	fmt.Println("\n=== DEADLOCK SCENARIO (COMMENTED TO AVOID CRASH) ===")

	// JANGAN jalankan kode ini - akan menyebabkan deadlock
	/*
		ch := make(chan string)

		fmt.Println("Mengirim ke unbuffered channel tanpa receiver...")
		ch <- "Ini akan menyebabkan deadlock"
		fmt.Println("Baris ini tidak akan pernah dieksekusi")
	*/

	fmt.Println("Kode deadlock dikomentari untuk keamanan.")
	fmt.Println("Deadlock terjadi karena:")
	fmt.Println("- Unbuffered channel membutuhkan sender dan receiver siap bersamaan")
	fmt.Println("- Tanpa goroutine receiver, sender akan blocking selamanya")
	fmt.Println("- Main goroutine adalah satu-satunya goroutine yang berjalan")
}

// demonstrateTimeoutToAvoidBlocking menunjukkan penggunaan timeout
func demonstrateTimeoutToAvoidBlocking() {
	fmt.Println("\n=== TIMEOUT TO AVOID EXCESSIVE BLOCKING ===")

	ch := make(chan string)

	// Simulasi sender yang lambat
	go func() {
		time.Sleep(3 * time.Second)
		ch <- "Data lambat"
	}()

	fmt.Println("Menunggu data dengan timeout 2 detik...")

	select {
	case data := <-ch:
		fmt.Println("Menerima:", data)
	case <-time.After(2 * time.Second):
		fmt.Println("TIMEOUT: Data tidak diterima dalam 2 detik")
		fmt.Println("Menghindari blocking berlebihan dengan timeout")
	}
}

// demonstrateNonBlocking menunjukkan operasi non-blocking dengan select default
func demonstrateNonBlocking() {
	fmt.Println("\n=== NON-BLOCKING OPERATIONS ===")

	ch := make(chan int, 1)

	// Operasi send non-blocking
	select {
	case ch <- 42:
		fmt.Println("Data berhasil dikirim (non-blocking)")
	default:
		fmt.Println("Channel penuh, tidak bisa mengirim (non-blocking)")
	}

	// Operasi receive non-blocking
	select {
	case data := <-ch:
		fmt.Println("Menerima data:", data, "(non-blocking)")
	default:
		fmt.Println("Tidak ada data tersedia (non-blocking)")
	}

	// Coba receive lagi ketika channel kosong
	select {
	case data := <-ch:
		fmt.Println("Menerima data:", data)
	default:
		fmt.Println("Channel kosong, tidak ada data (non-blocking)")
	}
}

// demonstrateAdvancedBlockingPatterns menunjukkan pola blocking yang lebih kompleks
func demonstrateAdvancedBlockingPatterns() {
	fmt.Println("\n=== ADVANCED BLOCKING PATTERNS ===")

	// Pattern 1: Pipeline dengan blocking natural
	fmt.Println("\n--- Pipeline Pattern ---")
	numbers := make(chan int)
	squares := make(chan int)

	// Stage 1: Generator
	go func() {
		defer close(numbers)
		for i := 1; i <= 5; i++ {
			fmt.Printf("Generator: Mengirim %d\n", i)
			numbers <- i
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Stage 2: Squarer
	go func() {
		defer close(squares)
		for num := range numbers {
			squared := num * num
			fmt.Printf("Squarer: %d -> %d\n", num, squared)
			squares <- squared
		}
	}()

	// Stage 3: Consumer
	fmt.Println("Consumer: Memulai konsumsi...")
	for square := range squares {
		fmt.Printf("Consumer: Menerima %d\n", square)
	}

	// Pattern 2: Worker Pool dengan bounded channel
	fmt.Println("\n--- Worker Pool Pattern ---")
	jobs := make(chan int, 3) // Buffered untuk batching
	results := make(chan int, 3)

	// Worker
	go func() {
		for job := range jobs {
			fmt.Printf("Worker: Memproses job %d\n", job)
			time.Sleep(200 * time.Millisecond)
			results <- job * 2
		}
	}()

	// Send jobs (akan blocking ketika buffer penuh)
	go func() {
		defer close(jobs)
		for i := 1; i <= 5; i++ {
			fmt.Printf("Sender: Mengirim job %d\n", i)
			jobs <- i
			fmt.Printf("Sender: Job %d berhasil dikirim\n", i)
		}
	}()

	// Receive results
	for i := 1; i <= 5; i++ {
		result := <-results
		fmt.Printf("Main: Menerima hasil %d\n", result)
	}
}

// demonstrateBlockingVisualization memberikan visualisasi timing blocking
func demonstrateBlockingVisualization() {
	fmt.Println("\n=== BLOCKING TIMING VISUALIZATION ===")

	ch := make(chan string)

	start := time.Now()

	go func() {
		time.Sleep(1 * time.Second)
		elapsed := time.Since(start)
		fmt.Printf("[%v] Sender: Mengirim data\n", elapsed.Round(time.Millisecond))
		ch <- "Data"
		elapsed = time.Since(start)
		fmt.Printf("[%v] Sender: Data berhasil dikirim\n", elapsed.Round(time.Millisecond))
	}()

	// Main goroutine akan blocking sampai data tersedia
	elapsed := time.Since(start)
	fmt.Printf("[%v] Main: Mulai menunggu data (akan blocking)\n", elapsed.Round(time.Millisecond))

	data := <-ch

	elapsed = time.Since(start)
	fmt.Printf("[%v] Main: Menerima data: %s\n", elapsed.Round(time.Millisecond), data)
}

// demonstrateChannelReceivePatterns menjelaskan perbedaan range vs direct receive
func demonstrateChannelReceivePatterns() {
	fmt.Println("\n=== CHANNEL RECEIVE PATTERNS ===")

	// Pattern 1: Range Loop (otomatis loop sampai channel ditutup)
	fmt.Println("\n--- Pattern 1: Range Loop ---")
	jobs1 := make(chan int)

	go func() {
		defer close(jobs1) // PENTING: close channel untuk menghentikan range loop
		for i := 1; i <= 3; i++ {
			fmt.Printf("Sender: Mengirim job %d\n", i)
			jobs1 <- i
			time.Sleep(200 * time.Millisecond)
		}
		fmt.Println("Sender: Semua job terkirim, closing channel...")
	}()

	fmt.Println("Range Receiver: Memulai range loop...")
	for job := range jobs1 { // OTOMATIS LOOP sampai channel ditutup
		fmt.Printf("Range Receiver: Menerima job %d\n", job)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Range Receiver: Channel ditutup, keluar dari loop")

	// Pattern 2: Direct Receive (manual one-by-one)
	fmt.Println("\n--- Pattern 2: Direct Receive ---")
	jobs2 := make(chan int, 3) // Buffered untuk demo

	// Kirim beberapa jobs sekaligus
	go func() {
		for i := 1; i <= 3; i++ {
			fmt.Printf("Sender: Mengirim job %d\n", i)
			jobs2 <- i
		}
		close(jobs2)
		fmt.Println("Sender: Semua job terkirim, channel ditutup")
	}()

	time.Sleep(100 * time.Millisecond) // Beri waktu sender selesai

	fmt.Println("Direct Receiver: Menerima satu per satu...")

	// Cara 1: Manual receive dengan loop
	for i := 1; i <= 3; i++ {
		job := <-jobs2 // MANUAL RECEIVE
		fmt.Printf("Direct Receiver: Menerima job %d (receive ke-%d)\n", job, i)
	}

	// Pattern 3: Direct Receive dengan Check Channel Status
	fmt.Println("\n--- Pattern 3: Direct Receive + Status Check ---")
	jobs3 := make(chan string, 2)

	go func() {
		jobs3 <- "TaskA"
		jobs3 <- "TaskB"
		close(jobs3)
		fmt.Println("Sender: Channel ditutup")
	}()

	time.Sleep(50 * time.Millisecond)

	fmt.Println("Status Checker: Receiving dengan status check...")
	for {
		job, ok := <-jobs3 // RECEIVE + STATUS CHECK
		if !ok {
			fmt.Println("Status Checker: Channel ditutup, tidak ada data lagi")
			break
		}
		fmt.Printf("Status Checker: Menerima %s, channel masih open\n", job)
	}

	// Pattern 4: Demonstrasi dalam Worker Pool Context
	fmt.Println("\n--- Pattern 4: Worker Pool Context ---")
	fmt.Println("Pertanyaan Anda: Apakah `for job := range jobs` sama dengan `job := <-jobs`?")

	jobs4 := make(chan int)

	// Worker dengan RANGE LOOP
	go func() {
		fmt.Println("Worker-Range: Mulai dengan range loop")
		for job := range jobs4 { // INI LOOP OTOMATIS
			fmt.Printf("Worker-Range: Processing job %d\n", job)
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println("Worker-Range: Channel ditutup, worker berhenti")
	}()

	// Jika kita gunakan DIRECT RECEIVE, harus manual loop
	jobs5 := make(chan int)
	go func() {
		fmt.Println("Worker-Direct: Mulai dengan direct receive")
		for { // HARUS MANUAL LOOP
			job := <-jobs5 // MANUAL RECEIVE
			fmt.Printf("Worker-Direct: Processing job %d\n", job)
			time.Sleep(100 * time.Millisecond)

			// MASALAH: Bagaimana tahu kapan berhenti?
			// Range loop otomatis berhenti saat channel closed
			// Direct receive tidak tahu kapan channel closed
		}
	}()

	// Kirim jobs ke kedua worker
	for i := 1; i <= 2; i++ {
		fmt.Printf("Main: Kirim job %d ke range worker\n", i)
		jobs4 <- i
		time.Sleep(150 * time.Millisecond)
	}
	close(jobs4)

	for i := 1; i <= 2; i++ {
		fmt.Printf("Main: Kirim job %d ke direct worker\n", i)
		jobs5 <- i
		time.Sleep(150 * time.Millisecond)
	}
	close(jobs5) // Direct worker tidak akan tahu channel ditutup!

	time.Sleep(200 * time.Millisecond)

	fmt.Println(`
💡 JAWABAN PERTANYAAN:
┌─────────────────────────────────────────────────────────────┐
│ TIDAK, keduanya BERBEDA:                                    │
│                                                             │
│ 1. "for job := range jobs"                                  │
│    ✅ OTOMATIS LOOP sampai channel ditutup                 │
│    ✅ OTOMATIS CHECK apakah channel closed                 │
│    ✅ OTOMATIS KELUAR dari loop saat channel closed        │
│    ✅ IDIOMATIK untuk worker pattern                       │
│                                                             │
│ 2. "job := <-jobs"                                          │
│    ❌ HANYA SATU KALI receive                              │
│    ❌ TIDAK otomatis loop                                   │
│    ❌ TIDAK otomatis cek channel status                     │
│    ❌ Perlu manual loop dan manual check                    │
│                                                             │
│ KESIMPULAN: Range loop adalah HIGH-LEVEL abstraction       │
│ yang SANGAT COCOK untuk worker patterns!                   │
└─────────────────────────────────────────────────────────────┘

🔍 DETAIL TEKNIS:
• Range loop: for job := range jobs { ... }
  Equivalent dengan:
  for {
      job, ok := <-jobs
      if !ok { break }
      // process job
  }

• Direct receive: job := <-jobs
  Hanya receive SEKALI, tidak loop, tidak check status
`)
}

func main() {
	fmt.Println("STUDI KASUS: KONSEP LOCKING/BLOCKING PADA CHANNEL")
	fmt.Println(strings.Repeat("=", 60))

	demonstrateUnbufferedBlocking()
	demonstrateBufferedBlocking()
	demonstrateReceiveBlocking()
	demonstrateDeadlock()
	demonstrateTimeoutToAvoidBlocking()
	demonstrateNonBlocking()
	demonstrateAdvancedBlockingPatterns()
	demonstrateBlockingVisualization()
	demonstrateChannelReceivePatterns()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("RINGKASAN KONSEP BLOCKING:")
	fmt.Println("1. Unbuffered channel: Send dan receive selalu blocking sampai ada partner")
	fmt.Println("2. Buffered channel: Send blocking saat buffer penuh, receive blocking saat buffer kosong")
	fmt.Println("3. Gunakan select dengan timeout untuk menghindari blocking berlebihan")
	fmt.Println("4. Gunakan select dengan default untuk operasi non-blocking")
	fmt.Println("5. Hindari deadlock dengan memastikan ada goroutine sender/receiver yang sesuai")
	fmt.Println("6. Blocking adalah mekanisme sinkronisasi natural - bukan bug tapi fitur!")
	fmt.Println("7. Pipeline dan worker pool memanfaatkan blocking untuk flow control")
}
