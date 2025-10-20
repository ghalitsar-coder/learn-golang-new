# Channel Blocking Demo

## Deskripsi

File ini berisi demonstrasi lengkap tentang konsep locking/blocking pada channel dalam Go. Demo ini menunjukkan berbagai skenario blocking dan cara mengelolanya dengan efektif.

## Cara Menjalankan

```bash
# Pindah ke direktori advance
cd materi/advance

# Jalankan demo
go run channel_blocking_demo.go
```

## Apa yang Akan Anda Lihat

### 1. Unbuffered Channel Blocking

- Menunjukkan bagaimana sender blocking sampai ada receiver
- Visualisasi timing dan urutan eksekusi

### 2. Buffered Channel Blocking

- Demonstrasi buffer yang penuh menyebabkan blocking
- Cara buffer bekerja sebagai queue

### 3. Receive Blocking

- Receiver blocking saat menunggu data
- Koordinasi antar goroutine

### 4. Deadlock Scenario

- Penjelasan mengapa deadlock terjadi
- Best practices untuk menghindarinya

### 5. Timeout Mechanisms

- Menggunakan `select` dengan `time.After`
- Menghindari blocking berlebihan

### 6. Non-blocking Operations

- Operasi dengan `select` dan `default`
- Pattern untuk operasi optional

### 7. Advanced Patterns

- Pipeline pattern dengan natural blocking
- Worker pool dengan bounded channels
- Timing visualization

## Konsep Penting yang Dipelajari

1. **Blocking adalah fitur, bukan bug** - Channel menggunakan blocking untuk koordinasi
2. **Unbuffered vs Buffered** - Perbedaan perilaku blocking
3. **Select statement** - Tool untuk mengelola multiple channels
4. **Timeout patterns** - Menghindari hanging operations
5. **Graceful shutdown** - Pattern untuk menghentikan goroutine dengan aman

## Output yang Diharapkan

Setiap demo section akan menampilkan:

- Timestamp atau urutan eksekusi
- Status blocking ("BLOCKING - menunggu...")
- Hasil operasi channel
- Penjelasan apa yang terjadi

## Tips untuk Pemahaman

1. **Perhatikan urutan output** - Ini menunjukkan kapan blocking terjadi
2. **Lihat timing** - Beberapa demo menggunakan sleep untuk visualisasi
3. **Bandingkan scenario** - Perbedaan unbuffered vs buffered
4. **Fokus pada koordinasi** - Bagaimana goroutine berkomunikasi

## Studi Lebih Lanjut

Setelah menjalankan demo, baca file `konsep_channel.md` untuk:

- Penjelasan teoritis yang mendalam
- Best practices dan anti-patterns
- Advanced patterns dan use cases
- Troubleshooting guide

## Modifikasi dan Eksperimen

Coba modifikasi:

- Buffer size pada buffered channels
- Timing pada sleep operations
- Jumlah goroutine dalam worker pool
- Timeout duration

Ini akan membantu memahami bagaimana parameter berbeda mempengaruhi blocking behavior.
