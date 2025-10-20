package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// PATTERN 1: wg.Wait() dan close() di MAIN THREAD (seperti kode Anda)
// ✅ AMAN jika buffer results CUKUP BESAR
func pattern1_MainThreadWait() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("PATTERN 1: wg.Wait() + close() di MAIN THREAD")
	fmt.Println(strings.Repeat("=", 60))

	jobs := make(chan int, 5)
	results := make(chan string, 5) // ⚠️ Buffer = jumlah jobs
	var wg sync.WaitGroup

	// Start 2 workers
	for w := 1; w <= 2; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				time.Sleep(100 * time.Millisecond)
				results <- fmt.Sprintf("Worker-%d: job-%d done", id, job)
			}
		}(w)
	}

	// Send 5 jobs
	for i := 1; i <= 5; i++ {
		jobs <- i
	}
	close(jobs)

	fmt.Println("✅ Jobs dikirim, tunggu workers selesai...")

	// KUNCI: Main thread BLOCKING di wg.Wait()
	wg.Wait()
	fmt.Println("✅ Semua worker selesai, tutup results")
	close(results)

	// Baru sekarang baca results
	fmt.Println("📊 Hasil:")
	for r := range results {
		fmt.Println("  ", r)
	}

	fmt.Println(`
💡 KAPAN AMAN:
   • Buffer results >= jumlah total hasil
   • Workers bisa menulis tanpa blocking
   • Main thread menunggu dulu, baru baca

⚠️  RISIKO DEADLOCK JIKA:
   • Buffer results < jumlah hasil
   • Worker blocking saat menulis ke results
   • wg.Wait() tidak akan selesai → DEADLOCK!
`)
}

// PATTERN 2: wg.Wait() + close() di GOROUTINE TERPISAH
// ✅ LEBIH AMAN - tidak bergantung pada buffer size
func pattern2_GoroutineClose() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("PATTERN 2: wg.Wait() + close() di GOROUTINE")
	fmt.Println(strings.Repeat("=", 60))

	jobs := make(chan int, 5)
	results := make(chan string) // ⚠️ Buffer KECIL atau unbuffered
	var wg sync.WaitGroup

	// Start 2 workers
	for w := 1; w <= 2; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				time.Sleep(100 * time.Millisecond)
				results <- fmt.Sprintf("Worker-%d: job-%d done", id, job)
			}
		}(w)
	}

	// Send 5 jobs
	for i := 1; i <= 5; i++ {
		jobs <- i
	}
	close(jobs)

	// KUNCI: Pindahkan wg.Wait() + close() ke GOROUTINE
	go func() {
		wg.Wait()
		fmt.Println("✅ Semua worker selesai, tutup results (dari goroutine)")
		close(results)
	}()

	// Main thread LANGSUNG baca results (tidak tunggu wg dulu)
	fmt.Println("📊 Hasil (dibaca sambil workers masih jalan):")
	for r := range results {
		fmt.Println("  ", r)
	}

	fmt.Println(`
💡 KEUNGGULAN:
   • Main thread langsung consume results
   • Workers tidak pernah blocking saat menulis
   • TIDAK bergantung pada buffer size
   • Pattern PALING AMAN dan IDIOMATIK

✅ RECOMMENDED untuk production code!
`)
}

// PATTERN 3: DEADLOCK SCENARIO (JANGAN LAKUKAN INI!)
func pattern3_DeadlockDemo() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("PATTERN 3: DEADLOCK SCENARIO (DEMO)")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println(`
❌ SKENARIO DEADLOCK:
   1. Buffer results KECIL (< jumlah jobs)
   2. Main thread tunggu wg.Wait() DULU
   3. Workers BLOCKING saat menulis ke results (buffer penuh)
   4. wg.Wait() tidak selesai karena workers stuck
   5. Main tidak pernah sampai ke "baca results"
   6. → DEADLOCK!

Contoh kode yang SALAH:
`)

	fmt.Println(`
func deadlockExample() {
    jobs := make(chan int, 10)
    results := make(chan string, 2) // ❌ Buffer TERLALU KECIL!
    var wg sync.WaitGroup
    
    for w := 1; w <= 3; w++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for job := range jobs {
                // 10 jobs, tapi buffer results cuma 2
                results <- fmt.Sprintf("job-%d done", job)
            }
        }(w)
    }
    
    for i := 1; i <= 10; i++ {
        jobs <- i
    }
    close(jobs)
    
    wg.Wait()        // ❌ STUCK DI SINI SELAMANYA!
    close(results)   // ❌ TIDAK PERNAH SAMPAI SINI
    
    for r := range results {
        fmt.Println(r) // ❌ TIDAK PERNAH DIJALANKAN
    }
}

🔥 YANG TERJADI:
   Step 1: Workers mulai proses jobs
   Step 2: Worker-1 menulis hasil ke-1 → OK (buffer: 1/2)
   Step 3: Worker-2 menulis hasil ke-2 → OK (buffer: 2/2)
   Step 4: Worker-3 mencoba menulis → BLOCKING (buffer PENUH)
   Step 5: Worker-3 stuck, tidak bisa defer wg.Done()
   Step 6: wg.Wait() menunggu Worker-3 selesai
   Step 7: Tidak ada yang baca results untuk kosongkan buffer
   Step 8: → DEADLOCK! Program hang selamanya
`)

	fmt.Println("💡 SOLUSI:")
	fmt.Println("   → Gunakan PATTERN 2 (wg.Wait di goroutine)")
	fmt.Println("   → ATAU pastikan buffer results >= total hasil")
}

