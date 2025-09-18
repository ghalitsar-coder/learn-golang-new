# Bagian 4: Package, Module, dan Testing

## 1. Package

Package adalah cara untuk mengelompokkan dan mengorganisir kode Go. Setiap file Go dimulai dengan deklarasi `package`.

### a. Membuat Package
Buat direktori baru untuk package Anda, misalnya `calculator`.

`calculator/calculator.go`:
```go
package calculator

// Add menjumlahkan dua angka
func Add(a, b int) int {
	return a + b
}

// Subtract mengurangkan dua angka
func Subtract(a, b int) int {
	return a - b
}
```

### b. Menggunakan Package
Untuk menggunakan package, kita perlu mengimport-nya. Struktur direktori proyek menjadi penting.

Struktur direktori:
```
myproject/
├── main.go
└── calculator/
    └── calculator.go
```

`main.go`:
```go
package main

import (
	"fmt"
	"myproject/calculator" // Path relatif terhadap module root
)

func main() {
	result := calculator.Add(5, 3)
	fmt.Println("Addition result:", result)
}
```

### c. Exported dan Unexported Identifiers
- Identifier (nama fungsi, variabel, struct, dll.) yang diawali huruf besar adalah *exported* (publik) dan bisa diakses dari package lain.
- Identifier yang diawali huruf kecil adalah *unexported* (private) dan hanya bisa diakses dalam package yang sama.

Contoh:
```go
// Dalam package mathutils
package mathutils

// Exported function
func Average(numbers []float64) float64 {
	// ...
}

// Unexported function
func sum(numbers []float64) float64 {
	// ...
	total := 0.0
	for _, num := range numbers {
		total += num
	}
	return total
}

func Average(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	return sum(numbers) / float64(len(numbers)) // Memanggil fungsi unexported
}
```

## 2. Go Modules

Go Modules adalah sistem manajemen dependensi resmi untuk Go. Ini memungkinkan kita untuk membuat proyek independen dengan versi dependensi yang spesifik.

### a. Inisialisasi Module
Di direktori proyek Anda, jalankan:
```bash
go mod init nama_module
```
Misalnya:
```bash
go mod init myproject
```

Ini akan membuat file `go.mod`:
```go
module myproject

go 1.21 // Versi Go yang digunakan
```

### b. Menambah Dependensi
Saat Anda mengimport package eksternal dan menjalankan `go run` atau `go build`, Go akan secara otomatis menambahkannya ke `go.mod`.

Contoh: Menambahkan package `github.com/google/uuid`:
```go
// main.go
package main

import (
	"fmt"
	"github.com/google/uuid"
)

func main() {
	id := uuid.New()
	fmt.Println("Generated UUID:", id)
}
```

Jalankan:
```bash
go mod tidy
```
Perintah ini akan membersihkan dan memperbarui `go.mod` dan `go.sum` dengan dependensi yang dibutuhkan.

### c. File `go.mod` dan `go.sum`
- `go.mod`: Mendefinisikan module path dan dependensi.
- `go.sum`: Berisi checksum kriptografis dari versi spesifik setiap dependensi untuk memastikan reproducibility.

Contoh `go.mod` setelah menambah dependensi:
```go
module myproject

go 1.21

require github.com/google/uuid v1.3.0
```

## 3. Unit Testing

Go memiliki sistem testing bawaan yang sangat kuat dan mudah digunakan.

### a. Membuat File Test
File test harus memiliki akhiran `_test.go`. Misalnya, untuk `calculator.go`, buat `calculator_test.go`.

`calculator/calculator.go`:
```go
package calculator

// Add menjumlahkan dua angka
func Add(a, b int) int {
	return a + b
}

// Subtract mengurangkan dua angka
func Subtract(a, b int) int {
	return a - b
}
```

