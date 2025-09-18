# Goroutine dalam Golang

Goroutine adalah salah satu fitur paling powerful dan fundamental dalam bahasa Go. Goroutine adalah lightweight thread yang dikelola oleh runtime Go, memungkinkan eksekusi konkuren dari fungsi dengan overhead yang sangat kecil.

## 1. Dasar Goroutine

### a. Konsep Lightweight Thread
Berbeda dengan thread sistem operasi (OS thread) yang mahal untuk dibuat dan dikontrol, goroutine:
- Dibuat dengan alokasi memori awal yang sangat kecil (sekitar 2KB)
- Dikelola oleh Go runtime, bukan OS
- Bisa berjumlah ribuan tanpa masalah kinerja signifikan
- Multiplexed ke sejumlah kecil thread OS

### b. Membuat Goroutine
Gunakan kata kunci `go` diikuti dengan pemanggilan fungsi untuk membuat goroutine.

```go
package main

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

**Hasil (urutan bisa bervariasi):**
```
Hello Alice (0)
Hello Bob (0)
Hello Alice (1)
Hello Bob (1)
Hello Alice (2)
Hello Bob (2)
Main function finished
```

## 2. Mekanisme dan Perilaku

### a. Scheduler Go
Go menggunakan scheduler M:N, yang berarti M goroutine dijadwalkan ke N thread OS. Ini memungkinkan:
- Efisiensi tinggi karena tidak terikat pada thread OS
- Scalability tinggi karena bisa membuat banyak goroutine
- Preemptive scheduling untuk fairness

### b. Stack Growth
Goroutine dimulai dengan stack yang kecil (sekitar 2KB) dan bisa tumbuh dan menyusut secara dinamis sesuai kebutuhan. Ini memungkinkan pembuatan ribuan goroutine tanpa menghabiskan memori.

### c. Goroutine dan Fungsi Anonim
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Goroutine dengan fungsi anonim
	go func() {
		fmt.Println("Anonymous goroutine executing")
	}()
	
	// Goroutine dengan fungsi anonim dan parameter
	name := "Charlie"
	go func(n string) {
		fmt.Printf("Hello from anonymous goroutine: %s\n", n)
	}(name)
	
	// Memberi waktu untuk eksekusi
	time.Sleep(100 * time.Millisecond)
}
```

### d. Goroutine Closure
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Closure dengan goroutine
	for i := 0; i < 3; i++ {
		go func() {
			fmt.Printf("Closure goroutine: %d\n", i) // HATI-HATI: ini akan selalu 3
		}()
	}
	
	// Versi yang benar
	for j := 0; j < 3; j++ {
		go func(val int) {
			fmt.Printf("Correct closure goroutine: %d\n", val)
		}(j)
	}
	
	time.Sleep(100 * time.Millisecond)
}
```

**Hasil:**
```
Closure goroutine: 3
Closure goroutine: 3
Closure goroutine: 3
Correct closure goroutine: 0
Correct closure goroutine: 1
Correct closure goroutine: 2
```

## 3. Sinkronisasi Goroutine

### a. Masalah Utama: Main Function Selesai Dulu
```go
package main

import (
	"fmt"
	"time"
)

