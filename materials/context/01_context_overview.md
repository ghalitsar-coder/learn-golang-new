# Context di Golang: Panduan Lengkap dan Komprehensif

## Daftar Isi
1. [Pendahuluan](#pendahuluan)
2. [Apa Itu Context?](#apa-itu-context)
3. [Mengapa Context Diperlukan?](#mengapa-context-diperlukan)
4. [Anatomi Interface Context](#anatomi-interface-context)
5. [Tipe-tipe Context](#tipe-tipe-context)
6. [Cara Kerja Internal Context](#cara-kerja-internal-context)
7. [Skenario Penggunaan](#skenario-penggunaan)
8. [Best Practices](#best-practices)
9. [Anti-Patterns](#anti-patterns)
10. [Contoh Implementasi Real-World](#contoh-implementasi-real-world)

---

## Pendahuluan

**Context** adalah salah satu fitur fundamental dalam pemrograman Go yang sering kurang dipahami oleh developer pemula. Dalam bahasa sederhana, context adalah mekanisme untuk mengontrol dan mengkoordinasikan goroutine-goroutine yang berjalan secara concurrent.

Bayangkan Anda adalah seorang manajer restoran yang memiliki banyak karyawan (goroutine). Ketika ada pesanan batal, Anda perlu memberitahu semua karyawan yang sedang mengerjakan pesanan tersebut untuk berhenti. Context memberikan Anda "radio komunikasi" untuk melakukan hal ini.

### Kapan Context Diperkenalkan?

Context diperkenalkan di Go 1.7 (Agustus 2016) sebagai bagian dari standard library melalui package `context`. Sebelumnya, package ini dikembangkan di `golang.org/x/net/context`.

### Import Statement

```go
import "context"
```

---

## Apa Itu Context?

### Definisi Formal

Context adalah **carrier** (pembawa) yang membawa:
1. **Deadline** - batas waktu kapan operasi harus selesai
2. **Cancellation Signal** - sinyal pembatalan untuk menghentikan operasi
3. **Request-scoped values** - nilai-nilai yang terkait dengan request tertentu

### Analogi Sederhana

Mari kita gunakan analogi **Sistem Pemesanan Pizza**:

```
🍕 Pesanan Pizza (Request)
├── Context = "Kartu Pesanan"
│   ├── Deadline: "Harus selesai dalam 30 menit"
│   ├── Cancel: "Tombol batalkan pesanan"
│   └── Values: "Info pelanggan, alamat, dll"
│
├── Koki (Goroutine 1) - Membuat adonan
├── Koki (Goroutine 2) - Menyiapkan topping
├── Koki (Goroutine 3) - Memanggang pizza
└── Kurir (Goroutine 4) - Mengantar pizza
```

Ketika pelanggan membatalkan pesanan, context akan memberitahu SEMUA goroutine untuk berhenti bekerja secara bersamaan.

### Visualisasi Tree Context

```
                    ┌─────────────────┐
                    │ context.Background │
                    │   (Root Context)   │
                    └─────────┬─────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
    ┌─────────▼─────────┐         ┌──────────▼──────────┐
    │ WithTimeout(5s)   │         │   WithCancel()      │
    │  (HTTP Handler)   │         │  (Background Job)   │
    └─────────┬─────────┘         └──────────┬──────────┘
              │                               │
    ┌─────────▼─────────┐         ┌──────────▼──────────┐
    │ WithValue(userID) │         │  WithDeadline(...)  │
    │  (Add user info)  │         │   (Sub-task)        │
    └─────────┬─────────┘         └─────────────────────┘
              │
    ┌─────────▼─────────┐
    │  Database Query   │
    │   (Child task)    │
    └───────────────────┘
```

---

## Mengapa Context Diperlukan?

### Problem 1: Goroutine Leak (Kebocoran Goroutine)

Tanpa context, goroutine bisa "bocor" dan terus berjalan meskipun sudah tidak diperlukan.

```go
// ❌ BURUK: Tanpa context - Goroutine Leak!
func fetchDataBad() {
    go func() {
        for {
            // Goroutine ini akan berjalan selamanya!
            // Tidak ada cara untuk menghentikannya
            time.Sleep(time.Second)
            fmt.Println("Still running...")
        }
    }()
}

// ✅ BAIK: Dengan context - Dapat dikontrol
func fetchDataGood(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Gracefully stopped!")
                return
            default:
                time.Sleep(time.Second)
                fmt.Println("Running...")
            }
        }
    }()
}
```

### Problem 2: Timeout Management

Tanpa context, sulit mengontrol berapa lama suatu operasi boleh berjalan.

```go
// ❌ BURUK: Tidak ada timeout
func slowOperation() {
    // Bisa hang selamanya!
    resp, err := http.Get("https://slow-api.example.com")
    // ...
}

// ✅ BAIK: Dengan timeout
func slowOperationWithTimeout(ctx context.Context) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    req, _ := http.NewRequestWithContext(ctx, "GET", "https://slow-api.example.com", nil)
    resp, err := http.DefaultClient.Do(req)
    // Akan timeout setelah 5 detik
}
```

### Problem 3: Cascading Cancellation

Ketika satu operasi dibatalkan, semua sub-operasi harus ikut dibatalkan.

```go
// Ilustrasi: Request Processing Pipeline
//
// HTTP Request
//     ├── Parse Request
//     ├── Authenticate User
//     ├── Fetch Data from Database
//     │       ├── Query 1
//     │       ├── Query 2
//     │       └── Query 3
//     └── Return Response
//
// Jika client disconnect di tengah-tengah,
// SEMUA operasi harus dihentikan!
```

### Problem 4: Request-Scoped Data

Meneruskan data yang terkait dengan request (seperti userID, requestID) ke seluruh call chain.

```go
// Tanpa context, harus passing di setiap function
func handler(userID string) {
    process(userID)
}
func process(userID string) {
    save(userID)
}
func save(userID string) {
    log(userID)
}

// Dengan context, lebih bersih
func handlerWithCtx(ctx context.Context) {
    process(ctx)
}
func processWithCtx(ctx context.Context) {
    save(ctx)
}
func saveWithCtx(ctx context.Context) {
    userID := ctx.Value("userID").(string)
    log(userID)
}
```

---

## Anatomi Interface Context

### Interface Definition

```go
type Context interface {
    // Deadline mengembalikan waktu kapan context ini harus selesai
    // ok == false jika tidak ada deadline
    Deadline() (deadline time.Time, ok bool)
    
    // Done mengembalikan channel yang akan di-close ketika context dibatalkan
    Done() <-chan struct{}
    
    // Err mengembalikan error yang menjelaskan mengapa context dibatalkan
    // nil jika belum dibatalkan
    Err() error
    
    // Value mengembalikan nilai yang terkait dengan key ini
    // nil jika key tidak ada
    Value(key any) any
}
```

### Penjelasan Detail Setiap Method

#### 1. `Deadline() (deadline time.Time, ok bool)`

```go
func explainDeadline(ctx context.Context) {
    deadline, ok := ctx.Deadline()
    
    if ok {
        fmt.Printf("Context akan expire pada: %v\n", deadline)
        timeLeft := time.Until(deadline)
        fmt.Printf("Waktu tersisa: %v\n", timeLeft)
    } else {
        fmt.Println("Context tidak memiliki deadline")
    }
}
```

**Kapan menggunakan:**
- Ketika ingin tahu berapa waktu tersisa
- Untuk memutuskan apakah masih cukup waktu untuk operasi tertentu
- Untuk logging dan debugging

#### 2. `Done() <-chan struct{}`

```go
func explainDone(ctx context.Context) {
    // Done() mengembalikan receive-only channel
    // Channel ini akan di-CLOSE (bukan di-send) ketika:
    // 1. cancel() dipanggil
    // 2. Timeout/deadline tercapai
    // 3. Parent context dibatalkan
    
    select {
    case <-ctx.Done():
        // Context sudah dibatalkan!
        fmt.Println("Context cancelled:", ctx.Err())
    default:
        // Context masih aktif
        fmt.Println("Context still active")
    }
}
```

**Pola Umum Penggunaan:**

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            // Cleanup dan exit
            return
        case <-time.After(time.Second):
            // Lakukan pekerjaan
            doWork()
        }
    }
}
```

#### 3. `Err() error`

```go
func explainErr(ctx context.Context) {
    err := ctx.Err()
    
    switch err {
    case nil:
        fmt.Println("Context masih aktif, belum dibatalkan")
    case context.Canceled:
        fmt.Println("Context dibatalkan secara manual via cancel()")
    case context.DeadlineExceeded:
        fmt.Println("Context timeout/deadline terlampaui")
    }
}
```

**Dua Jenis Error Context:**

```go
var (
    // Dikembalikan ketika cancel() dipanggil
    Canceled = errors.New("context canceled")
    
    // Dikembalikan ketika deadline/timeout tercapai
    DeadlineExceeded = errors.New("context deadline exceeded")
)
```

#### 4. `Value(key any) any`

```go
func explainValue(ctx context.Context) {
    // Value mencari key dari context saat ini ke parent-nya
    // Pencarian berhenti ketika menemukan key atau sampai root
    
    //     Background
    //         │
    //    WithValue(A=1)
    //         │
    //    WithValue(B=2)  ← ctx.Value("A") akan mencari ke atas dan menemukan 1
    //         │
    //    WithValue(A=3)  ← ctx.Value("A") akan menemukan 3 (shadowing!)
    
    value := ctx.Value("userID")
    if value != nil {
        userID := value.(string)
        fmt.Println("User ID:", userID)
    }
}
```

---

## Tipe-tipe Context

Go menyediakan beberapa fungsi untuk membuat context:

### 1. `context.Background()`

**Root context** yang paling umum digunakan. Tidak pernah dibatalkan, tidak memiliki deadline, tidak memiliki nilai.

```go
func main() {
    // Background() adalah titik awal untuk semua context chain
    ctx := context.Background()
    
    // Gunakan sebagai parent untuk context lainnya
    ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
}
```

**Kapan menggunakan:**
- Di `main()` function
- Di tests
- Sebagai parent context tingkat tertinggi

### 2. `context.TODO()`

Sama seperti `Background()`, tapi semantiknya berbeda - menandakan bahwa "belum tahu context mana yang harus digunakan".

```go
func legacyFunction() {
    // TODO: Nanti akan diganti dengan context yang proper
    ctx := context.TODO()
    newFunctionWithContext(ctx)
}
```

**Kapan menggunakan:**
- Ketika refactoring kode lama
- Placeholder sementara saat pengembangan
- Jangan gunakan di production code!

### 3. `context.WithCancel(parent Context)`

Membuat context yang dapat dibatalkan secara manual.

```go
func demonstrateWithCancel() {
    ctx := context.Background()
    
    // Membuat child context dengan kemampuan cancel
    childCtx, cancel := context.WithCancel(ctx)
    
    // PENTING: Selalu panggil cancel() untuk mencegah leak!
    defer cancel()
    
    go worker(childCtx)
    
    // Nanti, ketika ingin menghentikan worker:
    cancel() // Semua goroutine yang menggunakan childCtx akan dihentikan
}

func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Worker stopped:", ctx.Err())
            return
        default:
            fmt.Println("Working...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```

**Diagram Alur:**

```
                        ┌─────────────────┐
                        │    Background   │
                        └────────┬────────┘
                                 │
                                 ▼
                        ┌─────────────────┐
      cancel() ────────►│   WithCancel    │
                        └────────┬────────┘
                                 │
                    ┌────────────┼────────────┐
                    ▼            ▼            ▼
              ┌──────────┐ ┌──────────┐ ┌──────────┐
              │ Worker 1 │ │ Worker 2 │ │ Worker 3 │
              └──────────┘ └──────────┘ └──────────┘
              
    Ketika cancel() dipanggil, SEMUA worker berhenti!
```

### 4. `context.WithDeadline(parent Context, deadline time.Time)`

Membuat context yang akan dibatalkan pada waktu tertentu.

```go
func demonstrateWithDeadline() {
    ctx := context.Background()
    
    // Akan dibatalkan pada 16:00:00
    deadline := time.Date(2024, 1, 1, 16, 0, 0, 0, time.Local)
    deadlineCtx, cancel := context.WithDeadline(ctx, deadline)
    defer cancel()
    
    // Cek waktu tersisa
    d, ok := deadlineCtx.Deadline()
    if ok {
        fmt.Printf("Task harus selesai sebelum: %v\n", d)
        fmt.Printf("Waktu tersisa: %v\n", time.Until(d))
    }
    
    select {
    case <-deadlineCtx.Done():
        if deadlineCtx.Err() == context.DeadlineExceeded {
            fmt.Println("Deadline terlampaui!")
        }
    case result := <-doLongTask(deadlineCtx):
        fmt.Println("Task selesai:", result)
    }
}
```

### 5. `context.WithTimeout(parent Context, timeout time.Duration)`

Shortcut untuk `WithDeadline`. Lebih sering digunakan karena lebih intuitif.

```go
func demonstrateWithTimeout() {
    ctx := context.Background()
    
    // Timeout setelah 3 detik dari SEKARANG
    timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel() // SANGAT PENTING!
    
    // Equivalent dengan:
    // deadline := time.Now().Add(3 * time.Second)
    // deadlineCtx, cancel := context.WithDeadline(ctx, deadline)
    
    select {
    case <-timeoutCtx.Done():
        fmt.Println("Timeout!", timeoutCtx.Err())
    case <-time.After(3 * time.Second):
        fmt.Println("timeot atau masuk ?")
    }
}
```

**Visualisasi Timeline:**

```
Timeline:
    |──────────────────────────────────────────►
    0s         1s         2s         3s         4s         5s
    │          │          │          │          │          │
    ├──────────┴──────────┴──────────┤
    │     Context Aktif (3 detik)    │
    │                                │
    └── time.Now()                   └── Timeout! ctx.Done() di-close
```

### 6. `context.WithValue(parent Context, key, value any)`

Menyimpan key-value pair dalam context.

```go
// Definisikan custom type untuk key (best practice)
type contextKey string

const (
    userIDKey    contextKey = "userID"
    requestIDKey contextKey = "requestID"
    roleKey      contextKey = "role"
)

func demonstrateWithValue() {
    ctx := context.Background()
    
    // Menambahkan nilai ke context
    ctx = context.WithValue(ctx, userIDKey, "user-123")
    ctx = context.WithValue(ctx, requestIDKey, "req-abc-456")
    ctx = context.WithValue(ctx, roleKey, "admin")
    
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    // Mengambil nilai dari context
    userID, ok := ctx.Value(userIDKey).(string)
    if !ok {
        fmt.Println("userID tidak ditemukan")
        return
    }
    
    requestID := ctx.Value(requestIDKey).(string)
    role := ctx.Value(roleKey).(string)
    
    fmt.Printf("Processing request %s for user %s (role: %s)\n",
        requestID, userID, role)
}
```

**PERINGATAN PENTING:**

```go
// ❌ JANGAN gunakan built-in type sebagai key
ctx = context.WithValue(ctx, "userID", "123") // String key - BAD!
ctx = context.WithValue(ctx, 1, "value")       // Int key - BAD!

// ✅ Gunakan custom type
type myKey string
ctx = context.WithValue(ctx, myKey("userID"), "123") // GOOD!

// Ini mencegah collision dengan package lain yang mungkin menggunakan
// key yang sama dengan tipe berbeda
```

---

*Lanjutan di file berikutnya: `02_context_mechanism.md`*