// DEMO: Simulasi deadlock yang NYATA (dengan timeout untuk safety)
func pattern4_ActualDeadlockSimulation() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("PATTERN 4: SIMULASI DEADLOCK NYATA (dengan timeout)")
	fmt.Println(strings.Repeat("=", 60))

	done := make(chan bool)

	go func() {
		jobs := make(chan int, 10)
		results := make(chan string, 2) // ❌ Buffer KECIL (2) untuk 10 jobs
		var wg sync.WaitGroup

		// Start 3 workers
		for w := 1; w <= 3; w++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for job := range jobs {
					msg := fmt.Sprintf("Worker-%d: job-%d", id, job)
					fmt.Printf("  🔄 Worker-%d mencoba kirim hasil job-%d...\n", id, job)
					results <- msg // ⚠️ Akan BLOCKING setelah buffer penuh
					fmt.Printf("  ✅ Worker-%d berhasil kirim hasil job-%d\n", id, job)
				}
			}(w)
		}

		// Send 10 jobs
		for i := 1; i <= 10; i++ {
			jobs <- i
		}
		close(jobs)

		fmt.Println("\n⏳ Main goroutine: Tunggu workers selesai...")
		wg.Wait()                             // ❌ AKAN STUCK DI SINI!
		fmt.Println("✅ Semua worker selesai") // TIDAK AKAN DIPRINT
		close(results)
		done <- true
	}()

	// Safety: timeout 2 detik
	select {
	case <-done:
		fmt.Println("✅ Selesai tanpa deadlock")
	case <-time.After(2 * time.Second):
		fmt.Println("\n🚨 DEADLOCK TERDETEKSI!")
		fmt.Println("   Workers stuck menulis ke results channel")
		fmt.Println("   wg.Wait() tidak akan pernah selesai")
		fmt.Println("   Program hang selamanya (kita pakai timeout untuk demo)")
	}
}

func main() {
	fmt.Println("🎯 KAPAN PERLU MEMASUKKAN wg.Wait() + close() KE GOROUTINE?")
	fmt.Println(strings.Repeat("=", 60))

	// Pattern 1: Cara Anda (aman jika buffer cukup)
	pattern1_MainThreadWait()

	time.Sleep(500 * time.Millisecond)

	// Pattern 2: Lebih aman dan idiomatik
	pattern2_GoroutineClose()

	time.Sleep(500 * time.Millisecond)

	// Pattern 3: Penjelasan deadlock
	pattern3_DeadlockDemo()

	time.Sleep(500 * time.Millisecond)

	// Pattern 4: Simulasi deadlock nyata
	pattern4_ActualDeadlockSimulation()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📋 KESIMPULAN")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(`
❓ APAKAH PERLU MEMASUKKAN wg.Wait() + close() KE GOROUTINE?

JAWABAN SINGKAT:
   • Tidak WAJIB, tapi SANGAT DIREKOMENDASIKAN
   • Tergantung strategi consume results

📊 PERBANDINGAN:

┌─────────────────────────────────────────────────────────┐
│ PATTERN 1: wg.Wait() di main thread                    │
├─────────────────────────────────────────────────────────┤
│ ✅ Pros:                                                │
│    • Sederhana, linear flow                             │
│    • Mudah dipahami pemula                              │
│                                                         │
│ ❌ Cons:                                                │
│    • BERGANTUNG pada buffer size                        │
│    • RISIKO DEADLOCK jika buffer < total results       │
│    • Tidak scalable untuk unknown result count         │
│                                                         │
│ 💡 Kapan pakai:                                         │
│    • Jumlah results PASTI dan SEDIKIT                  │
│    • Buffer results >= total results                   │
│    • Code sederhana, bukan production critical         │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│ PATTERN 2: wg.Wait() + close() di goroutine            │
├─────────────────────────────────────────────────────────┤
│ ✅ Pros:                                                │
│    • TIDAK bergantung buffer size                       │
│    • TIDAK ada risiko deadlock                          │
│    • Main langsung consume → workers tidak blocking     │
│    • Scalable untuk any result count                   │
│    • IDIOMATIK Go pattern                               │
│                                                         │
│ ❌ Cons:                                                │
│    • Sedikit lebih kompleks (ada 1 goroutine ekstra)   │
│                                                         │
│ 💡 Kapan pakai:                                         │
│    • Production code                                    │
│    • Jumlah results tidak pasti atau besar             │
│    • Ingin buffer kecil atau unbuffered                │
│    • RECOMMENDED untuk semua kasus!                     │
└─────────────────────────────────────────────────────────┘

🎯 REKOMENDASI PRAKTIS:

1. Untuk learning/demo sederhana:
   → Pattern 1 OK (pastikan buffer cukup)

2. Untuk production code:
   → SELALU gunakan Pattern 2

3. Golden rule:
   → "Jika main thread harus tunggu sebelum consume,
      pindahkan wait+close ke goroutine terpisah"

📝 TEMPLATE PRODUCTION-READY:

func workerPoolProduction() {
    jobs := make(chan Job)
    results := make(chan Result) // unbuffered OK!
    var wg sync.WaitGroup
    
    // Start workers
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- processJob(job)
            }
        }()
    }
    
    // Close results when all workers done
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Main goroutine immediately consumes
    for result := range results {
        handleResult(result)
    }
}

✅ Pattern ini AMAN, SCALABLE, dan IDIOMATIK!
`)
}