`calculator/calculator_test.go`:
```go
package calculator

import (
	"testing"
)

// Fungsi test dimulai dengan Test dan mengambil *testing.T
func TestAdd(t *testing.T) {
	// Test case 1
	result := Add(2, 3)
	expected := 5
	if result != expected {
		// Gunakan t.Errorf untuk melaporkan error
		t.Errorf("Add(2, 3) = %d; expected %d", result, expected)
	}
	
	// Test case 2
	result = Add(-1, 1)
	expected = 0
	if result != expected {
		t.Errorf("Add(-1, 1) = %d; expected %d", result, expected)
	}
}

func TestSubtract(t *testing.T) {
	result := Subtract(5, 3)
	expected := 2
	if result != expected {
		t.Errorf("Subtract(5, 3) = %d; expected %d", result, expected)
	}
}
```

### b. Menjalankan Test
Di terminal, jalankan:
```bash
go test ./... # Menjalankan semua test dalam module
# atau
go test ./calculator # Menjalankan test dalam package calculator
# atau
go test -v # Menjalankan dengan output verbose
```

### c. Table Driven Tests
Cara idiomatik untuk menulis beberapa test case dalam satu fungsi.

`calculator_test.go` (versi diperbaiki):
```go
package calculator

import (
	"testing"
)

func TestAdd(t *testing.T) {
	// Mendefinisikan test cases
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"with zero", 5, 0, 5},
		{"negative and positive", -1, 1, 0},
		{"negative numbers", -2, -3, -5},
	}
	
	// Iterasi test cases
	for _, tt := range tests {
		// t.Run membuat subtest
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"positive numbers", 5, 3, 2},
		{"result negative", 3, 5, -2},
		{"with zero", 5, 0, 5},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subtract(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Subtract(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}
```

### d. Test Helpers dan Setup/Teardown
Untuk fungsi bantu dalam test:
```go
// helper function (unexported)
func assertEqual(t *testing.T, got, want int, operation string) {
	t.Helper() // Menandai fungsi ini sebagai test helper
	if got != want {
		t.Errorf("%s = %d, want %d", operation, got, want)
	}
}

// Menggunakan helper dalam test
func TestAddWithHelper(t *testing.T) {
	result := Add(2, 3)
	assertEqual(t, result, 5, "Add(2, 3)")
}
```

### e. Coverage
Go menyediakan alat untuk mengukur cakupan kode test:
```bash
go test -cover ./calculator
# atau untuk mendapatkan laporan detail
go test -coverprofile=coverage.out ./calculator
go tool cover -html=coverage.out # Membuka laporan HTML
```

## 4. Benchmarking

Benchmark adalah test khusus untuk mengukur performa kode.

### a. Membuat Benchmark
Fungsi benchmark dimulai dengan `Benchmark` dan mengambil `*testing.B`.

`calculator/calculator_bench_test.go`:
```go
package calculator

import (
	"testing"
)

// BenchmarkAdd mengukur performa fungsi Add
func BenchmarkAdd(b *testing.B) {
	// Reset timer untuk mengabaikan setup
	b.ResetTimer()
	
	// Loop b.N kali
	for i := 0; i < b.N; i++ {
		Add(100, 200)
	}
}

// BenchmarkAddLarge mengukur performa dengan data besar
func BenchmarkAddLarge(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := 0
		for _, n := range numbers {
			result = Add(result, n)
		}
	}
}
```

### b. Menjalankan Benchmark
```bash
go test -bench=. ./calculator
# atau
go test -bench=BenchmarkAdd ./calculator
# dengan alokasi memori
go test -bench=. -benchmem ./calculator
```

Output benchmark akan menunjukkan:
- Jumlah iterasi (`b.N`)
- Waktu rata-rata per operasi (`ns/op`)
- Alokasi memori per operasi (`B/op`)
- Jumlah alokasi per operasi (`allocs/op`)

Contoh output:
```
goos: linux
goarch: amd64
pkg: myproject/calculator
cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
BenchmarkAdd-8            	1000000000	         0.2744 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddLarge-8       	 1000000	      1092 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	myproject/calculator	2.345s
```

Dengan package, module, dan testing yang solid, kode Go menjadi lebih modular, dapat dikelola, dan dapat diandalkan.