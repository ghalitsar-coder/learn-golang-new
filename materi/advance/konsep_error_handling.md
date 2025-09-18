Berikut adalah artikel lengkap dalam format `file.md` yang membahas **Error Handling Terbaik di Go (Best Practices)** — cocok untuk dokumentasi, tim development, atau pembelajaran pribadi.

---

```markdown
# 🛡️ Best Practices Error Handling di Go

Go (Golang) memiliki pendekatan unik terhadap penanganan error: **error sebagai nilai**, bukan exception. Ini memberi kontrol penuh ke developer, tapi juga menuntut disiplin agar kode tetap bersih, aman, dan mudah dikelola.

Artikel ini membahas prinsip dan praktik terbaik dalam menangani error di Go — dari dasar hingga teknik lanjutan.

---

## 📌 1. Jangan Pernah Mengabaikan Error

Salah satu kesalahan paling umum adalah mengabaikan error dengan `_`.

### ❌ Salah:
```go
file, _ := os.Open("config.json") // error diabaikan!
data, _ := ioutil.ReadAll(file)
```

> Jika file tidak ada, program bisa crash atau berperilaku tak terduga.

### ✅ Benar:
```go
file, err := os.Open("config.json")
if err != nil {
    return fmt.Errorf("gagal membuka file config: %w", err)
}
defer file.Close()

data, err := ioutil.ReadAll(file)
if err != nil {
    return fmt.Errorf("gagal membaca file config: %w", err)
}
```

> **Prinsip:** Selalu periksa error. Jika Anda yakin error tidak penting, beri komentar jelas *mengapa*.

---

## 📌 2. Gunakan Early Return dan Hindari Nesting Berlebihan

Go menganut gaya “fail fast”. Gunakan `if err != nil` segera setelah pemanggilan fungsi, lalu `return` atau `continue`.

### ❌ Nested dan sulit dibaca:
```go
func process() error {
    file, err := os.Open("data.txt")
    if err == nil {
        data, err := ioutil.ReadAll(file)
        if err == nil {
            result, err := parse(data)
            if err == nil {
                return save(result)
            } else {
                return err
            }
        } else {
            return err
        }
    } else {
        return err
    }
}
```

### ✅ Flat, mudah dibaca, dan Go-idiomatic:
```go
func process() error {
    file, err := os.Open("data.txt")
    if err != nil {
        return fmt.Errorf("gagal membuka file: %w", err)
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return fmt.Errorf("gagal membaca file: %w", err)
    }

    result, err := parse(data)
    if err != nil {
        return fmt.Errorf("gagal parsing data: %w", err)
    }

    if err := save(result); err != nil {
        return fmt.Errorf("gagal menyimpan hasil: %w", err)
    }

    return nil
}
```

> **Prinsip:** “*Handle error, lalu keluar*.” Hindari `else` setelah `return`.

---

## 📌 3. Wrap Error dengan `%w` untuk Menjaga Konteks

Go 1.13+ mendukung wrapping error dengan `fmt.Errorf("%w", err)` dan pemeriksaan dengan `errors.Is()` dan `errors.As()`.

### ✅ Contoh wrapping:
```go
if err != nil {
    return fmt.Errorf("gagal menghubungi database: %w", err)
}
```

### ✅ Pemeriksaan error spesifik:
```go
err := connectDB()
if errors.Is(err, sql.ErrNoRows) {
    log.Println("Tidak ada data ditemukan")
} else if err != nil {
    log.Fatal("Error tak terduga:", err)
}
```

> **Prinsip:** Gunakan `%w` untuk menjaga *stack trace* asli dan memungkinkan pemeriksaan tipe error di atas.

---

## 📌 4. Buat Error Kustom Jika Diperlukan

Untuk error bisnis/logika aplikasi, buat tipe error kustom agar bisa dikenali dan ditangani secara spesifik.

### ✅ Contoh:
```go
var ErrUserNotFound = errors.New("user tidak ditemukan")

func GetUser(id int) (*User, error) {
    user := findUserInDB(id)
    if user == nil {
        return nil, ErrUserNotFound
    }
    return user, nil
}

