# Konsep-Konsep Kunci dalam Golang

Materi ini fokus pada konsep-konsep fundamental dan sering digunakan dalam pemrograman Go yang sering menjadi penyebab kebingungan atau error bagi pemula. Kita akan membahas perilaku masing-masing konsep melalui percobaan dan contoh error.

## 1. Array vs Slice

### Konsep
- **Array**: Struktur data dengan ukuran tetap. Tipe array ditentukan oleh panjang dan tipe elemen (`[5]int` ≠ `[10]int`).
- **Slice**: Struktur data dinamis yang merupakan referensi ke bagian dari array. Tipe slice hanya ditentukan oleh tipe elemen (`[]int`).

### Perilaku dan Percobaan

**Percobaan A: Pass by Value vs Reference**
```go
package main

import "fmt"

func modifyArray(arr [3]int) {
	arr[0] = 100
	fmt.Println("Inside modifyArray:", arr)
}

func modifySlice(slc []int) {
	slc[0] = 100
	fmt.Println("Inside modifySlice:", slc)
}

func main() {
	// Array
	myArr := [3]int{1, 2, 3}
	fmt.Println("Before modifyArray:", myArr)
	modifyArray(myArr)
	fmt.Println("After modifyArray:", myArr) // Apa yang terjadi?

	// Slice
	mySlice := []int{1, 2, 3}
	fmt.Println("Before modifySlice:", mySlice)
	modifySlice(mySlice)
	fmt.Println("After modifySlice:", mySlice) // Apa yang terjadi?
}
```
**Hasil Percobaan A:**
```
Before modifyArray: [1 2 3]
Inside modifyArray: [100 2 3]
After modifyArray: [1 2 3]  // Array tidak berubah karena pass by value
Before modifySlice: [1 2 3]
Inside modifySlice: [100 2 3]
After modifySlice: [100 2 3] // Slice berubah karena pass by reference
```

**Percobaan B: Menambah elemen**
```go
package main

import "fmt"

func main() {
	// Array - tidak bisa diubah ukurannya
	var arr [3]int
	// arr = append(arr, 4) // ERROR: first argument to append must be slice; have [3]int

	// Slice - bisa diubah ukurannya
	slc := []int{1, 2, 3}
	slc = append(slc, 4)
	fmt.Println("Slice after append:", slc) // [1 2 3 4]
}
```

**Error yang Sering Terjadi:**
```go
// JANGAN LAKUKAN INI
func processItems(items [100]int) {
	// Fungsi ini hanya menerima array dengan panjang 100
	// Sangat kaku dan tidak fleksibel
}

func main() {
	smallArray := [50]int{} // Panjang berbeda
	// processItems(smallArray) // ERROR: cannot use smallArray (type [50]int) as type [100]int
}
```
**Notifikasi Error:** Hindari penggunaan array dengan panjang tetap kecuali benar-benar diperlukan. Gunakan slice untuk fleksibilitas yang lebih baik.

---

## 2. Method Receiver: Value vs Pointer

### Konsep
Method bisa memiliki receiver berupa nilai (`T`) atau pointer (`*T`). Pemilihan receiver mempengaruhi apakah method bisa memodifikasi data dan efisiensi memori.

### Perilaku dan Percobaan

**Percobaan A: Modifikasi Data**
```go
package main

import "fmt"

type Counter struct {
	Value int
}

// Receiver nilai
func (c Counter) IncrementByValue() {
	c.Value++
	fmt.Println("Inside IncrementByValue:", c.Value)
}

// Receiver pointer
func (c *Counter) IncrementByPointer() {
	c.Value++
	fmt.Println("Inside IncrementByPointer:", c.Value)
}

func main() {
	// Dengan receiver nilai
	c1 := Counter{Value: 0}
	fmt.Println("Before IncrementByValue:", c1.Value)
	c1.IncrementByValue()
	fmt.Println("After IncrementByValue:", c1.Value) // Apa yang terjadi?

	// Dengan receiver pointer
	c2 := Counter{Value: 0}
	fmt.Println("Before IncrementByPointer:", c2.Value)
	c2.IncrementByPointer()
	fmt.Println("After IncrementByPointer:", c2.Value) // Apa yang terjadi?
}
```
**Hasil Percobaan A:**
```
Before IncrementByValue: 0
Inside IncrementByValue: 1
After IncrementByValue: 0  // Tidak berubah karena bekerja pada copy
Before IncrementByPointer: 0
Inside IncrementByPointer: 1
After IncrementByPointer: 1 // Berubah karena bekerja pada alamat memori yang sama
```

