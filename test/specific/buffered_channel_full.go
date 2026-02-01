package main

import (
	"fmt"
	"time"
)

// =============================================================================
// DEMO: BUFFERED CHANNEL FULL & BLOCKING
// =============================================================================
// Skenario:
// 1. Buffer size hanya 3.
// 2. Producer mengirim 10 item dengan sangat cepat.
// 3. Worker (Consumer) sangat lambat (1 detik per item).
//
// Kita akan melihat Producer "Macet" (Blocking) setelah mengirim 3 item,
// dan baru bisa kirim lagi setelah Worker mengambil item.
// =============================================================================

func main() {
	fmt.Println("🚦 DEMO BUFFERED CHANNEL BLOCKING")

	// Buffer hanya muat 3 surat
	bufferSize := 3
	jobs := make(chan int, bufferSize)

	// PRODUCER (Pengirim Cepat)
	go func() {
		for i := 1; i <= 10; i++ {
			fmt.Printf("📤 Producer mau kirim item #%d... ", i)

			// Baris di bawah ini akan BLOCKING jika buffer penuh (3 item)
			jobs <- i

			fmt.Printf("✅ BERHASIL Masuk Buffer (Len: %d/%d)\n", len(jobs), cap(jobs))
		}
		close(jobs)
		fmt.Println("🏁 Producer selesai tugas & pulang.")
	}()

	// // Beri jeda sebentar biar kelihatan Producer ngisi queue dulu
	// time.Sleep(2 * time.Second)
	// fmt.Println("\n🚧 PERHATIKAN: Producer macet di item #4 karena Buffer PENUH (3/3)!\nSekarang Worker bangun...\n")

	// WORKER (Penerima Lambat)
	consumerStart := time.Now()
	for job := range jobs {
		fmt.Printf("   📥 Worker ambil item #%d (Sedang diproses...)\n", job)
		time.Sleep(1 * time.Second) // Simulasi lelet
	}

	fmt.Printf("\n⏱️  Selesai dalam: %v\n", time.Since(consumerStart))
}
