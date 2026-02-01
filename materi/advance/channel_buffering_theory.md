# Deep Dive: Go Channel Buffering
## Unbuffered vs Buffered & Perilaku Ekstrim

Channel adalah mekanisme komunikasi core di Go. Memahami perilakunya di berbagai kondisi sangat krusial untuk menghindari **Deadlock** dan **Goroutine Leak**.

---

## 1. Unbuffered Channel (Kapasitas 0)
```go
ch := make(chan int)
```
Sifat: **Sync Broadcast / Handshake**

Pikirkan seperti **serah terima estafet**.
- Pelari A (Sender) tidak bisa melepaskan tongkat jika Pelari B (Receiver) belum memegangnya.
- Pelari B (Receiver) tidak bisa lari jika Pelari A belum memberikan tongkat.

**Aturan Main:**
- **Send (`ch <- val`)**: Akan **BLOCK** (berhenti) sampai ada goroutine lain yang melakukan Receive.
- **Receive (`<-ch`)**: Akan **BLOCK** sampai ada goroutine lain yang melakukan Send.
- **Jaminan**: Jika baris `ch <- val` selesai dieksekusi, Anda **100% yakin** data sudah diterima oleh pihak lain.

**Kapan dipakai?**
- Sinkronisasi ketat.
- Menjamin data processed before moving on.
- Signal/Trigger antar goroutine.

---

## 2. Buffered Channel (Kapasitas > 0)
```go
ch := make(chan int, 3) // Kapasitas 3
```
Sifat: **Queue / Mailbox**

Pikirkan seperti **kotak surat**.
- Pengirim bisa memasukkan surat lalu pergi kerja (selama kotak belum penuh).
- Pengirim tidak perlu tau apakah surat sudah dibaca atau belum saat itu juga.

**Aturan Main:**
- **Send (`ch <- val`)**: 
  - **TIDAK BLOCK** selama buffer belum penuh (`len < cap`).
  - **BLOCK** jika buffer penuh, sampai ada slot kosong (ada yang receive).
- **Receive (`<-ch`)**:
  - **TIDAK BLOCK** selama buffer ada isinya.
  - **BLOCK** jika buffer kosong, sampai ada data masuk.

**Kapan dipakai?**
- Decoupling (Producer lebih cepat dari Consumer atau sebaliknya).
- Mencegah blocking sesaat (Burst traffic).
- Mengumpulkan hasil (Fan-In) tanpa memblokir worker.

---

## 3. Matriks Perilaku (Sangat Penting!)

Tabel ini adalah "kitab suci" perilaku channel. Hafalkan efeknya.

| Operasi | Nil Channel (`var c chan int`) | Closed Channel | Open & Full | Open & Empty | Open & Normal |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Send** `c <-` | **Block Selamanya** 💀 | **PANIC** 💥 | **Block** 🛑 | Berhasil ✅ | Berhasil ✅ |
| **Receive** `<-c` | **Block Selamanya** 💀 | **Zero Value** (non-block) 💨 | Berhasil ✅ | **Block** 🛑 | Berhasil ✅ |
| **Close** `close(c)` | **PANIC** 💥 | **PANIC** 💥 | Berhasil | Berhasil | Berhasil |
| **Len** `len(c)` | 0 | Jumlah sisa item | Kapasitas | 0 | Jumlah item |
| **Cap** `cap(c)` | 0 | Kapasitas | Kapasitas | Kapasitas | Kapasitas |

### Penjelasan Detail Kondisi Kritis:

1.  **Block Selamanya (Deadlock Risk)**
    Terjadi jika Anda mengirim/menerima dari `nil` channel. Ini sering terjadi jika Anda lupa inisialisasi `make(chan ...)`.
    ```go
    var ch chan int // nil
    <-ch // Mampus, goroutine tidur selamanya.
    ```

2.  **Panic pada Send to Closed**
    Jangan pernah mengirim data ke channel yang sudah ditutup. `close()` adalah sinyal dari SENDER bahwa "tidak ada data lagi". Receiver tidak boleh mengirim balik (unidirectional logic).
    
3.  **Zero Value pada Receive from Closed**
    Ini fitur, bukan bug. Jika channel ditutup, receiver tidak akan crash, tapi akan menerima nilai default tipe data itu (0 untuk int, "" untuk string, nil untuk ptr) secara terus menerus tanpa henti (infinite loop jika tidak di-cek).
    
    **Cara Cek yang Benar (Comma OK Idiom):**
    ```go
    val, ok := <-ch
    if !ok {
        fmt.Println("Channel sudah tutup!")
        return
    }
    ```
    Atau pakailah `range`:
    ```go
    for val := range ch {
        // Otomatis berhenti loop jika channel closed & kosong
    }
    ```

---

## 🏗️ Studi Kasus: Buffer Size = Jumlah Job?

**Pertanyaan:** *"Jika saya punya 1000 Job, apakah sebaiknya saya buat `make(chan Job, 1000)`?"*

**Jawabannya: TERGANTUNG MEMORI.**

### ✅ Pros: Fire and Forget
Jika buffer size = jumlah job, maka **Producer tidak akan pernah blocking**.
```go
jobs := make(chan int, 1000)
for i:=0; i<1000; i++ {
    jobs <- i // WUSSH! Masuk semua instan.
}
close(jobs)
// Producer bisa langsung exit atau ngerjain hal lain.
```
Ini sangat nyaman untuk **Batch Processing kecil** (misal < 10.000 item).

### ❌ Cons: Boros Memori
Setiap slot buffer memakan memori, walaupun kosong.
Jika Anda punya **1 Juta Job** dan struct Job-nya besar (misal ada String panjang / Image bytes):
```go
jobs := make(chan BigStruct, 1000000) // ⚠️ RAM MELEDAK
```
Program Anda akan langsung mengalokasikan memori raksasa di awal, padahal Worker cuma bisa memproses 5 item sekaligus. **Mubazir**.

### 💡 Best Practice
Gunakan "Small Buffer" (misal 100) sebagai penampung sementara (Shock Absorber). 
Worker biasanya memproses lebih lambat daripada Producer, jadi buffer 1 Juta pun lama-lama akan penuh dan Producer tetap akan blocking juga.

**Rule of Thumb:**
- Job Sedikit (< 1000): Boleh Buffer = Job Size (Enak, codingan simple).
- Job Raksasa (> 10rb): Buffer Kecil saja (50-100). Biarkan Producer blocking sedikit demi menjaga RAM stabil.

---

## Kesimpulan Skenario

1. **Ingin Kirim data dan PASTIKAN data itu diproses dulu sebelum lanjut?**
   👉 Gunakan **Unbuffered Channel**.

2. **Ingin Kirim data "Fire and Forget" (biar cepat) dan rela antri sedikit jika penuh?**
   👉 Gunakan **Buffered Channel**.

3. **Ingin membatasi concurrency (Semaphore)?**
   👉 Gunakan **Buffered Channel** dengan kapasitas = limit.
