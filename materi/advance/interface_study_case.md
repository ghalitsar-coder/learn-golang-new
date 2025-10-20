# Interface dalam Go - Studi Kasus Praktis

## Apa itu Interface?
Interface dalam Go adalah tipe data yang mendefinisikan sekumpulan method signature (nama method, parameter, dan return type). Interface digunakan untuk menyamakan perilaku dari berbagai jenis tipe data.

**Penting:** Di Go, interface diimplementasikan secara **implisit**. Artinya, tidak perlu secara eksplisit menyatakan bahwa suatu struct mengimplementasi interface - cukup dengan memiliki method-method yang diminta oleh interface.

## Konsep Dasar Interface

Mari kita mulai dengan contoh sederhana:

```go
package main

import "fmt"

// 1. Mendefinisikan interface
type Animal interface {
    Speak() string
    Move() string
}

// 2. Membuat struct yang berbeda
type Dog struct {
    Name string
}

type Cat struct {
    Name string
}

type Bird struct {
    Name string
}

// 3. Menambahkan method-method yang sesuai dengan interface
func (d Dog) Speak() string {
    return d.Name + " menggonggong!"
}

func (d Dog) Move() string {
    return d.Name + " berlari!"
}

func (c Cat) Speak() string {
    return c.Name + " mengeong!"
}

func (c Cat) Move() string {
    return c.Name + " berjalan lembut!"
}

func (b Bird) Speak() string {
    return b.Name + " berkicau!"
}

func (b Bird) Move() string {
    return b.Name + " terbang!"
}

func main() {
    // 4. Menggunakan interface
    animals := []Animal{
        Dog{Name: "Buddy"},
        Cat{Name: "Whiskers"},
        Bird{Name: "Tweety"},
    }

    for _, animal := range animals {
        fmt.Println(animal.Speak())
        fmt.Println(animal.Move())
        fmt.Println("---")
    }
}
```

**Output:**
```
Buddy menggonggong!
Buddy berlari!
---
Whiskers mengeong!
Whiskers berjalan lembut!
---
Tweety berkicau!
Tweety terbang!
---
```

## Studi Kasus: Sistem Pembayaran

Sekarang mari kita lihat studi kasus dunia nyata: sistem pembayaran yang bisa menangani berbagai jenis metode pembayaran.

```go
package main

import (
    "fmt"
    "time"
)

// 1. Interface untuk metode pembayaran
type PaymentMethod interface {
    ProcessPayment(amount float64) (string, error)
    GetPaymentInfo() string
}

// 2. Berbagai implementasi struct untuk metode pembayaran
type CreditCard struct {
    CardNumber string
    ExpiryDate string
    CVV        string
}

type PayPal struct {
    Email string
    Password string
}

type BankTransfer struct {
    AccountNumber string
    BankName      string
}

// 3. Implementasi method untuk CreditCard
func (cc CreditCard) ProcessPayment(amount float64) (string, error) {
    // Simulasi pemrosesan pembayaran kartu kredit
    transactionID := fmt.Sprintf("CC-%d", time.Now().Unix())
    fmt.Printf("Memroses pembayaran kartu kredit sebesar $%.2f\n", amount)
    return transactionID, nil
}

func (cc CreditCard) GetPaymentInfo() string {
    return fmt.Sprintf("Kartu Kredit dengan nomor ending in %s", cc.CardNumber[len(cc.CardNumber)-4:])
}

// 4. Implementasi method untuk PayPal
func (pp PayPal) ProcessPayment(amount float64) (string, error) {
    // Simulasi pemrosesan pembayaran PayPal
    transactionID := fmt.Sprintf("PP-%d", time.Now().Unix())
    fmt.Printf("Memroses pembayaran PayPal sebesar $%.2f dari %s\n", amount, pp.Email)
    return transactionID, nil
}

func (pp PayPal) GetPaymentInfo() string {
    return fmt.Sprintf("Akun PayPal %s", pp.Email)
}

// 5. Implementasi method untuk BankTransfer
func (bt BankTransfer) ProcessPayment(amount float64) (string, error) {
    // Simulasi pemrosesan pembayaran transfer bank
    transactionID := fmt.Sprintf("BT-%d", time.Now().Unix())
    fmt.Printf("Memroses transfer bank sebesar $%.2f ke %s (Rekening: %s)\n", amount, bt.BankName, bt.AccountNumber)
    return transactionID, nil
}

func (bt BankTransfer) GetPaymentInfo() string {
    return fmt.Sprintf("Transfer Bank ke %s (Rekening: %s)", bt.BankName, bt.AccountNumber)
}

// 6. Fungsi yang menggunakan interface - ini bisa menerima semua tipe pembayaran
func ProcessOrder(paymentMethod PaymentMethod, amount float64) {
    fmt.Println("Menggunakan metode pembayaran:", paymentMethod.GetPaymentInfo())
    
    transactionID, err := paymentMethod.ProcessPayment(amount)
    if err != nil {
        fmt.Printf("Gagal memproses pembayaran: %v\n", err)
        return
    }
    
    fmt.Printf("Pembayaran berhasil! ID Transaksi: %s\n", transactionID)
    fmt.Println("---")
}

func main() {
    // 7. Berbagai metode pembayaran
    creditCard := CreditCard{
        CardNumber: "1234567890123456",
        ExpiryDate: "12/25",
        CVV:        "123",
    }
    
    paypal := PayPal{
        Email:    "user@example.com",
        Password: "password", // dalam dunia nyata, jangan simpan password di sini!
    }
    
    bankTransfer := BankTransfer{
        AccountNumber: "987654321",
        BankName:      "Bank ABC",
    }
    
    // 8. Proses pesanan dengan berbagai metode pembayaran - semua menggunakan fungsi yang sama!
    ProcessOrder(creditCard, 99.99)
    ProcessOrder(paypal, 49.99)
    ProcessOrder(bankTransfer, 149.50)
}
```

