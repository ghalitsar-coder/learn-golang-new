package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("========== DEMO 0: Bukti Select Hanya Ambil 1 Channel ==========\n")
	demo0_ProofSelectOnlyTakesOne()

	fmt.Println("\n========== DEMO 1: SATU Select dengan 3+ Channels ==========\n")
	demo1_OneSelectMultipleChannels()

	fmt.Println("\n========== DEMO 2: Loop Select untuk Menerima Semua Data ==========\n")
	demo2_LoopSelectToReceiveAll()

	fmt.Println("\n========== DEMO 3: Select dengan Timeout ==========\n")
	demo3_SelectWithTimeout()

	fmt.Println("\n========== DEMO 4: Select dengan Default (Non-blocking) ==========\n")
	demo4_SelectWithDefault()
}

// ===== DEMO 0: BUKTI bahwa Select hanya ambil 1 channel! =====
func demo0_ProofSelectOnlyTakesOne() {
	ch1 := make(chan string, 1)
	ch2 := make(chan int, 1)
	ch3 := make(chan bool, 1)

	// Kirim data ke SEMUA channel (langsung, tanpa goroutine)
	ch1 <- "Data dari ch1"
	ch2 <- 99
	ch3 <- true

	fmt.Println("✅ SEMUA channel sudah diisi data!")
	fmt.Println("   ch1 = \"Data dari ch1\"")
	fmt.Println("   ch2 = 99")
	fmt.Println("   ch3 = true")
	fmt.Println()

	// Coba ambil dengan SATU select
	fmt.Println("🔍 Menggunakan SATU select untuk mengambil data...")
	select {
	case msg := <-ch1:
		fmt.Printf("   📥 Dapat dari ch1: %s\n", msg)
	case num := <-ch2:
		fmt.Printf("   📥 Dapat dari ch2: %d\n", num)
	case flag := <-ch3:
		fmt.Printf("   📥 Dapat dari ch3: %v\n", flag)
	}

	fmt.Println()
	fmt.Println("❓ Apakah data dari channel lain ikut terambil?")
	fmt.Println("   Mari kita cek dengan select lagi...")
	fmt.Println()

	// Coba ambil lagi
	select {
	case msg := <-ch1:
		fmt.Printf("   📥 Select ke-2: Dapat dari ch1: %s\n", msg)
	case num := <-ch2:
		fmt.Printf("   📥 Select ke-2: Dapat dari ch2: %d\n", num)
	case flag := <-ch3:
		fmt.Printf("   📥 Select ke-2: Dapat dari ch3: %v\n", flag)
	default:
		fmt.Println("   🚫 Tidak ada lagi!")
	}

	// Coba ambil lagi
	select {
	case msg := <-ch1:
		fmt.Printf("   📥 Select ke-3: Dapat dari ch1: %s\n", msg)
	case num := <-ch2:
		fmt.Printf("   📥 Select ke-3: Dapat dari ch2: %d\n", num)
	case flag := <-ch3:
		fmt.Printf("   📥 Select ke-3: Dapat dari ch3: %v\n", flag)
	default:
		fmt.Println("   🚫 Tidak ada lagi!")
	}

	fmt.Println()
	fmt.Println("💡 KESIMPULAN:")
	fmt.Println("   - Select ke-1: Ambil dari 1 channel (random jika semua ready)")
	fmt.Println("   - Select ke-2: Ambil dari 1 channel lainnya")
	fmt.Println("   - Select ke-3: Ambil dari 1 channel terakhir")
	fmt.Println("   - Total: PERLU 3 KALI SELECT untuk ambil 3 channel!")
	fmt.Println("   - Atau: Gunakan LOOP SELECT seperti di demo 2!")
}

