# Real World Concurrency: Kapan Pakai Apa?

Berikut adalah breakdown penggunaan pattern concurrency di dunia nyata, serta jawaban mendalam tentang pertanyaan FIFO Anda.

---

## 🔥 Pertanyaan Anda: Apakah Pipeline Pasti FIFO?

**Jawabannya: YA dan TIDAK.**

### 1. YA, Pasti FIFO (Pure Pipeline)
Jika setiap stage hanya memiliki **1 Goroutine** (seperti contoh `assembly_line` kita), maka urutannya **PASTI FIFO**.
Go Channels bersifat antrian FIFO.
- Mobil 1 masuk Stage "Cat".
- Mobil 2 harus menunggu Mobil 1 diproses (atau masuk antrian buffer).
- Tidak mungkin Mobil 2 keluar dari Stage "Cat" sebelum Mobil 1.

### 2. TIDAK, Bisa Acak (Pipeline + Fan-Out)
Jika di dalam satu stage Anda menyebar pekerjaan ke banyak worker (Fan-Out) untuk mempercepat, maka **urutan bisa kacau**.
- Mobil 1 dikerjakan Worker A (Lambat).
- Mobil 2 dikerjakan Worker B (Cepat).
- Mobil 2 selesai duluan!

> **Solusi jika butuh Cepat & Urut:** Gunakan pola **re-sequencing** (memberi nomor urut pada data, lalu diurutkan kembali di akhir).

---

## 🌍 Real World Scenarios: 3 Case Utama

Menjawab kebingungan Anda tentang *"Kapan saya pakai Ribuan Goroutine dan kapan saya harus membatasinya?"*, berikut adalah panduan praktis berdasarkan sifat pekerjaan (Workload).

### 1. Case: IO Bound (Banyak Menunggu)
**Strategi:** 🔥 **Spawn Ribuan Goroutine (Unlimited)**
**Kenapa?** Karena Goroutine yang "menunggu" (Network Call/Sleep) tidak memakan CPU. Mereka diparkir oleh OS dan hanya makan memori kecil (2KB).

#### 🏢 Skenario: "Tokopedia" Flash Sale Notification
Bayangkan Anda harus mengirim **1 Juta Push Notification** ke HP user detik ini juga ("Promo 90% Dimulai!").

- **Proses:** 
  1. Buka koneksi HTTP ke server Apple/Google.
  2. Kirim Data.
  3. **MENUNGGU** konfirmasi (Latency: 200ms - 500ms).
- **Analisa:** 
  - Selama 500ms menunggu, CPU komputer Anda *nganggur*. Sayang sekali jika Anda cuma menjalankan 4 worker (karena takut CPU habis).
  - CPU tidak bekerja saat menunggu response network.
- **Implementasi:**
  - **Spawn 10.000+ Goroutine sekaligus!**
  - Biarkan mereka semua menunggu respons Google secara paralel.
  - Server Anda tidak akan meledak.

---

### 2. Case: CPU Bound (Banyak Mikir)
**Strategi:** 🛡️ **Limited Worker Pool (Sejumlah Core CPU)**
**Kenapa?** Jika semua goroutine melakukan matematika berat serentak melebihi jumlah otak (Core), CPU akan macet karena overhead _Context Switching_.

#### 🏢 Skenario: "YouTube" Video Transcoding
User upload video 1000 raw 4K. Anda harus mengkompresnya menjadi MP4 1080p.

- **Proses:**
  - FFmpeg melakukan jutaan perhitungan matriks per detik (Enkripsi/Kompresi).
  - CPU Usage: 100% Mentok.
- **Analisa:**
  - **Salah:** Spawn 1000 Goroutine. Akibat: 1000 proses berebut 8 Core. CPU sibuk ganti-ganti giliran daripada kerja. Laptop panas & lemot.
  - **Benar:** Pakai **Worker Pool = 8 Worker** (jika 8 Core).
- **Implementasi:**
  - Antrian video lain sabar menunggu di Channel.
  - Begitu 1 worker selesai, langsung sikat antrian berikutnya. 

---

### 3. Case: Resource Constraint (Tetangga Lemah)
**Strategi:** 🚧 **Limited Worker Pool (Sejumlah Limit External)**
**Kenapa?** CPU kita kuat, tapi **Sistem Tujuan** (Database/API Legacy) lemah.

#### 🏢 Skenario: Migrasi Data ke "PostgreSQL Kantor"
Insert 1 Juta Baris dari CSV ke Database.

- **Analisa:**
  - Laptop Anda kuat spawn 10.000 goroutine.
  - TAPI, Database diset `max_connections = 100`.
  - Jika ditembak 10.000 concurrent request, Database akan **Reject/Crash**.
- **Implementasi:**
  - Gunakan **Worker Pool = 80 Worker**.
  - Ini teknik **Rate Limiting** (Ngerem sengaja).
  - Tujuannya melindungi Database agar tidak mati.

---

## 🚀 Rangkuman Matrix Keputusan

| Apa yang dilakukan program? | Contoh | Strategi | Jumlah Goroutine? |
| :--- | :--- | :--- | :--- |
| **Menunggu** (IO Bound) | Request API, DB Query, Chat App | **Spawn per Request** | Ribuan / Sebanyak User |
| **Mikir Berat** (CPU Bound) | Video Encode, Image Resize, Hashing | **Worker Pool** | `runtime.NumCPU()` |
| **Akses Resource Rentan** | Insert DB, Access Legacy API | **Worker Pool** | `Limit Resource Tujuan` |
