# Pembaruan Go 1.25: Ringkasan Fitur Baru dan Perbandingan

## Ikhtisar

Go 1.25 membawa berbagai peningkatan, perbaikan, dan fitur baru. Salah satu perubahan paling signifikan bagi developer yang bekerja dengan konkurensi adalah penambahan metode `WaitGroup.Go` yang menyederhanakan penggunaan `sync.WaitGroup` dengan goroutine.

---

## `sync.WaitGroup.Go` - Perubahan Konkurensi Paling Penting

### Sebelum Go 1.25

Pattern umum untuk menunggu selesainya goroutine adalah sebagai berikut:

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup

    tasks := []string{"Tugas 1", "Tugas 2", "Tugas 3"}

    for _, task := range tasks {
        wg.Add(1) // Menambahkan counter sebelum membuat goroutine
        go func(t string) {
            defer wg.Done() // Mengurangi counter saat goroutine selesai
            // Lakukan pekerjaan
            fmt.Printf("Menjalankan %s\n", t)
        }(task) // Menyisipkan variabel ke closure
    }

    wg.Wait() // Menunggu semua goroutine selesai
    fmt.Println("Semua tugas selesai!")
}
```

### Sesudah Go 1.25

Go 1.25 memperkenalkan metode `WaitGroup.Go`, yang secara otomatis menangani `Add(1)` dan `Done()`:

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup

    tasks := []string{"Tugas 1", "Tugas 2", "Tugas 3"}

    for _, task := range tasks {
        wg.Go(func() { // Langsung menjalankan goroutine dalam konteks WaitGroup
            // Lakukan pekerjaan
            fmt.Printf("Menjalankan %s\n", task)
        })
    }

    wg.Wait() // Menunggu semua goroutine selesai
    fmt.Println("Semua tugas selesai!")
}
```

### Perbedaan dan Keuntungan

- **Lebih Ringkas**: Tidak perlu lagi menulis `wg.Add(1)` dan `defer wg.Done()` secara manual.
- **Lebih Aman**: Mengurangi risiko kesalahan, seperti lupa memanggil `Add` atau `Done`, atau memanggilnya di tempat yang salah.
- **Lebih Mudah Dibaca**: Kode menjadi lebih jelas dan fokus pada logika bisnis daripada manajemen konkurensi.
- **Mendeteksi Kesalahan**: `vet` analyzer baru (`waitgroup`) dapat mendeteksi penggunaan `WaitGroup.Add` yang salah tempat, yang sebelumnya bisa menyebabkan race condition.

---

## Fitur-Fitur Baru Lainnya (Ringkasan)

### 1. Package `testing/synctest`
- Membantu dalam pengujian kode konkuren.
- `Test` function menjalankan test dalam lingkungan terisolasi (bubble) dengan kontrol waktu.
- `Wait` function menunggu semua goroutine dalam bubble saat ini untuk block.

### 2. Perubahan GOMAXPROCS
- Kini mempertimbangkan batas cgroup CPU di Linux secara default.
- Nilai GOMAXPROCS diperbarui secara dinamis jika jumlah CPU logical atau batas cgroup berubah.
- Fungsi baru `runtime.SetDefaultGOMAXPROCS()`.

### 3. Flight Recorder API (`runtime/trace`)
- Menyediakan cara ringan untuk menangkap trace eksekusi runtime.
- Berguna untuk mendiagnosis event yang jarang terjadi.

### 4. Peningkatan Tools
- `go build -asan` sekarang mendeteksi kebocoran memori secara default.
- `go doc -http` membuka dokumentasi di browser.
- Directive `go mod ignore`.

### 5. Perubahan Lainnya
- GC eksperimental baru ("greenteagc").
- Perbaikan keamanan compiler.
- Support DWARF5.
- Perubahan pada `crypto/tls`, `net/http`, `slog`, dll.

---

## Kesimpulan

Go 1.25 membawa peningkatan signifikan terutama dalam manajemen konkurensi melalui `sync.WaitGroup.Go`, yang menyederhanakan dan membuat kode konkuren lebih aman dan mudah dibaca. Selain itu, ada banyak peningkatan lainnya di bidang tooling, runtime, dan keamanan. Pengguna yang sering menggunakan `sync.WaitGroup` dan goroutine seharusnya merasakan dampak positif langsung dari perubahan ini.