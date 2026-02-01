# Pass by Value vs Pass by Reference di Go

## 📌 Konsep Dasar

### **Pass by Value**
Ketika fungsi menerima parameter, Go **menyalin nilai** dari argument ke parameter. Perubahan pada parameter di dalam fungsi **tidak mempengaruhi** variable asli.

### **Pass by Reference**
Ketika fungsi menerima **pointer/reference**, fungsi dapat mengakses dan memodifikasi data asli di memory. Perubahan di dalam fungsi **mempengaruhi** variable asli.

---

## 🔍 Prinsip Fundamental Go

> **Go SELALU pass by value!**

Tapi... ada twist-nya:
- **Nilai primitif**: Copy nilai langsung
- **Pointer**: Copy alamat memory (referensi ke data asli)
- **Slice, Map, Channel**: Copy **header** yang berisi pointer ke underlying data

---

## 1️⃣ Pass by Value (Primitif)

### ✅ Tipe Data yang Pass by Value
- `int`, `float64`, `bool`, `string`
- `struct` (keseluruhan struct di-copy)
- `array` (keseluruhan array di-copy)

### Contoh:
```go
package main

import "fmt"

func ubahNilai(angka int) {
    angka = 100  // Hanya mengubah copy
    fmt.Println("Dalam fungsi:", angka) // 100
}

func main() {
    x := 50
    ubahNilai(x)
    fmt.Println("Di main:", x)  // Tetap 50 (tidak berubah)
}
```

**Output:**
```
Dalam fungsi: 100
Di main: 50
```

### Cara Kerja:
```
Memory:
main() scope:
  x = 50           [alamat: 0x1000]
  
ubahNilai() scope:
  angka = 50 (copy) [alamat: 0x2000]  ← Copy terpisah!
  angka = 100       [alamat: 0x2000]  ← Ubah copy saja
  
x tetap 50! ✅
```

---

## 2️⃣ Pass by Reference (Pointer)

### Contoh dengan Pointer:
```go
package main

import "fmt"

func ubahNilaiPointer(angka *int) {
    *angka = 100  // Mengubah nilai di alamat memory asli
    fmt.Println("Dalam fungsi:", *angka) // 100
}

func main() {
    x := 50
    ubahNilaiPointer(&x)  // Kirim alamat memory
    fmt.Println("Di main:", x)  // 100 (berubah!)
}
```

**Output:**
```
Dalam fungsi: 100
Di main: 100
```

### Cara Kerja:
```
Memory:
main() scope:
  x = 50              [alamat: 0x1000]
  &x = 0x1000         ← Alamat dari x
  
ubahNilaiPointer() scope:
  angka = 0x1000 (copy alamat) [alamat parameter: 0x2000]
  *angka → akses 0x1000
  *angka = 100        ← Ubah nilai di 0x1000
  
x berubah jadi 100! ✅
```

---

## 3️⃣ Struct: Value vs Pointer

### Pass by Value (Struct Copy):
```go
type Person struct {
    Name string
    Age  int
}

func updateAge(p Person) {
    p.Age = 30  // Copy di-update
    fmt.Println("Dalam fungsi:", p.Age)
}

func main() {
    person := Person{Name: "Alice", Age: 25}
    updateAge(person)
    fmt.Println("Di main:", person.Age)  // Tetap 25
}
```

**Output:**
```
Dalam fungsi: 30
Di main: 25
```

### Pass by Reference (Struct Pointer):
```go
func updateAgePointer(p *Person) {
    p.Age = 30  // Go otomatis dereference
    // Sama dengan: (*p).Age = 30
    fmt.Println("Dalam fungsi:", p.Age)
}

func main() {
    person := Person{Name: "Alice", Age: 25}
    updateAgePointer(&person)
    fmt.Println("Di main:", person.Age)  // Berubah jadi 30
}
```

**Output:**
```
Dalam fungsi: 30
Di main: 30
```

---

## 4️⃣ Array vs Slice (Penting!)

### Array: Full Copy
```go
func updateArray(arr [3]int) {
    arr[0] = 999
    fmt.Println("Dalam fungsi:", arr)
}

func main() {
    numbers := [3]int{1, 2, 3}
    updateArray(numbers)
    fmt.Println("Di main:", numbers)  // Tetap [1 2 3]
}
```

### Slice: Header Copy (Data Shared!)
```go
func updateSlice(s []int) {
    s[0] = 999  // Mengubah underlying array
    fmt.Println("Dalam fungsi:", s)
}

func main() {
    numbers := []int{1, 2, 3}
    updateSlice(numbers)
    fmt.Println("Di main:", numbers)  // Berubah jadi [999 2 3]
}
```

**Kenapa?**
```
Slice adalah struct:
type slice struct {
    ptr   *int  // Pointer ke array
    len   int
    cap   int
}

Yang di-copy adalah struct-nya, tapi pointer tetap sama!
```

### Gotcha: Append di Fungsi
```go
func appendSlice(s []int) {
    s = append(s, 4)  // Buat slice baru (realloc)
    fmt.Println("Dalam fungsi:", s)
}

func main() {
    numbers := []int{1, 2, 3}
    appendSlice(numbers)
    fmt.Println("Di main:", numbers)  // Tetap [1 2 3]
}
```

**Solusi: Return atau Pointer**
```go
func appendSliceReturn(s []int) []int {
    return append(s, 4)
}

func appendSlicePointer(s *[]int) {
    *s = append(*s, 4)
}

func main() {
    numbers := []int{1, 2, 3}
    
    // Opsi 1: Return
    numbers = appendSliceReturn(numbers)
    
    // Opsi 2: Pointer ke slice
    appendSlicePointer(&numbers)
    
    fmt.Println(numbers)  // [1 2 3 4 4]
}
```

