# Generic dalam Go

## Daftar Isi
1. [Pendahuluan](#pendahuluan)
2. [Konsep Dasar Generic](#konsep-dasar-generic)
3. [Type Parameters](#type-parameters)
4. [Constraints](#constraints)
5. [Predeclared Constraints](#predeclared-constraints)
6. [Contoh Implementasi](#contoh-implementasi)
7. [Generic Function](#generic-function)
8. [Generic Type](#generic-type)
9. [Kesimpulan](#kesimpulan)

## Pendahuluan

Generic adalah fitur yang diperkenalkan dalam Go 1.18 yang memungkinkan kita untuk menulis kode yang dapat bekerja dengan berbagai tipe data tanpa kehilangan keamanan tipe (type safety). Sebelum generic, kita sering menggunakan interface{} yang memerlukan casting dan berisiko terhadap type safety.

Generic memungkinkan kita untuk menulis fungsi dan tipe data yang fleksibel namun tetap menjaga type safety dan performa. Dengan generic, kita dapat menulis kode yang lebih reusable dan mengurangi duplikasi kode.

## Konsep Dasar Generic

Generic memungkinkan kita untuk menulis fungsi dan tipe data yang dapat bekerja dengan berbagai tipe tanpa mengorbankan kinerja atau keamanan tipe. Ini berbeda dari pendekatan tradisional di Go menggunakan interface{} yang memerlukan casting dan tidak type-safe.

Dalam konteks generic:
- Kita dapat mendefinisikan fungsi yang bekerja dengan berbagai tipe
- Kita dapat membuat tipe data yang dapat menampung berbagai tipe
- Kita dapat mempertahankan type safety selama kompilasi

### Type Parameters

Type parameters adalah cara untuk menentukan tipe yang akan digunakan dalam fungsi atau tipe generic. Type parameters ditulis dalam kurung siku [].

Contoh sederhana:
```go
func NamaFungsi[TipeParameter](parameter TipeParameter) {
    // implementasi
}
```

Dalam definisi ini:
- `TipeParameter` adalah nama dari type parameter
- `T` atau `Type` adalah nama yang umum digunakan untuk type parameter
- Type parameter didefinisikan dalam kurung siku sebelum daftar parameter fungsi

Type parameters dapat memiliki beberapa parameter:
```go
func NamaFungsi[T, U any](param1 T, param2 U) {
    // implementasi
}
```

## Constraints

Constraints menentukan tipe apa yang dapat digunakan dengan type parameter. Dalam Go, constraints didefinisikan menggunakan interfaces.

Sebuah constraint menentukan operasi apa yang dapat dilakukan pada tipe yang digunakan dengan type parameter. Misalnya, jika kita ingin membandingkan dua nilai dengan operator ==, maka tipe tersebut harus memenuhi constraint comparable.

### Mendefinisikan Custom Constraints

Kita dapat membuat custom constraints dengan mendefinisikan interface yang menentukan perilaku yang diperlukan:

```go
// Constraint untuk tipe yang mendukung operasi numerik dasar
type Number interface {
    int | int8 | int16 | int32 | int64 | 
    uint | uint8 | uint16 | uint32 | uint64 | 
    float32 | float64
}

// Fungsi generic yang hanya menerima tipe numerik
func Add[T Number](a, b T) T {
    return a + b
}

// Penggunaan:
func main() {
    result1 := Add(1, 2)        // int
    result2 := Add(1.5, 2.7)    // float64
    fmt.Println(result1, result2) // 3 4.2
}
```

### Approximation Constraint

Kita juga dapat menggunakan approximation constraint (~) untuk menerima tipe yang memiliki tipe dasar tertentu:

```go
type Integer interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}

func Sum[T Integer](slice []T) T {
    var sum T
    for _, v := range slice {
        sum += v
    }
    return sum
}
```

Dalam contoh ini, fungsi Sum akan menerima tipe int, int8, int16, int32, int64, dan juga tipe yang didefinisikan berdasarkan tipe-tipe tersebut (seperti type MyInt int).

### Union Constraint

Union constraint memungkinkan kita untuk menentukan beberapa tipe yang dapat digunakan:

```go
type StringOrInt interface {
    string | int
}

func PrintStringOrInt[T StringOrInt](value T) {
    fmt.Println(value)
}
```

## Predeclared Constraints

Go menyediakan beberapa constraints bawaan yang umum digunakan:

- `any`: Constraint yang memungkinkan semua tipe. Ini identik dengan interface{}.
- `comparable`: Constraint untuk tipe yang dapat dibandingkan dengan operator == dan !=, seperti integers, floats, strings, booleans, structs dengan comparable fields, arrays dengan comparable elements, dan channels.

Contoh penggunaan:
```go
// Fungsi yang menerima semua tipe
func PrintValue[T any](value T) {
    fmt.Println(value)
}

// Fungsi yang hanya menerima tipe yang dapat dibandingkan
func Contains[T comparable](slice []T, item T) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}
```

## Contoh Implementasi

### Generic Function

Berikut adalah beberapa contoh fungsi generic yang umum digunakan:

**Contoh 1: Fungsi Swap**
```go
func Swap[T any](a, b *T) {
    *a, *b = *b, *a
}

// Penggunaan:
func main() {
    x, y := 10, 20
    fmt.Println("Sebelum swap:", x, y)
    Swap(&x, &y)
    fmt.Println("Setelah swap:", x, y)
    
    s1, s2 := "hello", "world"
    fmt.Println("Sebelum swap:", s1, s2)
    Swap(&s1, &s2)
    fmt.Println("Setelah swap:", s1, s2)
}
```

**Contoh 2: Fungsi Map**
```go
func Map[T any, R any](slice []T, fn func(T) R) []R {
    result := make([]R, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Penggunaan:
func main() {
    numbers := []int{1, 2, 3, 4, 5}
    doubled := Map(numbers, func(n int) int {
        return n * 2
    })
    fmt.Println(doubled) // [2 4 6 8 10]
    
    words := []string{"go", "generic", "example"}
    uppercased := Map(words, func(s string) string {
        return strings.ToUpper(s)
    })
    fmt.Println(uppercased) // [GO GENERIC EXAMPLE]
}
```

### Generic Type

Generic juga dapat digunakan untuk membuat tipe data yang fleksibel:

**Contoh 1: Stack Generic**
```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

func (s *Stack[T]) IsEmpty() bool {
    return len(s.items) == 0
}

// Penggunaan:
func main() {
    // Stack untuk integer
    intStack := &Stack[int]{}
    intStack.Push(1)
    intStack.Push(2)
    val, ok := intStack.Pop()
    fmt.Println(val, ok) // 2 true
    
    // Stack untuk string
    stringStack := &Stack[string]{}
    stringStack.Push("hello")
    stringStack.Push("world")
    val2, ok2 := stringStack.Pop()
    fmt.Println(val2, ok2) // world true
}
```

**Contoh 2: Pair Generic**
```go
type Pair[T, U any] struct {
    First  T
    Second U
}

func NewPair[T, U any](first T, second U) Pair[T, U] {
    return Pair[T, U]{First: first, Second: second}
}

func (p Pair[T, U]) GetFirst() T {
    return p.First
}

func (p Pair[T, U]) GetSecond() U {
    return p.Second
}

// Penggunaan:
func main() {
    // Pair dengan int dan string
    pair1 := NewPair(1, "satu")
    fmt.Println(pair1.GetFirst(), pair1.GetSecond()) // 1 satu
    
    // Pair dengan float64 dan bool
    pair2 := NewPair(3.14, true)
    fmt.Println(pair2.GetFirst(), pair2.GetSecond()) // 3.14 true
}
```

## Kesimpulan

Generic adalah fitur penting dalam Go yang memungkinkan kita untuk menulis kode yang lebih fleksibel, reusable, dan tetap menjaga type safety. Dengan generic, kita dapat menghindari duplikasi kode dan membuat library yang lebih umum dan kuat.

Keuntungan menggunakan generic:
1. **Type safety** - Kesalahan tipe terdeteksi saat kompilasi
2. **Performa** - Tidak ada overhead runtime seperti pada interface{}
3. **Reusability** - Satu implementasi dapat digunakan untuk berbagai tipe
4. **Maintainability** - Mengurangi duplikasi kode

### Kapan Menggunakan Generic

Generic sangat berguna dalam situasi berikut:
- Membuat struktur data seperti list, stack, queue, tree, dll.
- Membuat fungsi utilitas yang bekerja dengan berbagai tipe (seperti Map, Filter, Reduce)
- Membuat library yang dapat digunakan oleh berbagai tipe data
- Menghindari duplikasi kode untuk operasi yang sama pada tipe berbeda

### Best Practices

1. Gunakan generic hanya ketika diperlukan - jangan memaksa penggunaan generic jika tidak memberikan manfaat
2. Gunakan constraints yang sesuai untuk memastikan type safety
3. Gunakan nama type parameter yang jelas dan deskriptif
4. Pertimbangkan penggunaan predeclared constraints (any, comparable) jika memadai
5. Gunakan approximation constraint (~) jika ingin menerima tipe yang didefinisikan berdasarkan tipe dasar

### Limitations

Meskipun generic sangat berguna, ada beberapa keterbatasan yang perlu diperhatikan:
- Generic tidak dapat digunakan dengan method receiver pada interface
- Tidak semua operasi tersedia untuk semua tipe dalam constraint
- Kompleksitas kode bisa meningkat jika generic digunakan secara berlebihan

Penggunaan generic yang tepat akan membuat kode Go Anda lebih bersih, aman, dan mudah dikelola. Generic adalah alat yang kuat, namun seperti halnya alat lainnya, penggunaannya harus bijak dan sesuai dengan kebutuhan.