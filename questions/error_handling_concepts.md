# Konsep Dasar Error Handling dalam Go

## Apa itu Error dalam Go?

Error dalam Go adalah nilai yang dikembalikan oleh fungsi ketika terjadi kesalahan. Error dalam Go mengikuti prinsip "errors are values", yang berarti error diperlakukan seperti nilai biasa dan harus diperiksa secara eksplisit. Ini berbeda dengan bahasa pemrograman lain yang menggunakan exception untuk menangani error.

## Interface Error

Dalam Go, error direpresentasikan oleh interface `error` yang didefinisikan dalam package `builtin`:

```go
type error interface {
    Error() string
}
```

Interface ini hanya memiliki satu method `Error() string` yang mengembalikan pesan error dalam bentuk string.

## Membuat Error dengan fmt.Errorf

`fmt.Errorf` adalah fungsi yang paling umum digunakan untuk membuat error baru dengan pesan yang diformat. Ini adalah cara yang idiomatic dalam Go untuk membuat error.

```go
if len(nums) == 0 {
    return 0, fmt.Errorf("array kosong")
}
```

Keuntungan menggunakan `fmt.Errorf`:
1. Memungkinkan pembuatan pesan error yang dapat diformat dan disesuaikan
2. Mudah dibaca dan dipahami
3. Sesuai dengan konvensi Go

## Wrapping Error dengan %w

Go 1.13 memperkenalkan fitur error wrapping yang memungkinkan kita untuk membungkus error asli dengan error baru:

```go
if err != nil {
    return fmt.Errorf("gagal membuka file: %w", err)
}
```

Dengan menggunakan `%w`, error asli tetap dipertahankan dan dapat diperiksa dengan `errors.Is()` atau `errors.As()`.

## Perbedaan Log dan Error Handling

### Package `log`
Package `log` digunakan untuk mencatat pesan ke output standar atau file log. Beberapa fungsi penting dalam package ini:

- `log.Println()`: Mencatat pesan dan melanjutkan eksekusi program
- `log.Printf()`: Mencatat pesan dengan format tertentu
- `log.Fatal()`: Mencatat pesan dan menghentikan program dengan kode keluar 1
- `log.Fatalf()`: Mencatat pesan dengan format tertentu dan menghentikan program

### Error Handling
Error handling adalah proses memeriksa dan menangani error yang dikembalikan oleh fungsi. Dalam Go, ini dilakukan dengan memeriksa nilai error yang dikembalikan:

```go
result, err := someFunction()
if err != nil {
    // Tangani error
}
```

## Best Practices Error Handling

1. **Selalu periksa error**: Jangan pernah mengabaikan error dengan `_`
2. **Gunakan early return**: Tangani error segera setelah diperiksa
3. **Gunakan fmt.Errorf dengan %w**: Untuk membungkus error asli
4. **Buat error kustom jika perlu**: Untuk error yang spesifik dalam aplikasi
5. **Dokumentasikan error**: Jelaskan error apa yang bisa dikembalikan fungsi
6. **Gunakan context untuk timeout**: Terutama dalam operasi I/O

## Perbandingan Error Handling vs Exception

Go menggunakan pendekatan "error as value" yang berbeda dengan bahasa lain yang menggunakan exception:

| Aspek | Go Error Handling | Exception-based |
|-------|-------------------|-----------------|
| Mekanisme | Return value | Throw/Catch |
| Pemeriksaan | Eksplisit | Implisit |
| Performance | Lebih cepat | Lebih lambat |
| Control | Lebih besar | Terbatas |

Pendekatan Go memberikan kontrol penuh kepada developer untuk menangani error sesuai kebutuhan, tetapi memerlukan disiplin untuk selalu memeriksa error.