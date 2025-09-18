# Mutex dalam Golang

Mutex (Mutual Exclusion) adalah mekanisme sinkronisasi yang digunakan untuk melindungi bagian kode kritis (critical section) agar hanya dapat diakses oleh satu goroutine pada satu waktu. Ini mencegah race condition dan memastikan data consistency dalam lingkungan concurrent.

## 1. Dasar Mutex

### a. Konsep Race Condition
Race condition terjadi ketika dua atau lebih goroutine mengakses dan memodifikasi data yang sama secara konkuren tanpa sinkronisasi yang tepat, menghasilkan nilai akhir yang tidak dapat diprediksi.

### b. Jenis Mutex dalam Go
Go menyediakan dua jenis mutex dalam package `sync`:
- `sync.Mutex`: Mutex dasar yang menyediakan lock eksklusif.
- `sync.RWMutex`: Reader-Writer Mutex yang memungkinkan multiple reader atau satu writer.

## 2. Mekanisme dan Perilaku sync.Mutex

### a. Penggunaan Dasar sync.Mutex
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	counter int
	mutex   sync.Mutex
)

func increment(name string, iterations int) {
	for i := 0; i < iterations; i++ {
		// Mengunci sebelum mengakses data bersama
		mutex.Lock()
		
		// Bagian kode kritis
		temp := counter
		time.Sleep(1 * time.Microsecond) // Simulasi proses kompleks
		counter = temp + 1
		
		// Membuka kunci setelah selesai
		mutex.Unlock()
		
		fmt.Printf("%s: counter = %d\n", name, counter)
	}
}

func main() {
	// Menjalankan beberapa goroutine
	go increment("Goroutine-1", 5)
	go increment("Goroutine-2", 5)
	
	// Memberi waktu untuk eksekusi
	time.Sleep(2 * time.Second)
	fmt.Printf("Final counter value: %d\n", counter)
}
```

**Hasil (dapat bervariasi):**
```
Goroutine-1: counter = 1
Goroutine-2: counter = 2
Goroutine-1: counter = 3
Goroutine-2: counter = 4
Goroutine-1: counter = 5
Goroutine-2: counter = 6
Goroutine-1: counter = 7
Goroutine-2: counter = 8
Goroutine-1: counter = 9
Goroutine-2: counter = 10
Final counter value: 10
```

### b. Tanpa Mutex (Race Condition)
```go
package main

import (
	"fmt"
	"time"
)

var counter int // Data bersama tanpa proteksi

func incrementUnsafe(name string, iterations int) {
	for i := 0; i < iterations; i++ {
		// Tidak ada sinkronisasi
		temp := counter
		time.Sleep(1 * time.Microsecond) // Simulasi proses kompleks
		counter = temp + 1
		
		fmt.Printf("%s: counter = %d\n", name, counter)
	}
}

func main() {
	go incrementUnsafe("Goroutine-1", 5)
	go incrementUnsafe("Goroutine-2", 5)
	
	time.Sleep(2 * time.Second)
	fmt.Printf("Final counter value (unsafe): %d\n", counter)
}
```

**Hasil (akan berbeda setiap kali dijalankan):**
```
Goroutine-1: counter = 1
Goroutine-2: counter = 1  // Race condition!
Goroutine-1: counter = 2
Goroutine-2: counter = 3
Goroutine-1: counter = 4
Goroutine-2: counter = 5
Goroutine-1: counter = 6
Goroutine-2: counter = 7
Goroutine-1: counter = 8
Goroutine-2: counter = 9
Final counter value (unsafe): 9  // Tidak konsisten
```

## 3. sync.RWMutex (Reader-Writer Mutex)

### a. Konsep
RWMutex memungkinkan:
- Multiple readers membaca data secara bersamaan
- Hanya satu writer yang dapat menulis (dan tidak ada reader yang dapat membaca saat itu)

### b. Penggunaan RWMutex
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type DataStore struct {
	data  map[string]string
	mutex sync.RWMutex
}

func NewDataStore() *DataStore {
	return &DataStore{
		data: make(map[string]string),
	}
}

func (ds *DataStore) Write(key, value string) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	
	fmt.Printf("Writing %s = %s\n", key, value)
	ds.data[key] = value
	time.Sleep(100 * time.Millisecond) // Simulasi penulisan
}

func (ds *DataStore) Read(key string) string {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	
	fmt.Printf("Reading %s\n", key)
	time.Sleep(50 * time.Millisecond) // Simulasi pembacaan
	return ds.data[key]
}

func main() {
	store := NewDataStore()
	
	// Goroutine writer
	go func() {
		store.Write("name", "Alice")
		store.Write("city", "New York")
	}()
	
	// Goroutine readers
	go func() {
		for i := 0; i < 3; i++ {
			name := store.Read("name")
			fmt.Printf("Reader 1 got name: %s\n", name)
		}
	}()
	
	go func() {
		for i := 0; i < 3; i++ {
			city := store.Read("city")
			fmt.Printf("Reader 2 got city: %s\n", city)
		}
	}()
	
	time.Sleep(2 * time.Second)
}
```

