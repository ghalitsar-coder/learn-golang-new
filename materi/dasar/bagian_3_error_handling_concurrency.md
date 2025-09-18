# Bagian 3: Error Handling, Defer, dan Concurrency Dasar

## 1. Error Handling

Error handling di Go dilakukan secara eksplisit dan merupakan bagian penting dari alur program. Berbeda dengan bahasa lain yang menggunakan exception, Go menggunakan tipe `error` sebagai return value.

### a. Tipe `error`
Tipe `error` adalah interface bawaan:
```go
type error interface {
	Error() string
}
```

### b. Mengembalikan Error
Fungsi yang bisa gagal biasanya memiliki tipe return `error` sebagai nilai kembalian terakhir.

Contoh:
```go
import (
	"errors"
	"fmt"
	"strconv"
)

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero is not allowed")
	}
	return a / b, nil
}

func main() {
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
	
	// Kasus error
	result2, err2 := divide(10, 0)
	if err2 != nil {
		fmt.Println("Error:", err2)
	} else {
		fmt.Println("Result:", result2)
	}
}
```

### c. Membuat Error Kustom
Menggunakan `fmt.Errorf` untuk membuat error dengan format string:
```go
import (
	"fmt"
)

func validateAge(age int) error {
	if age < 0 {
		return fmt.Errorf("invalid age: %d, age cannot be negative", age)
	}
	if age > 150 {
		return fmt.Errorf("invalid age: %d, age seems unrealistic", age)
	}
	return nil // Tidak ada error
}

func main() {
	err := validateAge(-5)
	if err != nil {
		fmt.Println("Validation Error:", err)
	}
}
```

### d. Error Checking dan Best Practices
- Selalu periksa error setelah fungsi yang bisa gagal dipanggil.
- Gunakan `errors.Is` dan `errors.As` (Go 1.13+) untuk pemeriksaan error yang lebih robust.
```go
import (
	"errors"
	"fmt"
	"os"
)

func main() {
	_, err := os.Open("nonexistent.txt")
	if err != nil {
		// Memeriksa error spesifik
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("File does not exist")
		} else {
			fmt.Println("Error opening file:", err)
		}
	}
}
```

## 2. Defer, Panic, dan Recover

### a. Defer
`defer` digunakan untuk menunda eksekusi sebuah fungsi hingga fungsi sekitarnya selesai. Ini sangat berguna untuk pembersihan sumber daya seperti menutup file atau koneksi database.

Contoh:
```go
func main() {
	defer fmt.Println("This will be printed last")
	fmt.Println("This will be printed first")
	
	// Multiple defer (LIFO - Last In, First Out)
	defer fmt.Println("Deferred 1")
	defer fmt.Println("Deferred 2")
	fmt.Println("This will be printed second")
	// Output:
	// This will be printed first
	// This will be printed second
	// Deferred 2
	// Deferred 1
	// This will be printed last
}
```

Contoh penggunaan praktis:
```go
func readFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close() // Menjamin file akan ditutup
	
	// Lakukan operasi pada file
	// File akan otomatis ditutup saat fungsi selesai
}
```

### b. Panic
`panic` digunakan untuk menghentikan alur normal program ketika terjadi kesalahan serius yang tidak bisa atau tidak seharusnya ditangani secara normal. Ini seperti melempar exception.

Contoh:
```go
func main() {
	fmt.Println("Starting program")
	panic("Something went wrong!")
	fmt.Println("This will not be printed")
}
```

### c. Recover
`recover` digunakan untuk menangkap `panic` dan memungkinkan program melanjutkan eksekusi. `recover` hanya berguna dalam fungsi yang di-defer.

Contoh:
```go
func recoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from panic:", r)
	}
}

func mightPanic() {
	defer recoverFromPanic()
	fmt.Println("About to panic")
	panic("This is a panic!")
	fmt.Println("This will not be printed")
}

func main() {
	mightPanic()
	fmt.Println("Program continues after panic recovery")
}
```

Best Practice:
- Gunakan `panic` dan `recover` dengan bijak. Error handling biasa (`error` return value) lebih disukai.
- `panic` umumnya digunakan untuk kesalahan pemrograman (bug), bukan kesalahan runtime yang bisa diprediksi.

## 3. Goroutine

