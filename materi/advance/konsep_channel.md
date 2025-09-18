# Channel dalam Golang

Channel adalah mekanisme komunikasi utama antar goroutine dalam Go. Channel memungkinkan pengiriman dan penerimaan nilai antar goroutine secara aman tanpa perlu explicit locking. Konsep ini didasarkan pada prinsip "Don't communicate by sharing memory; share memory by communicating."

## 1. Dasar Channel

### a. Konsep
Channel adalah conduit (saluran) bertipe yang menghubungkan goroutine. Operasi pada channel bersifat blocking:
- Mengirim ke channel memblokir pengirim sampai ada penerima
- Menerima dari channel memblokir penerima sampai ada pengirim

### b. Membuat Channel
```go
// Membuat unbuffered channel
ch := make(chan int)

// Membuat buffered channel dengan kapasitas
bufferedCh := make(chan string, 3)
```

### c. Operasi Dasar Channel
```go
// Mengirim data ke channel
ch <- 42

// Menerima data dari channel
value := <-ch

// Menerima data dan status channel
value, ok := <-ch
```

## 2. Mekanisme dan Perilaku

### a. Unbuffered Channel
Unbuffered channel tidak memiliki kapasitas buffer. Pengirim dan penerima harus "bertemu" untuk operasi berhasil.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Membuat unbuffered channel
	ch := make(chan string)
	
	// Goroutine pengirim
	go func() {
		fmt.Println("Sending message...")
		ch <- "Hello from goroutine!"
		fmt.Println("Message sent!")
	}()
	
	// Memberi waktu untuk melihat urutan
	time.Sleep(100 * time.Millisecond)
	
	// Menerima data dari channel
	message := <-ch
	fmt.Println("Received:", message)
}
```

**Hasil:**
```
Sending message...
Received: Hello from goroutine!
Message sent!
```

### b. Buffered Channel
Buffered channel memiliki kapasitas buffer. Pengiriman hanya memblokir jika buffer penuh.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Buffered channel dengan kapasitas 2
	ch := make(chan string, 2)
	
	// Mengirim data ke buffered channel (tidak memblokir sampai buffer penuh)
	ch <- "Message 1"
	fmt.Println("Sent Message 1")
	
	ch <- "Message 2"
	fmt.Println("Sent Message 2")
	
	// Buffer sekarang penuh, pengiriman berikutnya akan memblokir
	go func() {
		ch <- "Message 3" // Akan memblokir sampai ada yang menerima
		fmt.Println("Sent Message 3")
	}()
	
	// Memberi waktu untuk melihat blocking
	time.Sleep(100 * time.Millisecond)
	
	// Menerima data
	fmt.Println("Received:", <-ch)
	fmt.Println("Received:", <-ch)
	fmt.Println("Received:", <-ch)
}
```

**Hasil:**
```
Sent Message 1
Sent Message 2
Received: Message 1
Received: Message 2
Sent Message 3
Received: Message 3
```

### c. Channel Direction
Kita bisa menspesifikasikan arah channel dalam tanda tangan fungsi untuk type safety.

```go
package main

import "fmt"

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

## 3. Select Statement

`select` digunakan untuk menunggu operasi pada beberapa channel secara bersamaan.

### a. Dasar Select
```go
package main

import (
	"fmt"
	"time"
)

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

**Hasil:**
```
Received One
Received Two
```

### b. Select dengan Default Case
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "Delayed message"
	}()
	
	// Non-blocking select
	for {
		select {
		case msg := <-ch:
			fmt.Println("Received:", msg)
			return
		default:
			fmt.Println("No message yet...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
```

**Hasil:**
```
No message yet...
No message yet...
No message yet...
No message yet...
Received: Delayed message
```

### c. Select dengan Timeout
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	timeout := time.After(3 * time.Second)
	
	go func() {
		time.Sleep(5 * time.Second) // Lebih lama dari timeout
		ch <- "This message will not be received"
	}()
	
	select {
	case msg := <-ch:
		fmt.Println("Received:", msg)
	case <-timeout:
		fmt.Println("Timeout! No message received in time.")
	}
}
```

**Hasil:**
```
Timeout! No message received in time.
```

## 4. Closing Channel

Channel bisa ditutup menggunakan `close(channel)`. Ini penting untuk signal bahwa tidak akan ada lagi data dikirim.

### a. Menggunakan Close untuk Signal
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 3)
	
	// Goroutine yang mengirim dan menutup channel
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
			time.Sleep(500 * time.Millisecond)
		}
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
}
```

**Hasil:**
```
Received: 1
Received: 2
Received: 3
Received: 4
Received: 5
Channel closed
```

### b. Range Loop dengan Channel
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 3)
	
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
			time.Sleep(500 * time.Millisecond)
		}
		close(ch)
	}()
	
	// Menggunakan range loop (lebih idiomatik)
	for value := range ch {
		fmt.Println("Range received:", value)
	}
	
	fmt.Println("Channel closed, range loop finished")
}
```