**Percobaan B: Efisiensi Memori**
```go
package main

import (
	"fmt"
)

type BigStruct struct {
	Data [1000000]int // Struktur besar
}

// Receiver nilai - menyalin seluruh struct
func (bs BigStruct) GetValue() int {
	return bs.Data[0]
}

// Receiver pointer - hanya menyalin pointer (8 byte)
func (bs *BigStruct) SetValue(val int) {
	bs.Data[0] = val
}

func main() {
	bs := BigStruct{}
	
	// Memanggil method dengan receiver nilai
	// Ini akan menyalin seluruh struct (tidak efisien)
	val := bs.GetValue() 
	fmt.Println("Value:", val)
	
	// Memanggil method dengan receiver pointer
	// Ini hanya menyalin pointer (efisien)
	bs.SetValue(42)
	fmt.Println("New value:", bs.Data[0])
}
```

**Error yang Sering Terjadi:**
```go
// JANGAN LAKUKAN INI
type MyType struct {
	Value int
}

// Campuran receiver bisa menyebabkan kebingungan
func (m MyType) GetValue() int {
	return m.Value
}

func (m *MyType) SetValue(v int) {
	m.Value = v
}

func main() {
	var m MyType
	// m.GetValue() // OK
	// m.SetValue(10) // ERROR: cannot call pointer method on m
	
	// Harus menggunakan pointer
	p := &m
	p.SetValue(10) // OK
}
```
**Notifikasi Error:** Konsistenlah dalam menggunakan receiver (semua nilai atau semua pointer) dalam satu tipe untuk menghindari kebingungan dan error.

---

## 3. Error Handling

### Konsep
Go menggunakan tipe `error` sebagai return value untuk menunjukkan kesalahan. Ini berbeda dari exception di bahasa lain.

### Perilaku dan Percobaan

**Percobaan A: Lupa Mengecek Error**
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	// Sering terjadi: lupa mengecek error
	file, err := os.Open("nonexistent.txt")
	// if err != nil { // LUPA MENGECEK ERROR
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	
	// Program akan panic di sini karena file adalah nil
	defer file.Close()
	fmt.Println("File opened successfully")
}
```
**Hasil Percobaan A:**
```
panic: runtime error: invalid memory address or nil pointer dereference
```

**Percobaan B: Mengabaikan Error**
```go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Mengabaikan error saat konversi
	str := "abc"
	num, _ := strconv.Atoi(str) // Mengabaikan error dengan _
	fmt.Println("Number:", num) // Apa yang terjadi?
}
```
**Hasil Percobaan B:**
```
Number: 0  // Atoi mengembalikan 0 untuk input yang tidak valid
```

**Error yang Sering Terjadi:**
```go
// JANGAN LAKUKAN INI
func divide(a, b float64) float64 {
	if b == 0 {
		// Tidak bagus: panic untuk error yang bisa diprediksi
		panic("division by zero") 
	}
	return a / b
}

func main() {
	// Lebih baik menggunakan error return
	result := divide(10, 0) // Program akan crash
	fmt.Println(result)
}
```
**Notifikasi Error:** Gunakan tipe `error` sebagai return value untuk error yang bisa diprediksi. Gunakan `panic` hanya untuk kesalahan pemrograman yang serius.

---

## 4. Goroutine dan Channel

### Konsep
Goroutine adalah lightweight thread. Channel digunakan untuk komunikasi antar goroutine secara aman.

### Perilaku dan Percobaan

**Percobaan A: Goroutine Tanpa Sinkronisasi**
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
**Hasil Percobaan A:**
```
Main finished
// Tidak ada output dari goroutine
```

**Percobaan B: Deadlock dengan Channel**
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
**Hasil Percobaan B:**
```
fatal error: all goroutines are asleep - deadlock!
```

**Error yang Sering Terjadi:**
```go
// JANGAN LAKUKAN INI - Channel Ditutup Sebelum Selesai Digunakan
func main() {
	ch := make(chan int, 2)
	
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i // Bisa menyebabkan panic jika channel sudah ditutup
		}
	}()
	
	go func() {
		for i := 0; i < 3; i++ {
			fmt.Println(<-ch)
		}
		close(ch) // Menutup channel terlalu dini
	}()
	
	time.Sleep(1 * time.Second)
}
```
**Notifikasi Error:** Pastikan goroutine memiliki cara untuk menyinkronkan dengan program utama (misalnya dengan channel, `sync.WaitGroup`). Hati-hati saat menutup channel yang masih digunakan oleh goroutine lain.

---

## 5. Interface dan Type Assertion

### Konsep
Interface mendefinisikan perilaku (method). Type assertion digunakan untuk mengambil tipe konkret dari interface.

### Perilaku dan Percobaan

**Percobaan A: Type Assertion yang Gagal**
```go
package main