**Hasil (menunjukkan multiple readers bisa berjalan bersamaan):**
```
Writing name = Alice
Reading name
Reading name
Reader 1 got name: Alice
Reader 2 got name: Alice
Writing city = New York
Reading city
Reading city
Reader 1 got city: New York
Reader 2 got city: New York
```

## 4. Perilaku dan Gotchas

### a. Lupa Unlock
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex

func badFunction() {
	mutex.Lock()
	fmt.Println("Locked")
	// LUPA unlock - akan menyebabkan deadlock
	
	// Simulasi error
	if true {
		return // Fungsi keluar tanpa unlock
	}
	mutex.Unlock() // Tidak akan pernah dieksekusi
}

func main() {
	go badFunction()
	
	time.Sleep(100 * time.Millisecond)
	
	// Ini akan menyebabkan deadlock
	mutex.Lock()
	fmt.Println("This will never be printed")
	mutex.Unlock()
}
```

**Notifikasi Error:** JANGAN pernah lupa memanggil `Unlock()` setelah `Lock()`. Gunakan `defer` untuk memastikan unlock selalu dipanggil.

### b. Penggunaan Defer untuk Unlock
```go
package main

import (
	"fmt"
	"sync"
)

var (
	data  = make(map[string]int)
	mutex sync.Mutex
)

func safeOperation(key string, value int) {
	mutex.Lock()
	defer mutex.Unlock() // Memastikan unlock selalu dipanggil
	
	// Operasi kompleks yang bisa keluar di tengah jalan
	if key == "" {
		fmt.Println("Invalid key")
		return // Unlock tetap dipanggil karena defer
	}
	
	data[key] = value
	fmt.Printf("Set %s = %d\n", key, value)
}

func main() {
	safeOperation("", 10)   // Invalid key, tetap unlock
	safeOperation("key1", 20) // Valid operation
}
```

### c. Nested Lock (Deadlock Potensial)
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex

func functionA() {
	mutex.Lock()
	defer mutex.Unlock()
	
	fmt.Println("Function A acquired lock")
	time.Sleep(100 * time.Millisecond)
	
	// Memanggil fungsi lain yang juga membutuhkan lock
	// Bisa menyebabkan deadlock jika tidak hati-hati
	// functionB() // HATI-HATI DENGAN INI
}

func functionB() {
	mutex.Lock()
	defer mutex.Unlock()
	
	fmt.Println("Function B acquired lock")
}

func main() {
	go functionA()
	go functionB()
	
	time.Sleep(1 * time.Second)
}
```

**Notifikasi Error:** HINDARI nested lock pada mutex yang sama dalam call stack yang sama, karena bisa menyebabkan deadlock.

## 5. Best Practices

1.  **Gunakan defer untuk Unlock**: Ini memastikan mutex selalu dibuka, bahkan jika fungsi keluar secara prematur.
2.  **Minimalkan waktu lock**: Hanya kunci bagian kode yang benar-benar perlu proteksi.
3.  **Gunakan RWMutex jika sesuai**: Ketika Anda memiliki banyak pembaca dan sedikit penulis, RWMutex bisa memberikan performa yang lebih baik.
4.  **Hindari nested lock**: Jangan memanggil fungsi yang membutuhkan lock yang sama dalam fungsi yang sudah memiliki lock.
5.  **Pertimbangkan alternatif**: Untuk kasus sederhana, channel seringkali lebih idiomatik daripada mutex.

## 6. Pattern: Struktur dengan Mutex Tersemat

```go
package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	counter := &Counter{}
	
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	
	wg.Wait()
	fmt.Printf("Final counter value: %d\n", counter.Value())
}
```

## 7. Kesimpulan

Mutex adalah alat penting untuk sinkronisasi dalam pemrograman concurrent Go. Memahami cara kerja mutex sangat penting untuk:
- Mencegah race condition
- Memastikan data consistency
- Mengelola akses ke resource bersama

Dengan memahami mekanisme, perilaku, dan best practices dari mutex, Anda dapat menulis aplikasi concurrent yang aman dan efisien. Ingatlah bahwa meskipun mutex sangat berguna, dalam banyak kasus, channel bisa menjadi pilihan yang lebih idiomatik dalam Go untuk komunikasi antar goroutine.