func printNumbers() {
	for i := 1; i <= 3; i++ {
		fmt.Println(i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	go printNumbers() // Menjalankan goroutine
	// Program utama selesai sebelum goroutine selesai
	fmt.Println("Main finished")
	// Tidak ada output dari goroutine karena program sudah selesai
}
```

**Hasil:**
```
Main finished
// Tidak ada output dari goroutine
```

### b. Solusi 1: Menggunakan time.Sleep (Tidak Disarankan)
```go
package main

import (
	"fmt"
	"time"
)

func printNumbers() {
	for i := 1; i <= 3; i++ {
		fmt.Println(i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	go printNumbers()
	fmt.Println("Main finished")
	time.Sleep(500 * time.Millisecond) // Memberi waktu secara paksa
}
```

### c. Solusi 2: Menggunakan sync.WaitGroup (Disarankan)
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func printNumbers(wg *sync.WaitGroup) {
	defer wg.Done() // Memberitahu WaitGroup bahwa goroutine selesai
	
	for i := 1; i <= 3; i++ {
		fmt.Println(i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	var wg sync.WaitGroup
	
	// Menambah counter WaitGroup
	wg.Add(1)
	
	// Menjalankan goroutine
	go printNumbers(&wg)
	
	fmt.Println("Main finished")
	
	// Menunggu semua goroutine selesai
	wg.Wait()
}
```

### d. Solusi 3: Menggunakan Channel untuk Sinkronisasi
```go
package main

import (
	"fmt"
	"time"
)

func printNumbers(done chan bool) {
	for i := 1; i <= 3; i++ {
		fmt.Println(i)
		time.Sleep(100 * time.Millisecond)
	}
	done <- true // Mengirim sinyal bahwa goroutine selesai
}

func main() {
	done := make(chan bool)
	
	go printNumbers(done)
	
	fmt.Println("Main finished")
	
	// Menunggu sinyal dari goroutine
	<-done
}
```

## 4. Perilaku dan Gotchas

### a. Goroutine Leak
Goroutine yang tidak pernah selesai bisa menyebabkan kebocoran memori.

```go
package main

import (
	"fmt"
	"time"
)

func leakingGoroutine() {
	for {
		// Goroutine ini tidak pernah selesai
		// Akan terus berjalan dan mengkonsumsi resource
		time.Sleep(1 * time.Second)
		fmt.Println("Still running...")
	}
}

func main() {
	go leakingGoroutine()
	
	// Program utama selesai, tapi goroutine tetap berjalan
	fmt.Println("Main finished")
	time.Sleep(3 * time.Second)
	// Setelah ini, goroutine masih berjalan di background (goroutine leak)
}
```

**Notifikasi Error:** JANGAN biarkan goroutine berjalan tanpa kontrol. Pastikan ada cara untuk memberhentikan goroutine, seperti menggunakan channel untuk sinyal berhenti atau context.

### b. Closure Variable Capture
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	
	// Salah: semua goroutine mereferensikan variabel loop yang sama
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// i di-capture oleh closure, nilainya bisa tidak sesuai harapan
			fmt.Printf("Wrong: %d\n", i) 
		}()
	}
	
	// Benar: mengirim nilai ke parameter fungsi
	for j := 0; j < 3; j++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			fmt.Printf("Correct: %d\n", val)
		}(j)
	}
	
	wg.Wait()
	time.Sleep(100 * time.Millisecond)
}
```

### c. Goroutine dan Error Handling
Goroutine tidak bisa mengembalikan error secara langsung seperti fungsi biasa.

```go
package main

import (
	"fmt"
	"sync"
)

// Salah: error tidak bisa dikembalikan dari goroutine
func processWithErrorWrong() error {
	go func() {
		// Simulasi error
		// return fmt.Errorf("something went wrong") // ERROR: cannot return from goroutine
	}()
	return nil
}

// Benar: menggunakan channel untuk mengirim error
func processWithError(done chan error) {
	// Simulasi error
	err := fmt.Errorf("something went wrong")
	done <- err // Mengirim error melalui channel
}

func main() {
	done := make(chan error, 1) // Buffered channel
	
	go processWithError(done)
	
	// Menerima error dari goroutine
	if err := <-done; err != nil {
		fmt.Printf("Error from goroutine: %v\n", err)
	}
}
```

## 5. Best Practices

1.  **Gunakan sync.WaitGroup atau channel untuk sinkronisasi**, bukan `time.Sleep`.
2.  **Hindari goroutine leak** dengan memastikan ada cara untuk memberhentikan goroutine.
3.  **Gunakan parameter dalam closure** daripada mereferensikan variabel loop langsung.
4.  **Gunakan buffered channel** jika tidak ingin goroutine memblokir saat mengirim data.
5.  **Hindari membuat goroutine yang tidak dibutuhkan** - goroutine ringan tapi bukan gratis.
6.  **Gunakan context untuk manajemen lifetime** goroutine dalam aplikasi yang kompleks.

## 6. Pattern: Worker Pool

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	
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
	
	var wg sync.WaitGroup
	
	// Menjalankan worker
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}
	
	// Mengirim pekerjaan
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Memberitahu worker bahwa tidak ada pekerjaan lagi
	
	// Menunggu semua worker selesai
	wg.Wait()
	close(results)
	
	// Mencetak hasil
	for a := 1; a <= numJobs; a++ {
		fmt.Printf("Result: %d\n", <-results)
	}
}
```

## 7. Kesimpulan

Goroutine adalah fondasi dari concurrency di Go. Memahami cara kerja goroutine sangat penting untuk:
- Membangun aplikasi yang responsif dan efisien
- Memanfaatkan multi-core processor secara maksimal
- Mengelola banyak operasi I/O secara bersamaan

Dengan memahami mekanisme, perilaku, dan best practices dari goroutine, Anda dapat menulis aplikasi concurrent yang scalable dan maintainable. Ingatlah untuk selalu memperhatikan sinkronisasi dan manajemen lifetime goroutine untuk menghindari masalah seperti goroutine leak dan race condition.