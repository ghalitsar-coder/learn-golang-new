# Concurrency dalam Golang

Concurrency adalah kemampuan program untuk menangani beberapa komputasi secara bersamaan. Dalam konteks Go, concurrency berarti kemampuan untuk menjalankan beberapa goroutine secara bersamaan dengan cara yang aman dan efisien. Go memiliki dukungan concurrency bawaan yang sangat kuat melalui goroutine, channel, dan package sync.

## 1. Dasar Concurrency

### a. Konsep Concurrency vs Parallelism
- **Concurrency**: Menangani banyak tugas secara bersamaan (dalam waktu yang tumpang tindih). Ini tentang struktur dan komposisi program.
- **Parallelism**: Menjalankan banyak tugas secara bersamaan (secara fisik pada beberapa core). Ini tentang eksekusi simultan.

```go
// Contoh concurrency tanpa parallelism
// Satu core menangani banyak goroutine dengan context switching
package main

import (
	"fmt"
	"time"
)

func task(name string) {
	for i := 1; i <= 3; i++ {
		fmt.Printf("Task %s: step %d\n", name, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// Ini concurrency: beberapa tugas diatur untuk berjalan bersamaan
	// Meskipun mungkin hanya menggunakan satu core
	go task("A")
	go task("B")
	
	time.Sleep(1 * time.Second)
}
```

### b. Model CSP (Communicating Sequential Processes)
Go mengadopsi model CSP untuk concurrency:
- Goroutine adalah proses sekuen yang ringan
- Channel adalah media komunikasi antar goroutine
- Prinsip: "Don't communicate by sharing memory; share memory by communicating"

## 2. Mekanisme Concurrency dalam Go

### a. Goroutine
Unit eksekusi ringan yang dikelola oleh Go runtime.

```go
package main

import (
	"fmt"
	"time"
)

func backgroundTask(id int) {
	for i := 0; i < 5; i++ {
		fmt.Printf("Goroutine %d: step %d\n", id, i)
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	// Membuat beberapa goroutine
	for i := 1; i <= 3; i++ {
		go backgroundTask(i)
	}
	
	// Memberi waktu untuk eksekusi
	time.Sleep(2 * time.Second)
}
```

### b. Channel
Mekanisme komunikasi utama antar goroutine.

```go
package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- string, name string) {
	for i := 1; i <= 3; i++ {
		msg := fmt.Sprintf("%s-%d", name, i)
		ch <- msg
		time.Sleep(100 * time.Millisecond)
	}
}

func consumer(ch <-chan string, name string) {
	for msg := range ch {
		fmt.Printf("Consumer %s received: %s\n", name, msg)
	}
}

func main() {
	ch := make(chan string)
	
	// Producer dan consumer berjalan secara concurrent
	go producer(ch, "Producer1")
	go consumer(ch, "Consumer1")
	
	time.Sleep(1 * time.Second)
}
```

### c. Select Statement
Memungkinkan menunggu operasi pada beberapa channel secara bersamaan.

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
		ch1 <- "From channel 1"
	}()
	
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "From channel 2"
	}()
	
	// Menunggu dari dua channel secara bersamaan
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received:", msg2)
		}
	}
}
```

## 3. Synchronization dan Coordination

### a. sync.WaitGroup
Untuk menunggu sekelompok goroutine selesai.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	var wg sync.WaitGroup
	
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	
	fmt.Println("Waiting for workers...")
	wg.Wait()
	fmt.Println("All workers done")
}
```

### b. Mutex
Untuk melindungi akses ke data bersama.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	counter := &SafeCounter{}
	
	var wg sync.WaitGroup
	
	// Beberapa goroutine mengakses counter secara concurrent
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.Increment()
				time.Sleep(time.Millisecond)
			}
		}()
	}
	
	wg.Wait()
	fmt.Printf("Final counter value: %d\n", counter.Value())
}
```

### c. sync.Once
Untuk memastikan operasi hanya dilakukan sekali.

```go
package main

import (
	"fmt"
	"sync"
	"time
)

var (
	once sync.Once
	data string
)

