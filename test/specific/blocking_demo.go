package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("🚦 DEMO: BLOCKING VS NON-BLOCKING")
	fmt.Println("================================")

	// Skenario: Kita ingin memonitor sistem sambil menjalankan proses utama.

	// ❌ CARA SALAH (Tanpa Goroutine)
	// Uncomment baris di bawah ini untuk melihat efeknya (Program akan STUCK di sini!)
	// monitorSystem("Blocking Monitor")

	// ✅ CARA BENAR (Dengan Goroutine)
	// Kita wrap dengan 'go' agar dia jalan di "jalur sebelah" (background)
	// dan TIDAK MENGHALANGI baris kode di bawahnya.
	fmt.Println("1. Menjalankan Monitor di Background...")
	go monitorSystem("Background Monitor")

	// Proses Utama tetap bisa jalan!
	fmt.Println("2. Proses UTAMA mulai berjalan...")
	for i := 1; i <= 5; i++ {
		fmt.Printf("   🔨 Proses Utama sedang kerja... (%d/5)\n", i)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("3. Proses UTAMA selesai! Program berhenti (dan monitor otomatis mati).")
}

func monitorSystem(name string) {
	for {
		fmt.Printf("   👀 [%s] Cek Memory...\n", name)
		time.Sleep(500 * time.Millisecond)
	}
}