**Mengapa ini sangat berguna?**

1. **Fleksibilitas**: Kita bisa menambahkan metode pembayaran baru tanpa mengubah fungsi `ProcessOrder`
2. **Testability**: Kita bisa membuat mock object untuk testing
3. **Abstraksi**: Detail implementasi disembunyikan, yang penting adalah kontrak (interface)

## Studi Kasus: Sistem Logging

Mari kita lihat contoh lain dengan sistem logging yang bisa menulis ke berbagai tujuan:

```go
package main

import (
    "fmt"
    "os"
    "time"
)

// Interface logger
type Logger interface {
    Log(message string)
}

// Implementasi file logger
type FileLogger struct {
    FileName string
}

func (fl FileLogger) Log(message string) {
    f, err := os.OpenFile(fl.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer f.Close()
    
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
    f.WriteString(logEntry)
}

// Implementasi console logger
type ConsoleLogger struct {
    Prefix string
}

func (cl ConsoleLogger) Log(message string) {
    timestamp := time.Now().Format("15:04:05")
    fmt.Printf("[%s %s] %s\n", cl.Prefix, timestamp, message)
}

// Implementasi network logger (simulasi)
type NetworkLogger struct {
    Host string
    Port string
}

func (nl NetworkLogger) Log(message string) {
    // Dalam dunia nyata, ini akan mengirim data ke server logging
    fmt.Printf("[NETWORK] Mengirim log ke %s:%s - %s\n", nl.Host, nl.Port, message)
}

// Fungsi yang menggunakan logger
func ProcessData(logger Logger) {
    logger.Log("Memulai proses data...")
    // Simulasi proses data
    logger.Log("Proses data selesai!")
}

func main() {
    // Bisa menggunakan logger apapun! 
    fileLogger := FileLogger{FileName: "app.log"}
    consoleLogger := ConsoleLogger{Prefix: "APP"}
    networkLogger := NetworkLogger{Host: "logs.server.com", Port: "514"}
    
    // Semua ini bekerja karena interface!
    ProcessData(fileLogger)    // Log ke file
    ProcessData(consoleLogger) // Log ke konsol
    ProcessData(networkLogger) // Log ke jaringan
}
```

## Interface Kosong (Empty Interface) - `interface{}`

Interface kosong (`interface{}`) adalah interface yang tidak memiliki method sama sekali. Di Go, **semua tipe data mengimplementasi interface kosong**.

```go
package main

import "fmt"

func PrintAnything(value interface{}) {
    fmt.Println("Nilai:", value)
    fmt.Printf("Tipe: %T\n", value)
}

func main() {
    // Bisa menerima tipe apapun!
    PrintAnything(42)
    PrintAnything(3.14)
    PrintAnything("Hello")
    PrintAnything(true)
}
```

## Type Assertion

Kadang kita perlu mengakses nilai sebenarnya dari interface:

```go
package main

import "fmt"

func main() {
    var value interface{} = "Hello, Go!"

    // Type assertion - aman
    if str, ok := value.(string); ok {
        fmt.Printf("Nilai adalah string: %s\n", str)
    }

    // Type assertion - bisa panic jika salah
    str := value.(string)
    fmt.Printf("Nilai: %s\n", str)
}
```

## Praktik Terbaik

1. **Interface seharusnya kecil**: Idealnya hanya berisi 1-3 method
2. **Gunakan nama yang berakhiran -er**: `Reader`, `Writer`, `Stringer`
3. **Tulis interface dari kebutuhan pengguna, bukan implementasi**
4. **Gunakan interface untuk abstraksi, bukan untuk semua hal**

## Kesimpulan

Interface memungkinkan:
- **Polymorphism**: Menangani berbagai tipe secara konsisten
- **Abstraction**: Menyembunyikan detail implementasi
- **Testing**: Mudah membuat mock object
- **Fleksibilitas**: Mudah mengganti implementasi

Dengan interface, kita bisa menulis kode yang lebih modular, lebih mudah diuji, dan lebih fleksibel!

## Latihan

Coba buat interface `Shape` dengan metode `Area()` dan `Perimeter()`, lalu implementasikan untuk `Rectangle`, `Circle`, dan `Triangle`. Buat fungsi yang bisa menerima berbagai bentuk dan menampilkan informasinya.