## 5. Perilaku dan Gotchas

### a. Deadlock dengan Channel
```go
package main

func main() {
	ch := make(chan string)
	// Mengirim ke channel tanpa goroutine penerima
	// Akan menyebabkan deadlock
	ch <- "Hello" // Program akan hang/blokir selamanya
	<-ch
}
```

**Hasil:**
```
fatal error: all goroutines are asleep - deadlock!
```

**Notifikasi Error:** JANGAN mengirim ke unbuffered channel tanpa goroutine penerima yang siap, karena akan menyebabkan deadlock.

### b. Membaca dari Channel yang Ditutup
```go
package main

import "fmt"

func main() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	
	// Membaca dari channel yang ditutup
	for i := 0; i < 5; i++ {
		value, ok := <-ch
		fmt.Printf("Value: %d, OK: %t\n", value, ok)
	}
}
```

**Hasil:**
```
Value: 1, OK: true
Value: 2, OK: true
Value: 0, OK: false
Value: 0, OK: false
Value: 0, OK: false
```

**Notifikasi Error:** Ketika membaca dari channel yang ditutup, Anda akan mendapatkan zero value dari tipe channel dan `ok` akan bernilai `false`.

### c. Menutup Channel yang Sudah Ditutup
```go
package main

import "fmt"

func main() {
	ch := make(chan int)
	close(ch)
	
	// Menutup channel yang sudah ditutup akan panic
	// close(ch) // panic: close of closed channel
	
	// Cara aman memeriksa apakah channel ditutup
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	
	close(ch) // Ini akan panic
}
```

**Notifikasi Error:** JANGAN menutup channel yang sudah ditutup. Gunakan mekanisme sinkronisasi jika diperlukan untuk memastikan hanya satu goroutine yang menutup channel.

## 6. Best Practices

1.  **Gunakan buffered channel** jika Anda tahu kapasitas maksimum dan ingin menghindari blocking pengirim.
2.  **Gunakan unbuffered channel** untuk sinkronisasi ketat antar goroutine.
3.  **Selalu close channel** jika tidak ada lagi data yang akan dikirim, terutama jika penerima menggunakan range loop.
4.  **Gunakan `select` dengan `default`** untuk operasi non-blocking.
5.  **Gunakan `select` dengan `time.After`** untuk implementasi timeout.
6.  **Gunakan channel direction** dalam tanda tangan fungsi untuk type safety.
7.  **Hindari deadlock** dengan memastikan ada goroutine penerima untuk setiap pengirim.

## 7. Pattern: Fan-in dan Fan-out

### a. Fan-out (satu channel ke banyak goroutine)
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(time.Second) // Simulasi kerja
		fmt.Printf("Worker %d finished job %d\n", id, job)
	}
}

func main() {
	const numJobs = 5
	const numWorkers = 3
	
	jobs := make(chan int, numJobs)
	
	var wg sync.WaitGroup
	
	// Menjalankan worker (fan-out)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}
	
	// Mengirim pekerjaan
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)
	
	// Menunggu semua worker selesai
	wg.Wait()
}
```

### b. Fan-in (banyak channel ke satu channel)
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(ch chan<- int, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for i := 1; i <= 3; i++ {
		ch <- id*10 + i
		time.Sleep(500 * time.Millisecond)
	}
}

func fanIn(channels []<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	
	// Fungsi untuk menerima dari satu channel
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}
	
	// Menjalankan goroutine untuk setiap channel input
	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}
	
	// Menutup channel output setelah semua input selesai
	go func() {
		wg.Wait()
		close(out)
	}()
	
	return out
}

func main() {
	// Membuat beberapa channel input
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	
	channels := []<-chan int{ch1, ch2, ch3}
	
	// Menggabungkan channel (fan-in)
	merged := fanIn(channels)
	
	var wg sync.WaitGroup
	
	// Menjalankan producer
	wg.Add(3)
	go producer(ch1, 1, &wg)
	go producer(ch2, 2, &wg)
	go producer(ch3, 3, &wg)
	
	// Menutup channel ketika producer selesai
	go func() {
		wg.Wait()
		close(ch1)
		close(ch2)
		close(ch3)
	}()
	
	// Menerima dari channel gabungan
	for value := range merged {
		fmt.Println("Received:", value)
	}
}
```

## 8. Kesimpulan

Channel adalah fondasi dari komunikasi concurrent dalam Go. Memahami cara kerja channel sangat penting untuk:
- Membangun aplikasi concurrent yang aman dan efisien
- Menghindari race condition dan deadlock
- Mengelola alur data antar goroutine
- Membangun sistem yang scalable dan maintainable

Dengan memahami mekanisme, perilaku, dan best practices dari channel, Anda dapat menulis aplikasi Go yang robust dan idiomatik. Channel memungkinkan pendekatan "share memory by communicating" yang merupakan inti dari filosofi concurrency dalam Go.