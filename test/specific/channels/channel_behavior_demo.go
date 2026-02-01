package main

import (
	"fmt"
	"time"
)

// =============================================================================
// DEMO: PERILAKU CHANNEL (BUFFERED VS UNBUFFERED)
// =============================================================================
// Channel adalah pipa komunikasi antar goroutine.
// Ada dua jenis:
// 1. Unbuffered: Pipa tanpa penampungan (Synchronous/Blocking)
// 2. Buffered: Pipa dengan tangki penampungan (Asynchronous/Non-blocking until full)
// =============================================================================

func main() {
	// Pilih skenario yang ingin dijalankan dengan uncomment:

	// --- BAGIAN 1: UNBUFFERED ---
	demoUnbufferedBlocking()
	// demoUnbufferedHandshake()

	// --- BAGIAN 2: BUFFERED ---
	// demoBufferedCapacity()
	// demoBufferedBlockingFull()

	// --- BAGIAN 3: EDGE CASES (Review ini penting!) ---
	// demoClosedChannel()
	// demoNilChannel() // ⚠️ Hati-hati, ini akan Deadlock permanen!
}

// -----------------------------------------------------------------------------
// 1. UNBUFFERED CHANNEL (Capacity = 0)
// -----------------------------------------------------------------------------
// Sifat: "Touch and Go". Pengirim TIDAK BISA menaruh data jika tidak ada yang mengambil.
// Pengirim dan Penerima harus "ketemu" di waktu yang sama (Sinkronisasi).

func demoUnbufferedBlocking() {
	fmt.Println("\n🔴 DEMO: Unbuffered Channel Blocking")
	fmt.Println("====================================")

	ch := make(chan string) // Tidak ada buffer size

	fmt.Println("1. Main: Mencoba mengirim data ke channel...")

	// ⚠️ BLOCKING POINT:
	// Kode di bawah ini akan DEADLOCK (Panic) jika dijalankan di thread utama tanpa goroutine lain.
	// Kenapa? Karena 'ch <- data' akan menunggu selamanya sampai ada yang 'receive'.
	// Jika tidak ada goroutine lain yang siap menerima, program macet.

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("   🤖 Worker: Saya siap menerima data!")
		msg := <-ch // Receive (membebaskan sender)
		fmt.Printf("   🤖 Worker: Diterima '%s'\n", msg)
	}()

	ch <- "Data Penting" // Block di sini sampai Worker bangun (2 detik)
	fmt.Println("2. Main: Data berhasil terkirim! (Sender lepas dari blocking)")
}

func demoUnbufferedHandshake() {
	fmt.Println("\n🤝 DEMO: Unbuffered Handshake (Guarantee Delivery)")
	fmt.Println("================================================")

	ch := make(chan int)

	go func() {
		fmt.Println("   📦 Sender: Mengirim paket...")
		ch <- 100 // Block sampai diterima
		fmt.Println("   📦 Sender: Paket SUDAH diterima. Tugas selesai.")
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("📥 Receiver: Menerima paket...")
	val := <-ch
	fmt.Printf("📥 Receiver: Paket %d diterima.\n", val)

	// Poin Penting:
	// Di Unbuffered, jika Sender print "Tugas selesai", kita PASTI tau Receiver sudah terima.
	// Itu jaminan pengiriman (Delivery Guarantee).
}

// -----------------------------------------------------------------------------
// 2. BUFFERED CHANNEL (Capacity > 0)
// -----------------------------------------------------------------------------
// Sifat: "Mailbox". Pengirim bisa menaruh surat di kotak surat lalu pergi.
// Tidak perlu ketemu langsung dengan penerima, ASALKAN kotak surat belum penuh.

func demoBufferedCapacity() {
	fmt.Println("\n🟢 DEMO: Buffered Channel (Non-Blocking Send)")
	fmt.Println("===========================================")

	// Kapasitas 3 slot
	ch := make(chan string, 3)

	fmt.Println("1. Mengirim 3 data (Kotak belum penuh)...")
	ch <- "A" // Masuk slot 1 (Tidak blocking)
	ch <- "B" // Masuk slot 2 (Tidak blocking)
	ch <- "C" // Masuk slot 3 (Tidak blocking)

	fmt.Println("2. Selesai mengirim 3 data. Tidak perlu receiver standby!")
	fmt.Printf("   Jumlah data di buffer: %d/3\n", len(ch))

	fmt.Println("3. Sekarang kita ambil datanya...")
	fmt.Println("   <-", <-ch)
	fmt.Println("   <-", <-ch)
	fmt.Println("   <-", <-ch)
}

func demoBufferedBlockingFull() {
	fmt.Println("\n🟠 DEMO: Buffered Channel (Blocking When Full)")
	fmt.Println("============================================")

	ch := make(chan int, 2) // Kapasitas cuma 2

	fmt.Println("Isi slot 1...")
	ch <- 1
	fmt.Println("Isi slot 2...")
	ch <- 2

	fmt.Println("Mencoba isi slot 3 (Buffer Penuh!)...")

	// Goroutine bantu ambil supaya main tidak deadlock
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("   🚑 Helper: Mengosongkan satu slot...")
		<-ch
	}()

	ch <- 3 // ⚠️ BLOCKING di sini sampai ada slot kosong
	fmt.Println("Berhasil isi slot 3!")
}

// -----------------------------------------------------------------------------
// 3. SPECIAL CASES (Channel Tertutup & Nil)
// -----------------------------------------------------------------------------

func demoClosedChannel() {
	fmt.Println("\n💀 DEMO: Closed Channel Behavior")
	fmt.Println("================================")

	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	close(ch) // Tutup channel

	fmt.Println("Channel ditutup. Mari kita baca:")

	// Baca 1 (Data sisa masih ada)
	val, ok := <-ch
	fmt.Printf("1. Value: %d, Open: %t\n", val, ok)

	// Baca 2 (Data sisa masih ada)
	val, ok = <-ch
	fmt.Printf("2. Value: %d, Open: %t\n", val, ok)

	// Baca 3 (Data habis & Channel Closed)
	// Return Zero Value (0) dan false
	val, ok = <-ch
	fmt.Printf("3. Value: %d, Open: %t (Zero Value!)\n", val, ok)

	// ⚠️ SEND ke Closed Channel = PANIC
	// ch <- 30 // Uncomment untuk lihat panic
}

func demoNilChannel() {
	fmt.Println("\n👻 DEMO: Nil Channel (The Silent Killer)")
	fmt.Println("========================================")

	var ch chan int // Belum di-make, nilai default adalah nil

	fmt.Println("Mencoba kirim ke nil channel...")
	// ⚠️ SEND ke Nil Channel = BLOCK SELAMANYA (Bukan panic, tapi hang)
	ch <- 1

	// ⚠️ RECEIVE dari Nil Channel juga BLOCK SELAMANYA
	// <-ch

	fmt.Println("Baris ini tidak akan pernah tercapai.")
}