Goroutine adalah lightweight thread yang dikelola oleh runtime Go. Goroutine sangat murah untuk dibuat dan bisa berjumlah ribuan tanpa masalah kinerja signifikan.

### a. Membuat Goroutine
Gunakan kata kunci `go` diikuti dengan pemanggilan fungsi untuk membuat goroutine.

Contoh:
```go
import (
	"fmt"
	"time"
)

func sayHello(name string) {
	for i := 0; i < 3; i++ {
		fmt.Printf("Hello %s (%d)\n", name, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// Menjalankan fungsi dalam goroutine
	go sayHello("Alice")
	
	// Menjalankan fungsi dalam goroutine lain
	go sayHello("Bob")
	
	// Memberi waktu agar goroutine selesai dieksekusi
	// (dalam praktik nyata, biasanya menggunakan channel atau WaitGroup)
	time.Sleep(1 * time.Second)
	fmt.Println("Main function finished")
}
```

### b. Perbedaan Goroutine dan Thread OS
- Goroutine dikelola oleh runtime Go, sedangkan thread dikelola oleh OS.
- Goroutine jauh lebih ringan, bisa dibuat dalam jumlah besar.
- Goroutine multiplexed ke sejumlah kecil thread OS.

## 4. Channel

Channel digunakan untuk komunikasi dan sinkronisasi antar goroutine. Channel memungkinkan pengiriman dan penerimaan nilai antar goroutine secara aman.

### a. Membuat Channel
```go
// Membuat channel untuk tipe int
ch := make(chan int)

// Membuat buffered channel dengan kapasitas 3
bufferedCh := make(chan string, 3)
```

### b. Mengirim dan Menerima Data
```go
// Mengirim data ke channel
ch <- 42

// Menerima data dari channel
value := <-ch
```

Contoh sederhana:
```go
func main() {
	// Membuat channel
	ch := make(chan string)
	
	// Goroutine pengirim
	go func() {
		ch <- "Hello from goroutine!"
	}()
	
	// Menerima data dari channel
	message := <-ch
	fmt.Println(message)
}
```

### c. Buffered Channel
Buffered channel memiliki kapasitas buffer. Pengiriman ke buffered channel hanya akan memblokir jika buffer penuh.

```go
func main() {
	// Buffered channel dengan kapasitas 2
	ch := make(chan int, 2)
	
	ch <- 1
	ch <- 2
	// ch <- 3 // Ini akan memblokir karena buffer penuh
	
	fmt.Println(<-ch) // Output: 1
	fmt.Println(<-ch) // Output: 2
	// fmt.Println(<-ch) // Ini akan memblokir karena tidak ada data
}
```

### d. Channel Direction
Kita bisa menspesifikasikan arah channel (hanya kirim atau hanya terima) dalam tanda tangan fungsi untuk type safety.

```go
// Channel hanya bisa menerima data
func sender(ch chan<- string) {
	ch <- "Message from sender"
}

// Channel hanya bisa menerima data
func receiver(ch <-chan string) {
	message := <-ch
	fmt.Println("Received:", message)
}

func main() {
	ch := make(chan string)
	
	go sender(ch)
	receiver(ch)
}
```

### e. Select Statement
`select` digunakan untuk menunggu operasi pada beberapa channel secara bersamaan.

```go
func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "One"
	}()
	
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "Two"
	}()
	
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received", msg2)
		}
	}
}
```

### f. Closing Channel
Channel bisa ditutup menggunakan `close(channel)`. Penerima bisa memeriksa apakah channel sudah ditutup.

```go
func main() {
	ch := make(chan int, 3)
	
	go func() {
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch) // Menutup channel
	}()
	
	// Menerima data dengan pemeriksaan channel closed
	for {
		value, ok := <-ch
		if !ok {
			fmt.Println("Channel closed")
			break
		}
		fmt.Println("Received:", value)
	}
	
	// Atau menggunakan range loop (lebih idiomatik)
	// ch2 := make(chan int, 3)
	// go func() {
	// 	ch2 <- 10
	// 	ch2 <- 20
	// 	close(ch2)
	// }()
	// 
	// for value := range ch2 {
	// 	fmt.Println("Range received:", value)
	// }
}
```

Goroutine dan channel adalah fondasi dari concurrency di Go, memungkinkan penulisan program yang efisien dan scalable.