func initialize() {
	fmt.Println("Initializing...")
	time.Sleep(100 * time.Millisecond)
	data = "Initialized data"
	fmt.Println("Initialization complete")
}

func getData() string {
	once.Do(initialize)
	return data
}

func main() {
	var wg sync.WaitGroup
	
	// Banyak goroutine memanggil getData secara concurrent
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			result := getData()
			fmt.Printf("Goroutine %d got: %s\n", id, result)
		}(i)
	}
	
	wg.Wait()
}
```

## 4. Pattern dan Best Practices

### a. Worker Pool Pattern
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, j)
		time.Sleep(time.Second) // Simulasi kerja
		fmt.Printf("Worker %d finished job %d\n", id, j)
		results <- j * 2
	}
}

func main() {
	const numJobs = 5
	const numWorkers = 3
	
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	
	// Menjalankan worker
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}
	
	// Mengirim pekerjaan
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)
	
	// Mengumpulkan hasil
	for a := 1; a <= numJobs; a++ {
		result := <-results
		fmt.Printf("Result: %d\n", result)
	}
}
```

### b. Pipeline Pattern
```go
package main

import "fmt"

// Stage 1: Generator
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// Stage 2: Processor
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// Stage 3: Consumer
func consumer(in <-chan int) {
	for n := range in {
		fmt.Println(n)
	}
}

func main() {
	// Membangun pipeline
	nums := generator(1, 2, 3, 4, 5)
	squares := square(nums)
	
	// Mengkonsumsi hasil
	consumer(squares)
}
```

### c. Fan-in dan Fan-out
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(out chan<- int, id int) {
	for i := 1; i <= 3; i++ {
		out <- id*10 + i
		time.Sleep(100 * time.Millisecond)
	}
	close(out)
}

func fanIn(channels []<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	
	// Goroutine untuk setiap channel input
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for n := range c {
				out <- n
			}
		}(ch)
	}
	
	// Menutup channel output ketika semua selesai
	go func() {
		wg.Wait()
		close(out)
	}()
	
	return out
}

func main() {
	// Membuat beberapa channel
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	
	// Menjalankan producer
	go producer(ch1, 1)
	go producer(ch2, 2)
	go producer(ch3, 3)
	
	// Fan-in
	channels := []<-chan int{ch1, ch2, ch3}
	merged := fanIn(channels)
	
	// Menerima dari channel gabungan
	for value := range merged {
		fmt.Println("Received:", value)
	}
}
```

## 5. Perilaku dan Gotchas

### a. Race Condition
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var counter int

func unsafeIncrement(wg *sync.WaitGroup) {
	defer wg.Done()
	temp := counter
	time.Sleep(1 * time.Microsecond) // Simulasi proses
	counter = temp + 1
}

func main() {
	var wg sync.WaitGroup
	
	// Menjalankan beberapa goroutine yang mengakses counter secara concurrent
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go unsafeIncrement(&wg)
	}
	
	wg.Wait()
	fmt.Printf("Unsafe final counter: %d\n", counter) // Hasil tidak konsisten
}
```

**Notifikasi Error:** JANGAN mengakses data bersama dari beberapa goroutine tanpa sinkronisasi yang tepat. Ini akan menyebabkan race condition.

### b. Goroutine Leak
```go
package main

import (
	"fmt"
	"time"
)

func leakingGoroutine() {
	for {
		// Goroutine ini tidak pernah selesai
		time.Sleep(1 * time.Second)
		fmt.Println("Still running...")
	}
}

func main() {
	go leakingGoroutine()
	
	// Program utama selesai, tapi goroutine tetap berjalan
	fmt.Println("Main finished")
	time.Sleep(3 * time.Second)
	// Setelah ini, goroutine masih berjalan (goroutine leak)
}
```

**Notifikasi Error:** JANGAN biarkan goroutine berjalan tanpa kontrol. Pastikan ada cara untuk memberhentikan goroutine, seperti menggunakan channel untuk sinyal berhenti atau context.

