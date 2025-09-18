# Bagian 1: Dasar-Dasar Golang

## 1. Instalasi dan Setup Lingkungan Kerja

### a. Mengunduh Go
1.  Kunjungi situs resmi [https://golang.org/dl/](https://golang.org/dl/).
2.  Unduh versi terbaru sesuai dengan sistem operasi Anda (Windows, macOS, Linux).
3.  Jalankan installer dan ikuti petunjuk instalasi.

### b. Verifikasi Instalasi
Buka terminal atau command prompt dan jalankan:
```bash
go version
```
Jika berhasil, akan menampilkan versi Go yang terinstal.

### c. Setup Workspace
Secara default, workspace Go berada di direktori `GOPATH`. Namun, sejak Go 1.11+, penggunaan modul lebih disarankan sehingga `GOPATH` tidak wajib lagi. Untuk saat ini, kita bisa menggunakan direktori proyek langsung.

Buat folder baru untuk proyek belajar kita:
```bash
mkdir learn-golang
cd learn-golang
```

### d. Text Editor / IDE
Gunakan editor teks favorit Anda. Beberapa pilihan populer yang memiliki dukungan baik untuk Go antara lain:
- Visual Studio Code dengan ekstensi Go
- Goland (IDE khusus Go dari JetBrains)
- Vim/Neovim dengan plugin coc-go atau vim-go

## 2. Struktur Dasar Program Golang

Program Go dimulai dari fungsi `main()` dalam package `main`.

Contoh program sederhana:
```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
```

Penjelasan:
- `package main`: Mendeklarasikan bahwa file ini adalah bagian dari package utama (`main`). Package `main` adalah titik masuk eksekusi program.
- `import "fmt"`: Mengimpor package `fmt` (format) yang digunakan untuk input/output standar.
- `func main()`: Fungsi utama tempat eksekusi program dimulai.
- `fmt.Println(...)`: Memanggil fungsi `Println` dari package `fmt` untuk mencetak teks ke layar diikuti baris baru.

Untuk menjalankan program:
1. Simpan kode di atas dengan nama `hello.go`.
2. Jalankan di terminal dengan perintah:
   ```bash
   go run hello.go
   ```

## 3. Variabel dan Tipe Data

### a. Deklarasi Variabel
Ada beberapa cara untuk mendeklarasikan variabel di Go:

1.  **Deklarasi panjang:**
    ```go
    var nama string = "John Doe"
    var umur int = 30
    ```
2.  **Deklarasi dengan type inference (inferensi tipe):**
    Go bisa secara otomatis menentukan tipe data berdasarkan nilai yang diberikan.
    ```go
    var nama = "John Doe" // Tipe string ditentukan otomatis
    var umur = 30         // Tipe int ditentukan otomatis
    ```
3.  **Short variable declaration (`:=`) - hanya dalam fungsi:**
    Cara ini paling umum digunakan dalam fungsi.
    ```go
    func main() {
        nama := "John Doe"
        umur := 30
    }
    ```
    Catatan: `:=` digunakan untuk deklarasi dan inisialisasi sekaligus, sedangkan `=` digunakan untuk assignment (penugasan) ke variabel yang sudah dideklarasikan.

### b. Tipe Data Dasar
- **Numbers:**
  - `int`: Bilangan bulat (signed) ukuran tergantung platform (32 atau 64 bit).
  - `int8`, `int16`, `int32`, `int64`: Bilangan bulat signed dengan ukuran spesifik.
  - `uint8`, `uint16`, `uint32`, `uint64`: Bilangan bulat unsigned dengan ukuran spesifik.
  - `float32`, `float64`: Bilangan desimal.
  - `complex64`, `complex128`: Bilangan kompleks.
- **String:** Deretan karakter, dinyatakan dengan tanda kutip ganda (`"..."`) atau backticks (`` `...` ``) untuk raw string literal.
- **Boolean:** `bool`, bernilai `true` atau `false`.

Contoh:
```go
func main() {
	var isReady bool = true
	var count int = 10
	var price float64 = 99.99
	var message string = "Welcome to Go!"
	
	// Short declaration
	isLoggedIn := false
	userName := "Alice"
	
	fmt.Println(isReady, count, price, message)
	fmt.Println(isLoggedIn, userName)
}
```

## 4. Operator

### a. Operator Aritmatika
- `+` (Penjumlahan)
- `-` (Pengurangan)
- `*` (Perkalian)
- `/` (Pembagian)
- `%` (Modulo/Sisa bagi)

### b. Operator Perbandingan
- `==` (Sama dengan)
- `!=` (Tidak sama dengan)
- `<` (Kurang dari)
- `<=` (Kurang dari atau sama dengan)
- `>` (Lebih dari)
- `>=` (Lebih dari atau sama dengan)

### c. Operator Logika
- `&&` (AND)
- `||` (OR)
- `!` (NOT)

### d. Operator Lainnya
- `=` (Assignment)
- `:=` (Short variable declaration)
- `+=`, `-=`, `*=`, `/=`, `%=` (Compound assignment)

Contoh:
```go
func main() {
	a := 10
	b := 5
	
	// Aritmatika
	sum := a + b
	diff := a - b
	product := a * b
	quotient := a / b
	remainder := a % b
	
	fmt.Println("Aritmatika:", sum, diff, product, quotient, remainder)
	
	// Perbandingan
	isEqual := a == b
	isGreater := a > b
	
	fmt.Println("Perbandingan:", isEqual, isGreater)
	
	// Logika
	x := true
	y := false
	andResult := x && y
	orResult := x || y
	notResult := !x
	
	fmt.Println("Logika:", andResult, orResult, notResult)
}
```

## 5. Percabangan (Branching)

### a. If/Else Statement
Digunakan untuk mengeksekusi blok kode jika suatu kondisi terpenuhi.

Sintaks dasar:
```go
if kondisi {
	// kode jika kondisi benar
} else if kondisi_lain {
	// kode jika kondisi_lain benar
} else {
	// kode jika semua kondisi salah
}
```

Contoh:
```go
func main() {
	age := 20
	
	if age >= 18 {
		fmt.Println("Dewasa")
	} else {
		fmt.Println("Anak-anak")
	}
	
	score := 85
	if score >= 90 {
		fmt.Println("Grade: A")
	} else if score >= 80 {
		fmt.Println("Grade: B")
	} else if score >= 70 {
		fmt.Println("Grade: C")
	} else {
		fmt.Println("Grade: D")
	}
}
```

### b. Switch Statement
Alternatif dari if/else yang lebih rapi untuk banyak kondisi.

Sintaks dasar:
```go
switch ekspresi {
case nilai1:
	// kode jika ekspresi == nilai1
case nilai2:
	// kode jika ekspresi == nilai2
default:
	// kode jika tidak ada yang cocok
}
```

Contoh:
```go
func main() {
	day := "Monday"
	
	switch day {
	case "Monday":
		fmt.Println("It's Monday!")
	case "Tuesday":
		fmt.Println("It's Tuesday!")
	default:
		fmt.Println("It's another day.")
	}
	
	// Switch tanpa ekspresi (mirip dengan if/else)
	num := 7
	switch {
	case num < 0:
		fmt.Println("Negative number")
	case num > 0:
		fmt.Println("Positive number")
	default:
		fmt.Println("Zero")
	}
}
```

## 6. Perulangan (Looping)

### a. For Loop
Go hanya memiliki satu jenis perulangan yaitu `for`.

Sintaks dasar:
```go
for inisialisasi; kondisi; iterasi {
	// kode yang diulang
}
```

Contoh:
```go
func main() {
	// Perulangan standar
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}
	
	// Perulangan while-style (hanya kondisi)
	j := 0
	for j < 3 {
		fmt.Println(j)
		j++
	}
	
	// Perulangan forever (infinite loop)
	// Hati-hati, gunakan break untuk keluar
	k := 0
	for {
		if k >= 2 {
			break // Keluar dari loop
		}
		fmt.Println("Infinite loop iteration", k)
		k++
	}
}
```

### b. Range Loop
Digunakan untuk mengiterasi elemen-elemen dalam array, slice, map, atau channel.

Contoh dengan slice:
```go
func main() {
	numbers := []int{1, 2, 3, 4, 5}
	
	for index, value := range numbers {
		fmt.Printf("Index: %d, Value: %d\n", index, value)
	}
	
	// Jika hanya butuh value, gunakan _
	for _, value := range numbers {
		fmt.Println("Value:", value)
	}
	
	// Jika hanya butuh index
	for index := range numbers {
		fmt.Println("Index:", index)
	}
}
```

## 7. Fungsi (Functions)

Fungsi adalah blok kode yang dapat dipanggil untuk melakukan tugas tertentu. Fungsi membantu dalam modularisasi kode.

### a. Deklarasi Fungsi
Sintaks dasar:
```go
func namaFungsi(parameter1 tipe1, parameter2 tipe2) tipeReturn {
	// kode fungsi
	return nilai
}
```

Contoh:
```go
// Fungsi tanpa parameter dan return
func greet() {
	fmt.Println("Hello!")
}

// Fungsi dengan parameter
func greetPerson(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// Fungsi dengan return value
func add(a int, b int) int {
	return a + b
}

// Fungsi dengan multiple return values
func divide(dividend, divisor int) (int, error) {
	if divisor == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return dividend / divisor, nil
}

func main() {
	greet()
	greetPerson("Bob")
	result := add(5, 3)
	fmt.Println("Addition result:", result)
	
	quotient, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Division result:", quotient)
	}
	
	_, zeroErr := divide(10, 0) // Mengabaikan hasil pertama
	if zeroErr != nil {
		fmt.Println("Expected error:", zeroErr)
	}
}
```

### b. Variadic Functions
Fungsi yang dapat menerima jumlah argumen yang bervariasi.

Contoh:
```go
func sum(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

func main() {
	fmt.Println(sum(1, 2, 3))       // Output: 6
	fmt.Println(sum(4, 5, 6, 7, 8)) // Output: 30
}
```

### c. Anonymous Functions dan Closures
Anonymous function adalah fungsi tanpa nama. Closure adalah fungsi yang mereferensikan variabel dari lingkup luar.

Contoh:
```go
func main() {
	// Anonymous function
	func() {
		fmt.Println("This is an anonymous function")
	}()
	
	// Assign anonymous function ke variabel
	greeting := func(name string) string {
		return "Hello, " + name
	}
	fmt.Println(greeting("Charlie"))
	
	// Closure
	makeCounter := func() func() int {
		count := 0
		return func() int {
			count++
			return count
		}
	}
	
	counter := makeCounter()
	fmt.Println(counter()) // Output: 1
	fmt.Println(counter()) // Output: 2
}
```