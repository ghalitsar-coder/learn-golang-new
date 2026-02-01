# Comparative Analysis: Worker Pool vs Pipeline
## Trade-offs & Disadvantages

Setiap pola concurrency punya "harga" yang harus dibayar. Tidak ada satu pola yang memecahkan semua masalah (No Silver Bullet).

---

## 1. Worker Pool 🏗️

**Konsep:**
Antrian satu pintu, dikeroyok oleh banyak worker identik.

### ✅ Kelebihan (Why we love it)
- **Simple:** Codingannya mudah dipahami.
- **Load Balancing:** Semua CPU core terpakai rata.
- **Resource Control:** Bisa membatasi jumlah Goroutine agar DB tidak meledak.

### ❌ Kekurangan / Trade-offs (The Ugly Truth)

#### 1. Hilangnya Urutan (Loss of Ordering)
Ini masalah terbesar.
- Input: `[Job 1, Job 2, Job 3]`
- Output bisa jadi: `[Job 2 Selesai, Job 3 Selesai, Job 1 Selesai]` (Random).
- **Why?** Worker B mungkin lebih cepat kerjanya daripada Worker A.
- **Impact:** Tidak cocok untuk data yang sekuensial (misal: Saldo Bank).

#### 2. Kompleksitas Error Handling
Jika Worker ke-5 error, siapa yang handle?
- Apakah dia harus retry sendiri?
- Apakah dia lapor ke Main?
- Bagaimana cara membatalkan 99 worker lain jika 1 worker error fatal? (Butuh `context.Context` & `errGroup`).

#### 3. Idle Resources (Over-Provisioning)
Jika kita spawn 100 Worker tapi job cuma ada 5 per menit, maka 99 Worker cuma makan memori (walau kecil) sambil "bengong".

---

## 2. Pipeline 🏭

**Konsep:**
Satu flow data dibagi menjadi beberapa tahap (Stage) yang berjalan sambung-menyambung.
Input -> [Stage 1] -> [Stage 2] -> [Stage 3] -> Output.

### ✅ Kelebihan (Why we love it)
- **Separation of Concern:** Codingan jadi rapi. Fungsi `Cuci()` terpisah dari fungsi `Bilas()`.
- **Stream Processing:** Data diproses real-time tanpa menunggu semua batch selesai.
- **Modular:** Gampang menambah stage baru di tengah jalan.

### ❌ Kekurangan / Trade-offs (The Ugly Truth)

#### 1. "The Weakest Link" (Lelet satu, lelet semua)
Kecepatan pipeline ditentukan oleh **Stage Paling Lambat**.
- Stage 1 (Cepat): 10ms.
- Stage 2 (Lambat): 1000ms.
- Stage 3 (Cepat): 10ms.
- **Akibat:** Stage 1 dan 3 akan terhambat karena Stage 2 macet (Bottleneck). Anda harus tuning manual stage mana yang butuh fan-out.

#### 2. Latency Overhead
Untuk satu biji data instan, Pipeline lebih lama daripada langsung dikerjakan.
- Analoginya: Kalau masak mie instan cuma 1 bungkus, lebih cepat dimasak sendiri daripada dioper ke 5 orang koki (oper-operannya butuh waktu).
- Pipeline bagus untuk **Throughput** (Massal), tapi jelek untuk **Latency** (Satuan).

#### 3. Complexity & Deadlock Risk
Lebih susah di-debug.
- Jika channel macet di Stage 4, Stage 1 bisa ikut macet (Backpressure).
- Risiko "Circular Dependency" (Stage A nunggu B, B nunggu A) yang bikin Deadlock.

---

## 🚦 Rangkuman Matrix Keputusan

| Fitur | Worker Pool | Pipeline |
| :--- | :--- | :--- |
| **Ordering** | ❌ Hancur (Acak) | ✅ Terjamin (FIFO) (jika stage tunggal) |
| **Complexity** | 🟢 Low (Mudah) | 🔴 High (Jalur pipa rumit) |
| **Bottleneck** | Mudah dideteksi | Sulit dideteksi (per stage) |
| **Best For** | Job Independen (Resize Image, Email Blast) | Data Flow (ETL, Video Stream, Parsing) |
