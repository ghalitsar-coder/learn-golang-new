package main

import (
	"fmt"
	"time"
)

// =============================================================================
// DEMO: PIPELINE ADALAH "ASSEMBLY LINE" (VERTICAL CONCURRENCY)
// =============================================================================
// Pipeline itu Concordant (Berjalan bersamaan antar stage),
// TAPI Sequential (Berurutan) untuk DATA YANG SAMA.
//
// Analogi Pabrik Mobil:
// - Stage 1: Pasang Rangka
// - Stage 2: Pasang Mesin
// - Stage 3: Cat Mobil
//
// Mobil A sedang dicat, SEMENTARA Mobil B dipasang mesin, SEMENTARA Mobil C dipasang rangka.
// =============================================================================

func main() {
	fmt.Println("🏭 DEMO: PIPELINE CONCURRENCY")

	start := time.Now()

	// Data masuk satu per satu
	input := generate(5)

	// Stage Heavy Processing (Misal: butuh 1 detik per item)
	// Jika sequential murni (tanpa goroutine): 5 item x 1 detik = 5 detik.
	// Jika pipeline concurrent: Stage 1 dan 2 jalan bareng.
	// Total waktu ~= Waktu terlama 1 stage + overhead.

	// Stage 1: Cuci (1 detik)
	stage1 := process("🚿 CUCI ", input, 1000)

	// Stage 2: Bilas (1 detik)
	stage2 := process("🧽 BILAS", stage1, 1000)

	// Stage 3: Kering (1 detik)
	final := process("☀️ KERING", stage2, 1000)

	for res := range final {
		fmt.Printf("   ✅ SELESAI: %v\n", res)
	}

	fmt.Printf("\n⏱️ Total Waktu: %v\n", time.Since(start))
	// Perhatikan: Total waktu bukan 15 detik (5x3), tapi sekitar 6-7 detik.
	// Ini bukti bahwa mereka berjalan CONCURRENTLY (berlapis).
}

func generate(max int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 1; i <= max; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

func process(name string, in <-chan int, delayMs int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			fmt.Printf("%s: Mulai item %d\n", name, n)
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
			fmt.Printf("%s: Selesai item %d\n", name, n)
			out <- n
		}
		close(out)
	}()
	return out
}