// ===== DEMO 1: SATU Select bisa handle BANYAK channel sekaligus! =====
func demo1_OneSelectMultipleChannels() {
	ch1 := make(chan string)
	ch2 := make(chan int)
	ch3 := make(chan bool)

	// Kirim data dari 3 goroutine berbeda
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "Data dari channel 1"
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- 42
	}()

	go func() {
		time.Sleep(150 * time.Millisecond)
		ch3 <- true
	}()

	fmt.Println("Menunggu data dari salah satu channel...")

	// SATU select bisa monitor 3+ channels!
	// Select akan eksekusi case PERTAMA yang ready (paling cepat)
	select {
	case msg := <-ch1:
		fmt.Printf("✅ Dapat dari ch1: %s\n", msg)
	case num := <-ch2:
		fmt.Printf("✅ Dapat dari ch2: %d\n", num)
	case flag := <-ch3:
		fmt.Printf("✅ Dapat dari ch3: %v\n", flag)
	}

	fmt.Println("⚠️  Catatan: Select hanya mengambil SATU case yang ready pertama!")
	fmt.Println("    Data dari channel lain masih ada di buffer (belum diambil)")
}

// ===== DEMO 2: Untuk menerima SEMUA channel, gunakan LOOP! =====
func demo2_LoopSelectToReceiveAll() {
	ch1 := make(chan string)
	ch2 := make(chan int)
	ch3 := make(chan bool)
	done := make(chan bool)

	// Kirim data dari 3 goroutine
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "Hello from channel 1"
		close(ch1)
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- 99
		close(ch2)
	}()

	go func() {
		time.Sleep(150 * time.Millisecond)
		ch3 <- false
		close(ch3)
	}()

	// Goroutine untuk menerima SEMUA data
	go func() {
		received := 0
		target := 3

		for received < target {
			select {
			case msg, ok := <-ch1:
				if ok {
					fmt.Printf("📥 [%d] Dapat dari ch1: %s\n", received+1, msg)
					received++
				}
			case num, ok := <-ch2:
				if ok {
					fmt.Printf("📥 [%d] Dapat dari ch2: %d\n", received+1, num)
					received++
				}
			case flag, ok := <-ch3:
				if ok {
					fmt.Printf("📥 [%d] Dapat dari ch3: %v\n", received+1, flag)
					received++
				}
			}
		}

		fmt.Println("✅ Semua data dari 3 channel sudah diterima!")
		done <- true
	}()

	<-done
}

// ===== DEMO 3: Select dengan Timeout =====
func demo3_SelectWithTimeout() {
	ch1 := make(chan string)
	ch2 := make(chan int)
	ch3 := make(chan bool)

	// Goroutine lambat (akan timeout)
	go func() {
		time.Sleep(2 * time.Second) // Terlalu lama!
		ch1 <- "Data terlambat"
	}()

	fmt.Println("Menunggu data dengan timeout 500ms...")

	select {
	case msg := <-ch1:
		fmt.Printf("✅ Dapat dari ch1: %s\n", msg)
	case num := <-ch2:
		fmt.Printf("✅ Dapat dari ch2: %d\n", num)
	case flag := <-ch3:
		fmt.Printf("✅ Dapat dari ch3: %v\n", flag)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("⏱️  TIMEOUT! Tidak ada channel yang ready dalam 500ms")
	}
}

// ===== DEMO 4: Select dengan Default (Non-blocking) =====
func demo4_SelectWithDefault() {
	ch1 := make(chan string)
	ch2 := make(chan int)
	ch3 := make(chan bool)

	fmt.Println("Coba ambil data dari channel (non-blocking)...")

	// Tidak ada goroutine yang kirim data, langsung default!
	select {
	case msg := <-ch1:
		fmt.Printf("Dapat dari ch1: %s\n", msg)
	case num := <-ch2:
		fmt.Printf("Dapat dari ch2: %d\n", num)
	case flag := <-ch3:
		fmt.Printf("Dapat dari ch3: %v\n", flag)
	default:
		fmt.Println("🚫 Tidak ada channel yang ready, eksekusi default!")
	}

	// Sekarang kirim data ke ch2
	go func() {
		ch2 <- 123
	}()

	time.Sleep(10 * time.Millisecond) // Beri waktu goroutine kirim data

	fmt.Println("\nCoba lagi setelah ada data...")
	select {
	case msg := <-ch1:
		fmt.Printf("Dapat dari ch1: %s\n", msg)
	case num := <-ch2:
		fmt.Printf("✅ Dapat dari ch2: %d\n", num)
	case flag := <-ch3:
		fmt.Printf("Dapat dari ch3: %v\n", flag)
	default:
		fmt.Println("Tidak ada channel yang ready")
	}
}