// Di handler/pemanggil:
user, err := GetUser(123)
if errors.Is(err, ErrUserNotFound) {
    return http.StatusNotFound, "User tidak ada"
}
if err != nil {
    return http.StatusInternalServerError, "Server error"
}
```

> **Prinsip:** Gunakan `errors.New` untuk error statis, atau buat tipe struct jika butuh metadata tambahan.

---

## 📌 5. Logging vs Returning Error

- **Return error** → untuk memberi tahu *caller* bahwa ada masalah yang perlu ditangani.
- **Log error** → untuk mencatat kejadian ke sistem (file, stdout, monitoring).

### ✅ Contoh kombinasi:
```go
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
    err := handleRequest(r)
    if err != nil {
        log.Printf("Gagal menangani request: %v", err) // log untuk developer/ops
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}
```

> **Prinsip:** Jangan log error lalu diam — pastikan error juga dikembalikan atau ditangani sesuai konteks.

---

## 📌 6. Gunakan `defer` untuk Cleanup, Tapi Hati-hati dengan Error-nya

`defer` sering dipakai untuk `Close()`, tapi `Close()` juga bisa error!

### ✅ Contoh aman:
```go
file, err := os.Create("output.txt")
if err != nil {
    return err
}
defer func() {
    if closeErr := file.Close(); closeErr != nil {
        log.Printf("Gagal menutup file: %v", closeErr)
    }
}()
```

> **Prinsip:** Jika `Close()` bisa gagal dan penting, jangan abaikan error-nya — minimal log.

---

## 📌 7. Hindari Panic untuk Error yang Dapat Diantisipasi

`panic` hanya untuk kondisi *unrecoverable*: bug, invariant rusak, atau kondisi darurat.

### ❌ Jangan:
```go
if user == nil {
    panic("User tidak boleh nil!") // Ini bisa dicegah dengan validasi!
}
```

### ✅ Seharusnya:
```go
if user == nil {
    return nil, errors.New("user tidak boleh nil")
}
```

> **Prinsip:** Gunakan `panic` hanya untuk hal yang benar-benar tidak bisa dilanjutkan. Gunakan `error` untuk semua kondisi yang bisa diantisipasi.

---

## 📌 8. Gunakan Tools: `errcheck`, `staticcheck`, `golangci-lint`

Gunakan linter untuk mendeteksi error yang diabaikan atau pola buruk.

Contoh `.golangci.yml`:
```yaml
linters:
  enable:
    - errcheck
    - goerr113
    - errorlint
```

> **Prinsip:** Otomasi pengecekan error handling di CI/CD pipeline.

---

## 📌 9. Dokumentasikan Error yang Dikembalikan

Dokumentasi membantu pengguna fungsi tahu error apa yang mungkin terjadi.

### ✅ Contoh:
```go
// GetUser mengambil user berdasarkan ID.
// Mengembalikan ErrUserNotFound jika user tidak ada.
// Mengembalikan error database jika terjadi kegagalan query.
func GetUser(id int) (*User, error) {
    // ...
}
```

> **Prinsip:** Dokumentasi error = komunikasi antar developer.

---

## 📌 10. Gunakan Context untuk Timeout dan Cancellation

Untuk operasi I/O (HTTP, DB, gRPC), selalu gunakan `context.Context` untuk menghindari hanging.

### ✅ Contoh:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.Do(req.WithContext(ctx))
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        return fmt.Errorf("timeout: %w", err)
    }
    return fmt.Errorf("gagal request: %w", err)
}
```

> **Prinsip:** Semua operasi jaringan/database harus timeout. Jangan biarkan sistem nge-hang.

---

## 🎯 Kesimpulan: Prinsip Utama Error Handling di Go

| Prinsip | Penjelasan |
|--------|------------|
| **Errors are values** | Perlakukan error seperti nilai biasa — cek, wrap, kembalikan. |
| **Fail fast** | Tangani error sedini mungkin, jangan tunda. |
| **Be explicit** | Jangan sembunyikan error. Jangan gunakan panic untuk error biasa. |
| **Preserve context** | Gunakan `%w` untuk menjaga informasi error asli. |
| **Log wisely** | Log untuk observabilitas, return error untuk kontrol flow. |
| **Test error paths** | Unit test harus mencakup kasus error juga. |

---

## 🧪 Bonus: Contoh Lengkap Fungsi dengan Error Handling Baik

```go
import (
    "context"
    "fmt"
    "log"
    "time"
)

var ErrTimeout = fmt.Errorf("permintaan melebihi batas waktu")

func FetchUserData(ctx context.Context, userID string) ([]byte, error) {
    // Validasi input
    if userID == "" {
        return nil, fmt.Errorf("userID tidak boleh kosong")
    }

    // Buat context dengan timeout
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    // Simulasi panggilan eksternal
    data, err := externalAPI(ctx, userID)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
        }
        return nil, fmt.Errorf("gagal mengambil data user: %w", err)
    }

    if len(data) == 0 {
        return nil, fmt.Errorf("data kosong untuk user: %s", userID)
    }

    return data, nil
}

func externalAPI(ctx context.Context, userID string) ([]byte, error) {
    // Simulasi jeda dan error
    select {
    case <-time.After(2 * time.Second):
        return []byte(`{"name":"John"}`), nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}
```

---

## 📚 Referensi

- [Go Blog: Errors are values](https://blog.golang.org/errors-are-values)
- [Go 1.13 Error Handling](https://blog.golang.org/go1.13-errors)
- [Uber Go Style Guide - Error Handling](https://github.com/uber-go/guide/blob/master/style.md#error-strings)
- [Effective Go - Errors](https://golang.org/doc/effective_go.html#errors)

---

✅ Dengan menerapkan praktik di atas, kode Go Anda akan lebih **robust**, **mudah di-debug**, dan **mudah dipelihara** — bahkan saat skala aplikasi membesar.

Happy coding with proper error handling! 🚀
```

---

### 📥 Cara Gunakan:
Simpan sebagai `error-handling-best-practices-go.md` dan bagikan ke tim Anda, atau gunakan sebagai referensi pribadi.

Ingin versi PDF atau HTML? Bisa di-generate dengan tools seperti `pandoc` atau VSCode + Markdown Preview.

---

Jika Anda ingin versi yang lebih ringkas, versi slide, atau contoh unit test untuk error handling — beri tahu saya, saya bantu buatkan!