### c. Deadlock
```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex
	
	// Mencoba mengunci mutex yang sama dua kali dalam satu goroutine
	mu.Lock()
	fmt.Println("First lock acquired")
	
	// Ini akan menyebabkan deadlock
	mu.Lock() // Tidak bisa mengunci lagi karena sudah dikunci
	fmt.Println("Second lock acquired")
}
```

**Notifikasi Error:** JANGAN mencoba mengunci mutex yang sama dua kali dalam satu goroutine karena akan menyebabkan deadlock.

## 6. Error Handling dalam Concurrent Programming

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// Fungsi yang bisa gagal dalam goroutine
func riskyTask(id int) error {
	// Simulasi kemungkinan error
	if id%2 == 0 {
		return fmt.Errorf("task %d failed", id)
	}
	
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Task %d completed successfully\n", id)
	return nil
}

func main() {
	const numTasks = 5
	
	// Channel untuk mengirim error
	errCh := make(chan error, numTasks)
	var wg sync.WaitGroup
	
	// Menjalankan tugas secara concurrent
	for i := 1; i <= numTasks; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()
			if err := riskyTask(taskID); err != nil {
				errCh <- err // Mengirim error ke channel
			}
		}(i)
	}
	
	// Goroutine untuk menutup channel error ketika semua selesai
	go func() {
		wg.Wait()
		close(errCh)
	}()
	
	// Mengumpulkan error
	for err := range errCh {
		fmt.Printf("Error occurred: %v\n", err)
	}
}
```

## 7. Monitoring dan Debugging Concurrency

### a. Menggunakan runtime.NumGoroutine()
```go
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Printf("Initial goroutines: %d\n", runtime.NumGoroutine())
	
	// Menjalankan beberapa goroutine
	for i := 0; i < 5; i++ {
		go func(id int) {
			time.Sleep(2 * time.Second)
			fmt.Printf("Goroutine %d finished\n", id)
		}(i)
	}
	
	fmt.Printf("After starting goroutines: %d\n", runtime.NumGoroutine())
	
	time.Sleep(3 * time.Second)
	fmt.Printf("After completion: %d\n", runtime.NumGoroutine())
}
```

### b. Menggunakan pprof untuk profiling
```go
// Untuk mengaktifkan profiling, tambahkan ini di main():
// import _ "net/http/pprof"
// import "net/http"
// go func() {
// 	http.ListenAndServe(":6060", nil)
// }()

// Kemudian gunakan:
// go tool pprof http://localhost:6060/debug/pprof/goroutine
```

## 8. Best Practices

1.  **Gunakan channel untuk komunikasi** antar goroutine daripada shared memory jika memungkinkan.
2.  **Gunakan mutex atau atomic** untuk melindungi akses ke data bersama.
3.  **Hindari goroutine leak** dengan memastikan ada cara untuk memberhentikan goroutine.
4.  **Gunakan sync.WaitGroup** untuk menunggu sekelompok goroutine selesai.
5.  **Gunakan context** untuk manajemen lifetime dan cancellation.
6.  **Gunakan select dengan default** untuk operasi non-blocking.
7.  **Gunakan buffered channel** jika Anda tahu kapasitas maksimum.
8.  **Tutup channel** ketika tidak ada lagi data yang akan dikirim.
9.  **Hindari nested lock** untuk mencegah deadlock.
10. **Gunakan defer** untuk memastikan unlock atau cleanup selalu dilakukan.

## 9. Kesimpulan

Concurrency dalam Go adalah pendekatan yang sangat powerful dan idiomatik untuk membangun aplikasi yang responsif dan efisien. Dengan memahami:
- Mekanisme goroutine, channel, dan sync package
- Pattern seperti worker pool, pipeline, fan-in/fan-out
- Perilaku dan gotchas dalam pemrograman concurrent
- Best practices untuk menulis kode concurrent yang aman dan efisien

Anda dapat membangun aplikasi Go yang scalable, maintainable, dan dapat memanfaatkan penuh kemampuan multi-core processor modern. Concurrency bukan hanya tentang menjalankan banyak hal secara bersamaan, tetapi juga tentang mengelola kompleksitas dan memastikan keselamatan data dalam lingkungan concurrent.