package main

import (
	"fmt"
	"time"
)

// =============================================================================
// KONSEP: PIPELINE PATTERN
// =============================================================================
// Pipeline adalah rangkaian "stage" (tahap) pemrosesan data.
// Output dari Stage 1 menjadi Input untuk Stage 2, dst.
// Mirip ban berjalan di pabrik.
// =============================================================================

func main() {
	fmt.Println("🏭 DEMO: PIPELINE PATTERN")
	fmt.Println("=========================")
	start := time.Now()

	// Stage 1: Generate Angka
	// Menghasilkan sekumpulan angka
	nums := gen(1, 2, 3, 4, 5)

	// Stage 2: Validasi / Filter
	// Hanya loloskan angka genap
	// Input: nums -> Output: filtered
	filtered := filterEven(nums)

	// Stage 3: Proses Berat (Kuadrat)
	// Input: filtered -> Output: squared
	squared := sq(filtered)

	// Stage 4: Final Processing / Consumer
	// Input: squared
	// Kita baca hasil akhirnya di sini
	for n := range squared {
		fmt.Printf("   📦 Hasil Akhir: %d\n", n)
	}

	fmt.Printf("\n⏱️  Selesai dalam: %v\n", time.Since(start))
}

// -----------------------------------------------------------------------------
// STAGES
// -----------------------------------------------------------------------------

// Stage 1: Generator
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			fmt.Printf("1️⃣  Generated: %d\n", n)
			out <- n
			time.Sleep(100 * time.Millisecond) // Simulasi delay
		}
		close(out)
	}()
	return out
}

// Stage 2: Filter (Hanya Genap)
func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if n%2 == 0 {
				fmt.Printf("   2️⃣  Filter OK: %d\n", n)
				out <- n
			} else {
				fmt.Printf("   2️⃣  Filter SKIP: %d\n", n)
			}
		}
		close(out)
	}()
	return out
}

// Stage 3: Squaring (Kuadrat)
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			fmt.Printf("      3️⃣  Squaring: %d -> %d\n", n, n*n)
			time.Sleep(200 * time.Millisecond) // Simulasi proses agak lama
			out <- n * n
		}
		close(out)
	}()
	return out
}