import "fmt"

func main() {
	var i interface{} = "hello"
	
	// Type assertion yang berhasil
	s := i.(string)
	fmt.Println(s)
	
	// Type assertion yang gagal
	// n := i.(int) // Panic: interface conversion: interface {} is string, not int
	
	// Cara aman menggunakan comma ok idiom
	if n, ok := i.(int); ok {
		fmt.Println("Number is", n)
	} else {
		fmt.Println("Value is not an int")
	}
}
```
**Hasil Percobaan A:**
```
hello
Value is not an int
```

**Percobaan B: Interface Kosong**
```go
package main

import "fmt"

func printType(v interface{}) {
	// Menggunakan type switch
	switch val := v.(type) {
	case int:
		fmt.Printf("Integer: %d\n", val)
	case string:
		fmt.Printf("String: %s\n", val)
	case bool:
		fmt.Printf("Boolean: %t\n", val)
	default:
		fmt.Printf("Unknown type: %T\n", val)
	}
}

func main() {
	printType(42)
	printType("Go")
	printType(true)
	printType(3.14)
}
```
**Hasil Percobaan B:**
```
Integer: 42
String: Go
Boolean: true
Unknown type: float64
```

**Error yang Sering Terjadi:**
```go
// JANGAN LAKUKAN INI - Menggunakan interface{} secara berlebihan
func processAnything(data interface{}) {
	// Kehilangan type safety
	// Sulit untuk maintain dan prone error
	// Gunakan generics (Go 1.18+) jika memungkinkan
}

// Lebih baik:
// func processData[T any](data T) {
//     // Type safe
// }
```
**Notifikasi Error:** Gunakan type assertion dengan aman menggunakan comma ok idiom. Pertimbangkan penggunaan generics daripada `interface{}` untuk type safety yang lebih baik.

---

## 6. Defer, Panic, dan Recover

### Konsep
`defer` menunda eksekusi. `panic` menghentikan alur normal. `recover` menangkap `panic`.

### Perilaku dan Percobaan

**Percobaan A: Urutan Eksekusi Defer**
```go
package main

import "fmt"

func main() {
	defer fmt.Println("Defer 1")
	defer fmt.Println("Defer 2")
	defer fmt.Println("Defer 3")
	
	fmt.Println("Main function")
	// Defer dieksekusi dalam urutan LIFO (Last In, First Out)
}
```
**Hasil Percobaan A:**
```
Main function
Defer 3
Defer 2
Defer 1
```

**Percobaan B: Recover dari Panic**
```go
package main

import "fmt"

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
**Hasil Percobaan B:**
```
About to panic
Recovered from panic: This is a panic!
Program continues after panic recovery
```

**Error yang Sering Terjadi:**
```go
// JANGAN LAKUKAN INI - Menggunakan panic untuk error normal
func divide(a, b int) int {
	if b == 0 {
		panic("division by zero") // Tidak sesuai
	}
	return a / b
}

// Lebih baik:
func divideSafe(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}
```
**Notifikasi Error:** Gunakan `panic` hanya untuk kesalahan serius yang tidak bisa atau tidak seharusnya ditangani. Untuk error normal, gunakan tipe `error` sebagai return value.

---

Dengan memahami konsep-konsep ini dan perilakunya melalui percobaan, Anda akan lebih siap untuk menulis kode Go yang efektif dan menghindari error umum.