package main

import (
	"fmt"
	"time"
)

// =============================================================================
// PEMBUKTIAN: PIPELINE ADALAH STREAMING (BUKAN BATCHING)
// =============================================================================
// User Question: "Apakah stage 2 menunggu semua value terkumpul dulu?"
// Answer: TIDAK. Data mengalir satu per satu (estafet).
//
// Mari kita buktikan dengan log waktu.
// Jika Batching: Kita akan melihat semua "Filter" selesai, baru "Process" mulai.
// Jika Streaming: Kita akan melihat "Filter" dan "Process" muncul bergantian.
// =============================================================================

func main() {
	fmt.Println("🌊 DEMO: STREAMING PROOF")
	fmt.Println("========================")

	// Stage 1: Generator
	// Menghasilkan 5 angka
	nums := generateStream(5)

	// Stage 2: Filter
	// Menambahkan delay 500ms
	filtered := filterStream(nums)

	// Stage 3: Consumer
	// Langsung cetak begitu data sampai
	for val := range filtered {
		fmt.Printf("   🏁 CONSUMER: Terima data %d (Saat ini: %s)\n",
			val, time.Now().Format("15:04:05.000"))
	}
}

func generateStream(max int) <-chan int {
	// 1.buat channel
	out := make(chan int)
	// 2.buat goroutine
	go func() {
		// 3.buat loop
		for i := 1; i <= max; i++ {
			fmt.Printf("1️⃣  GENERATOR: Buat %d\n", i)
			// 4.kirim ke channel
			out <- i // Kirim ke stage berikutnya (akan block jika stage 2 belum siap terima)
			time.Sleep(200 * time.Millisecond)
		} 
		// 5.lakukan sampai selesai lalu tutup channel , di asumsi saya out masih berisi 5 data/sesuai max
		close(out)
		fmt.Println("✅ GENERATOR SELESAI")
	}()
	// 6. return channel , nah dipikiran saya masih ada 5 data di dalam channel out
	return out
}

func filterStream(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			// Simulasi proses "Mikir"
			fmt.Printf("   2️⃣  FILTER   : Sedang memproses %d... \n", n)
			time.Sleep(500 * time.Millisecond)

			// Kirim ke consumer
			fmt.Printf("   2️⃣  FILTER   : Kirim %d ke Consumer\n", n)
			out <- n
		}
		close(out)
		fmt.Println("✅ FILTER SELESAI")
	}()
	return out
}
