package main

import (
	"fmt"
	"time"
)

// =============================================================================
// DEMO: CLOSURE VS PARAMETER DI GOROUTINE
// =============================================================================
// Pertanyaan: "Apa perlu dipassing parameter atau bisa langsung akses?"
// Jawaban:
// 1. Channel (`out`) -> AMAN diakses langsung (Closure).
// 2. Loop Variable (`i`) -> HATI-HATI (Terutama di Go versi < 1.22).
// =============================================================================

func main() {
	fmt.Println("🕵️ DEMO: CLOSURE TRAP")

	// KAPASITAS CHANNEL 0 (Unbuffered)
	// Kita pakai unbuffered biar kelihatan efek locking/blocking-nya kalau ada.

	// KASUS 1: Channel (Safe Closure)
	// Kita tidak perlu passing 'ch' sebagai parameter.
	// Goroutine bisa 'melihat' variabel 'ch' milik 'main'.
	ch := make(chan string)
	go func() {
		// Mengakses 'ch' dari outer scope secara langsung
		ch <- "Halo dari Closure!"
	}()
	fmt.Println(<-ch)

	// KASUS 2: Loop Variable (The Classic Trap)
	// Coba jalankan ini. Di Go lama (<= 1.21), outputnya mungkin "5, 5, 5, 5, 5"
	// Di Go baru (>= 1.22), outputnya sudah benar "0, 1, 2, 3, 4" (Fixed!)

	fmt.Println("\n--- Tes Loop Variable Capture ---")
	for i := 0; i < 3; i++ {
		go func() {
			// Mengakses 'i' langsung adalah resiko di Go versi lama!
			// Karena saat goroutine ini jalan, loop 'i' mungkin sudah berubah.
			fmt.Printf("Unsafe Access: %d\n", i)
		}()
	}

	time.Sleep(100 * time.Millisecond) // Tunggu goroutines

	// KASUS 3: Solusi Cara Aman (Parameter Passing)
	// Ini cara paling robust, berlaku untuk SEMUA versi Go.
	// Kita "mengunci" nilai 'i' dengan mengirimnya sebagai parameter.

	fmt.Println("\n--- Tes Parameter Passing (Safe) ---")
	for i := 0; i < 3; i++ {
		go func(val int) {
			// 'val' adalah salinan dari 'i' SAAT goroutine DIBUAT.
			// Jadi tidak peduli 'i' di luar berubah jadi berapa.
			fmt.Printf("Safe Parameter: %d\n", val)
		}(i) // <--- Pass 'i' di sini
	}

	time.Sleep(100 * time.Millisecond)
}
