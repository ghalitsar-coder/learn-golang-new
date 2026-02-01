# Deep Dive: Fan-Out & Fan-In Pattern

Menjawab pertanyaan Anda tentang mekanisme channel di Fan-Out/Fan-In.

---

## 1. Apakah Fan-Out harus dari 1 Channel?
**Jawabannya: BIASANYA YA (90% Skenario).**

Fan-Out artinya "Menyebar Beban". Konsep dasarnya adalah **Load Balancing**. 
- Bayangkan antrian loket bank (1 Channel input).
- Ada 5 Teller (5 Goroutines/Workers).
- Semua Teller mata-matanya tertuju ke 1 antrian yang sama. Begitu ada 1 nasabah datang, Teller yang nganggur langsung ambil.

**Jadi flow-nya:** `1 Channel Source` -> `N Workers`.

---

## 2. Fan-Out vs Worker Pool: Apa Bedanya?

Ini adalah kebingungan yang sangat umum karena bentuk kodenya mirip.
Sebenarnya, **Worker Pool adalah SALAH SATU bentuk Fan-Out (Bounded Fan-Out).**

### 🧠 Konsep Hirarkinya:
1.  **Fan-Out (Payung Besar):** Konsep "Menyebar 1 pekerjaan ke Banyak Goroutine".
2.  **Unbounded Fan-Out (Spawn per Job):**
    - Setiap data masuk -> Langsung `go func()`.
    - Jika ada 1 Juta data -> 1 Juta Goroutine.
    - **Sifat:** Cepat, tapi boros memori.
3.  **Bounded Fan-Out (Worker Pool):**
    - Data masuk -> Ditunggu oleh **Sejumlah Tetap** Goroutine (misal 5 Worker).
    - Jika ada 1 Juta data -> Antri diproses oleh 5 Worker.
    - **Sifat:** Terkendali, hemat memori.

**Kesimpulan:**
Code `fan_out_fan_in.go` yang saya buat sebelumnya adalah **Fan-Out tipe Worker Pool**.
Jadi Anda benar seratus persen, itu adalah Worker Pool. Dan Worker Pool adalah implementasi paling aman dari Fan-Out.

---

## 3. Multiplexing (Banyak Channel Input)
Jika Anda punya `Channel A`, `Channel B`, `Channel C` dan ingin semua diproses, Anda punya 2 pilihan:

### Opsi A: Merge Dulu (Fan-In di awal)
Anda buat 1 goroutine khusus untuk menggabungkan Data A, B, C ke dalam 1 `Channel TOTAL`. Baru setelah itu worker mengambil dari `Channel TOTAL`.

### Opsi B: Worker yang "Rakus" (Select)
Worker-nya diprogram untuk mendengarkan banyak telinga sekaligus.
```go
for {
    select {
    case job := <-chanA:
        process(job)
    case job := <-chanB:
        process(job)
    }
}
```

---

## 4. Apakah "Fan-In" (Penggabungan) harus Manual?
**YA.** Channel di Go tidak otomatis bergabung seperti sungai. Anda harus membuat "Pipa Sambungan" (Goroutine) yang memindahkan isi dari banyak channel ke satu channel.

Pola standarnya:
1.  Bikin `mergedChan`.
2.  Bikin loops untuk baca `chanA` -> kirim ke `mergedChan`.
3.  Bikin loops untuk baca `chanB` -> kirim ke `mergedChan`.
4.  Pakai `WaitGroup` untuk tau kapan A dan B habis, baru tutup `mergedChan`.

---

## 5. Ringkasan Kapan Pakai Apa?

| Pola | Kapan Pakai? |
| :--- | :--- |
| **Unbounded Fan-Out** | Trafik sedikit, butuh latensi super rendah (misal: Notifikasi user). |
| **Worker Pool (Bounded)** | Trafik tinggi, resource terbatas (misal: Image Processing, DB Insert). |
| **Fan-In** | Menggabungkan hasil dari Fan-Out untuk dibuat report/summary. |
