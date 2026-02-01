package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("========== DEMO: Perilaku Select dalam Loop ==========\n")
	fmt.Println("Pertanyaan: Apakah select dibuat 3x sekaligus atau 1x lalu block?\n")

	demoSelectBlockingBehavior()
}

func demoSelectBlockingBehavior() {
	ch1 := make(chan string)
	ch2 := make(chan int)
	ch3 := make(chan bool)

	// Kirim data BERTAHAP dengan delay berbeda
	go func() {
		fmt.Println("⏰ [T+0ms]   Goroutine 1: Akan kirim data ke ch1 dalam 2 detik...")
		time.Sleep(2 * time.Second)
		fmt.Println("📤 [T+2000ms] Goroutine 1: Mengirim data ke ch1...")
		ch1 <- "Hello from ch1"
	}()

	go func() {
		fmt.Println("⏰ [T+0ms]   Goroutine 2: Akan kirim data ke ch2 dalam 4 detik...")
		time.Sleep(4 * time.Second)
		fmt.Println("📤 [T+4000ms] Goroutine 2: Mengirim data ke ch2...")
		ch2 <- 99
	}()

	go func() {
		fmt.Println("⏰ [T+0ms]   Goroutine 3: Akan kirim data ke ch3 dalam 6 detik...")
		time.Sleep(6 * time.Second)
		fmt.Println("📤 [T+6000ms] Goroutine 3: Mengirim data ke ch3...")
		ch3 <- true
	}()

	fmt.Println()
	fmt.Println("🔄 Memulai loop select...")
	fmt.Println()

	received := 0
	target := 3
	startTime := time.Now()

	for received < target {
		elapsed := time.Since(startTime).Milliseconds()

		fmt.Printf("🔍 [T+%dms] Loop iterasi ke-%d: MASUK select (mencoba ambil data...)\n",
			elapsed, received+1)

		// Select akan BLOCK di sini sampai ada channel yang ready!
		select {
		case msg := <-ch1:
			elapsed := time.Since(startTime).Milliseconds()
			fmt.Printf("📥 [T+%dms] ✅ Select SELESAI! Dapat dari ch1: %s\n\n", elapsed, msg)
			received++

		case num := <-ch2:
			elapsed := time.Since(startTime).Milliseconds()
			fmt.Printf("📥 [T+%dms] ✅ Select SELESAI! Dapat dari ch2: %d\n\n", elapsed, num)
			received++

		case flag := <-ch3:
			elapsed := time.Since(startTime).Milliseconds()
			fmt.Printf("📥 [T+%dms] ✅ Select SELESAI! Dapat dari ch3: %v\n\n", elapsed, flag)
			received++
		}

		// Setelah select selesai, loop akan lanjut ke iterasi berikutnya
	}

	elapsed := time.Since(startTime).Milliseconds()
	fmt.Printf("🎉 [T+%dms] Loop SELESAI! Semua %d data sudah diterima.\n\n", elapsed, target)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("💡 KESIMPULAN:")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("❌ SALAH: Select TIDAK dibuat 3x sekaligus dan wait parallel")
	fmt.Println()
	fmt.Println("✅ BENAR: Yang terjadi adalah:")
	fmt.Println("   1. Loop iterasi ke-1: Masuk select → BLOCK/IDLE")
	fmt.Println("   2. Tunggu sampai ada channel ready (2 detik)")
	fmt.Println("   3. Select ambil data dari channel yang ready")
	fmt.Println("   4. Loop iterasi ke-2: Masuk select BARU → BLOCK/IDLE")
	fmt.Println("   5. Tunggu sampai ada channel ready (2 detik lagi)")
	fmt.Println("   6. Select ambil data...")
	fmt.Println("   7. Dan seterusnya...")
	fmt.Println()
	fmt.Println("📌 POIN PENTING:")
	fmt.Println("   - Select BUKAN dibuat 3x sekaligus")
	fmt.Println("   - Select dibuat 1x per loop, lalu BLOCK sampai ada yang ready")
	fmt.Println("   - Setelah dapat data, loop lanjut dan buat select BARU")
	fmt.Println("   - Total waktu: 2s + 2s + 2s = ~6 detik (SEQUENTIAL, bukan parallel)")
	fmt.Println()
	fmt.Println("🔄 Analogi: Seperti antrian kasir")
	fmt.Println("   - Kasir (select) melayani customer pertama → TUNGGU sampai selesai")
	fmt.Println("   - Baru melayani customer kedua → TUNGGU lagi")
	fmt.Println("   - Bukan kasir melayani 3 customer sekaligus!")
}