---

## 5️⃣ Map dan Channel (Reference Types)

### Map: Selalu Reference Behavior
```go
func updateMap(m map[string]int) {
    m["key"] = 100  // Mengubah map asli
}

func main() {
    myMap := map[string]int{"key": 50}
    updateMap(myMap)
    fmt.Println(myMap["key"])  // 100 (berubah!)
}
```

**Alasan:**
Map adalah pointer ke hash table. Copy pointer tetap nunjuk ke data sama.

### Channel: Juga Reference
```go
func sendData(ch chan int) {
    ch <- 42  // Kirim ke channel asli
}

func main() {
    ch := make(chan int)
    go sendData(ch)
    fmt.Println(<-ch)  // 42
}
```

---

## 📊 Tabel Ringkasan

| Tipe Data      | Pass by   | Modifikasi Mempengaruhi Asli? | Catatan                        |
|----------------|-----------|--------------------------------|--------------------------------|
| `int`, `float` | Value     | ❌ Tidak                       | Copy penuh                     |
| `string`       | Value     | ❌ Tidak                       | Immutable                      |
| `bool`         | Value     | ❌ Tidak                       | Copy penuh                     |
| `array`        | Value     | ❌ Tidak                       | Copy seluruh array             |
| `struct`       | Value     | ❌ Tidak                       | Copy seluruh field             |
| `*pointer`     | Value     | ✅ Ya (via dereference)        | Copy alamat, akses data sama   |
| `slice`        | Value     | ✅ Ya (modify elements)        | Copy header, data shared       |
|                |           | ❌ Tidak (append/reassign)     | Kecuali pakai pointer ke slice |
| `map`          | Value     | ✅ Ya                          | Copy pointer ke hash table     |
| `channel`      | Value     | ✅ Ya                          | Copy pointer ke channel        |

---

## 🎯 Kapan Menggunakan Apa?

### Gunakan **Pass by Value** ketika:
- ✅ Data kecil (int, bool, small struct)
- ✅ Tidak perlu modifikasi data asli
- ✅ Ingin immutability
- ✅ Fungsi pure (no side effects)

### Gunakan **Pass by Pointer** ketika:
- ✅ Data besar (large struct) - hindari copy overhead
- ✅ Perlu modifikasi data asli
- ✅ Struct dengan method yang modify state
- ✅ Perlu `nil` sebagai value (pointer bisa nil)

### Contoh Best Practice:
```go
// Small struct: by value
type Point struct {
    X, Y int
}

func (p Point) Distance() float64 {
    return math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
}

// Large struct: by pointer
type Config struct {
    Database   DatabaseConfig
    Cache      CacheConfig
    Logger     LoggerConfig
    // ... banyak field lain
}

func (c *Config) Validate() error {
    // Validasi tanpa copy overhead
    return nil
}

func (c *Config) Update(key, value string) {
    // Modifikasi langsung
}
```

---

## ⚠️ Gotchas & Common Mistakes

### 1. Nil Pointer Dereference
```go
func updatePerson(p *Person) {
    p.Name = "Bob"  // PANIC jika p == nil!
}

// Solusi: Check nil
func updatePersonSafe(p *Person) {
    if p == nil {
        return
    }
    p.Name = "Bob"
}
```

### 2. Loop Variable Pointer
```go
// ❌ SALAH
var pointers []*int
numbers := []int{1, 2, 3}
for _, num := range numbers {
    pointers = append(pointers, &num)  // Semua pointer sama!
}
// Semua pointer nunjuk ke variable terakhir

// ✅ BENAR
for _, num := range numbers {
    n := num  // Copy ke variable baru
    pointers = append(pointers, &n)
}

// Atau (Go 1.22+):
for _, num := range numbers {
    pointers = append(pointers, &num)  // Fixed di Go 1.22
}
```

### 3. Slice Append di Goroutine
```go
// ❌ Race condition
func appendConcurrent(s []int, val int) {
    s = append(s, val)  // Tidak aman!
}

// ✅ Gunakan mutex atau channel
var mu sync.Mutex
func appendSafe(s *[]int, val int) {
    mu.Lock()
    *s = append(*s, val)
    mu.Unlock()
}
```

---

## 🧪 Eksperimen: Ukur Memory Address

```go
package main

import "fmt"

func main() {
    // Primitif
    x := 42
    fmt.Printf("x address: %p, value: %d\n", &x, x)
    
    byValue := func(n int) {
        fmt.Printf("param address: %p, value: %d\n", &n, n)
    }
    byValue(x)
    // Address berbeda! Copy terpisah
    
    // Pointer
    byPointer := func(n *int) {
        fmt.Printf("param address: %p, points to: %p, value: %d\n", 
            &n, n, *n)
    }
    byPointer(&x)
    // Parameter address beda, tapi points to sama!
    
    // Slice
    s := []int{1, 2, 3}
    fmt.Printf("slice header address: %p\n", &s)
    fmt.Printf("underlying array address: %p\n", s)
    
    modifySlice := func(slice []int) {
        fmt.Printf("param header address: %p\n", &slice)
        fmt.Printf("param underlying array: %p\n", slice)
    }
    modifySlice(s)
    // Header address beda, underlying array sama!
}
```

---

## 📚 Kesimpulan

1. **Go selalu pass by value**, tapi beberapa tipe menyimpan pointer internal
2. **Primitif & array**: Copy penuh, aman dari side effects
3. **Pointer**: Copy alamat, akses data sama
4. **Slice/Map/Channel**: Copy header, data shared
5. **Struct besar**: Pakai pointer untuk efisiensi
6. **Modifikasi asli**: Butuh pointer atau tipe reference (slice/map)

**Golden Rule:**
> Jika ragu, print memory address dengan `%p` untuk understand behavior!
