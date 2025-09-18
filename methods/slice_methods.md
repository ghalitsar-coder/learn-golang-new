# Slice Methods dalam Go

Dalam Go, slice adalah struktur data yang mirip dengan array namun lebih fleksibel karena ukurannya bisa berubah. Untuk memanipulasi slice, kita menggunakan package `slices` (sejak Go 1.21) dan fungsi-fungsi bawaan Go.

## 1. Map

Dalam JavaScript, map digunakan untuk mengubah setiap elemen dalam array. Dalam Go, kita perlu mengimplementasikan sendiri fungsi ini:

```go
package main

import "fmt"

// Map mengaplikasikan fungsi transformasi ke setiap elemen slice
func Map[T any, R any](slice []T, transform func(T) R) []R {
    result := make([]R, len(slice))
    for i, v := range slice {
        result[i] = transform(v)
    }
    return result
}

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // Mengalikan setiap angka dengan 2
    doubled := Map(numbers, func(n int) int {
        return n * 2
    })
    
    fmt.Println(doubled) // Output: [2 4 6 8 10]
    
    // Mengubah angka menjadi string
    strings := Map(numbers, func(n int) string {
        return fmt.Sprintf("Number: %d", n)
    })
    
    fmt.Println(strings) // Output: [Number: 1 Number: 2 Number: 3 Number: 4 Number: 5]
}
```

## 2. Filter

Filter digunakan untuk memilih elemen-elemen yang memenuhi kondisi tertentu:

```go
package main

import "fmt"

// Filter memilih elemen-elemen yang memenuhi predicate
func Filter[T any](slice []T, predicate func(T) bool) []T {
    result := make([]T, 0)
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // Memilih hanya angka genap
    evenNumbers := Filter(numbers, func(n int) bool {
        return n%2 == 0
    })
    
    fmt.Println(evenNumbers) // Output: [2 4 6 8 10]
    
    // Memilih angka lebih besar dari 5
    greaterThanFive := Filter(numbers, func(n int) bool {
        return n > 5
    })
    
    fmt.Println(greaterThanFive) // Output: [6 7 8 9 10]
}
```

## 3. Reduce

Reduce digunakan untuk mengakumulasi nilai dari elemen-elemen slice:

```go
package main

import "fmt"

// Reduce mengakumulasi nilai dari elemen-elemen slice
func Reduce[T any, R any](slice []T, initial R, accumulator func(R, T) R) R {
    result := initial
    for _, v := range slice {
        result = accumulator(result, v)
    }
    return result
}

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // Menghitung total
    sum := Reduce(numbers, 0, func(acc int, n int) int {
        return acc + n
    })
    
    fmt.Println(sum) // Output: 15
    
    // Menghitung perkalian
    product := Reduce(numbers, 1, func(acc int, n int) int {
        return acc * n
    })
    
    fmt.Println(product) // Output: 120
    
    // Menggabungkan string
    words := []string{"Hello", " ", "World", "!"}
    sentence := Reduce(words, "", func(acc string, word string) string {
        return acc + word
    })
    
    fmt.Println(sentence) // Output: Hello World!
}
```

## 4. Find dan FindIndex

Find digunakan untuk mencari elemen pertama yang memenuhi kondisi:

```go
package main

import "fmt"

// Find mencari elemen pertama yang memenuhi predicate
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
    var zero T
    for _, v := range slice {
        if predicate(v) {
            return v, true
        }
    }
    return zero, false
}

// FindIndex mencari index elemen pertama yang memenuhi predicate
func FindIndex[T any](slice []T, predicate func(T) bool) int {
    for i, v := range slice {
        if predicate(v) {
            return i
        }
    }
    return -1
}

func main() {
    numbers := []int{1, 3, 5, 8, 9, 12}
    
    // Mencari angka genap pertama
    even, found := Find(numbers, func(n int) bool {
        return n%2 == 0
    })
    
    if found {
        fmt.Printf("First even number: %d\n", even) // Output: First even number: 8
    }
    
    // Mencari index angka lebih besar dari 10
    index := FindIndex(numbers, func(n int) bool {
        return n > 10
    })
    
    if index != -1 {
        fmt.Printf("Index of number > 10: %d\n", index) // Output: Index of number > 10: 5
    }
}
```

## 5. Some dan Every

Some memeriksa apakah setidaknya satu elemen memenuhi kondisi, dan Every memeriksa apakah semua elemen memenuhi kondisi:

```go
package main

import "fmt"

// Some memeriksa apakah setidaknya satu elemen memenuhi predicate
func Some[T any](slice []T, predicate func(T) bool) bool {
    for _, v := range slice {
        if predicate(v) {
            return true
        }
    }
    return false
}

// Every memeriksa apakah semua elemen memenuhi predicate
func Every[T any](slice []T, predicate func(T) bool) bool {
    for _, v := range slice {
        if !predicate(v) {
            return false
        }
    }
    return true
}

func main() {
    numbers := []int{2, 4, 6, 8, 10}
    
    // Memeriksa apakah ada angka ganjil
    hasOdd := Some(numbers, func(n int) bool {
        return n%2 != 0
    })
    
    fmt.Println(hasOdd) // Output: false
    
    // Memeriksa apakah semua angka genap
    allEven := Every(numbers, func(n int) bool {
        return n%2 == 0
    })
    
    fmt.Println(allEven) // Output: true
    
    mixedNumbers := []int{1, 2, 3, 4, 5}
    
    // Memeriksa apakah ada angka genap
    hasEven := Some(mixedNumbers, func(n int) bool {
        return n%2 == 0
    })
    
    fmt.Println(hasEven) // Output: true
    
    // Memeriksa apakah semua angka genap
    allEvenMixed := Every(mixedNumbers, func(n int) bool {
        return n%2 == 0
    })
    
    fmt.Println(allEvenMixed) // Output: false
}
```

## 6. Sort

Mengurutkan elemen-elemen dalam slice:

```go
package main

import (
    "fmt"
    "slices" // Package slices baru tersedia sejak Go 1.21
)

func main() {
    numbers := []int{5, 2, 8, 1, 9, 3}
    
    // Mengurutkan slice
    slices.Sort(numbers)
    fmt.Println(numbers) // Output: [1 2 3 5 8 9]
    
    // Mengurutkan dalam urutan terbalik
    slices.SortFunc(numbers, func(a, b int) int {
        return b - a // Mengurangi b dengan a untuk urutan terbalik
    })
    fmt.Println(numbers) // Output: [9 8 5 3 2 1]
    
    // Untuk versi Go sebelum 1.21, gunakan package sort:
    /*
    import "sort"
    
    numbers := []int{5, 2, 8, 1, 9, 3}
    sort.Ints(numbers)
    fmt.Println(numbers) // Output: [1 2 3 5 8 9]
    */
}
```

## 7. Reverse

Membalik urutan elemen dalam slice:

```go
package main

import (
    "fmt"
    "slices" // Package slices baru tersedia sejak Go 1.21
)

func main() {
    letters := []string{"a", "b", "c", "d", "e"}
    
    // Membalik urutan slice
    slices.Reverse(letters)
    fmt.Println(letters) // Output: [e d c b a]
    
    // Untuk versi Go sebelum 1.21, implementasi manual:
    /*
    func reverse[T any](slice []T) {
        for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
            slice[i], slice[j] = slice[j], slice[i]
        }
    }
    
    letters := []string{"a", "b", "c", "d", "e"}
    reverse(letters)
    fmt.Println(letters) // Output: [e d c b a]
    */
}
```

## 8. Concat dan Flatten

Menggabungkan slice dan meratakan slice bersarang:

```go
package main

import "fmt"

// Concat menggabungkan beberapa slice
func Concat[T any](slices ...[]T) []T {
    result := make([]T, 0)
    for _, slice := range slices {
        result = append(result, slice...)
    }
    return result
}

func main() {
    slice1 := []int{1, 2, 3}
    slice2 := []int{4, 5, 6}
    slice3 := []int{7, 8, 9}
    
    // Menggabungkan slice
    combined := Concat(slice1, slice2, slice3)
    fmt.Println(combined) // Output: [1 2 3 4 5 6 7 8 9]
    
    // Untuk flatten, kita bisa menggunakan pendekatan manual:
    nested := [][]int{{1, 2}, {3, 4}, {5, 6}}
    flattened := make([]int, 0)
    for _, inner := range nested {
        flattened = append(flattened, inner...)
    }
    fmt.Println(flattened) // Output: [1 2 3 4 5 6]
}
```

## 9. Unique

Menghapus elemen duplikat dari slice:

```go
package main

import "fmt"

// Unique menghapus elemen duplikat dari slice
func Unique[T comparable](slice []T) []T {
    seen := make(map[T]bool)
    result := make([]T, 0)
    
    for _, v := range slice {
        if !seen[v] {
            seen[v] = true
            result = append(result, v)
        }
    }
    
    return result
}

func main() {
    numbers := []int{1, 2, 2, 3, 3, 3, 4, 5, 5}
    
    uniqueNumbers := Unique(numbers)
    fmt.Println(uniqueNumbers) // Output: [1 2 3 4 5]
    
    words := []string{"apple", "banana", "apple", "orange", "banana"}
    uniqueWords := Unique(words)
    fmt.Println(uniqueWords) // Output: [apple banana orange]
}
```

## Contoh Implementasi Handler/Helper

Berikut contoh bagaimana kita bisa membuat handler/helper untuk slice:

```go
package main

import "fmt"

// SliceHelper adalah struct untuk mengelompokkan method-method slice
type SliceHelper struct{}

// Map mengaplikasikan fungsi transformasi ke setiap elemen slice
func (sh SliceHelper) Map[T any, R any](slice []T, transform func(T) R) []R {
    result := make([]R, len(slice))
    for i, v := range slice {
        result[i] = transform(v)
    }
    return result
}

// Filter memilih elemen-elemen yang memenuhi predicate
func (sh SliceHelper) Filter[T any](slice []T, predicate func(T) bool) []T {
    result := make([]T, 0)
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

// Reduce mengakumulasi nilai dari elemen-elemen slice
func (sh SliceHelper) Reduce[T any, R any](slice []T, initial R, accumulator func(R, T) R) R {
    result := initial
    for _, v := range slice {
        result = accumulator(result, v)
    }
    return result
}

// Find mencari elemen pertama yang memenuhi predicate
func (sh SliceHelper) Find[T any](slice []T, predicate func(T) bool) (T, bool) {
    var zero T
    for _, v := range slice {
        if predicate(v) {
            return v, true
        }
    }
    return zero, false
}

// Some memeriksa apakah setidaknya satu elemen memenuhi predicate
func (sh SliceHelper) Some[T any](slice []T, predicate func(T) bool) bool {
    for _, v := range slice {
        if predicate(v) {
            return true
        }
    }
    return false
}

// Every memeriksa apakah semua elemen memenuhi predicate
func (sh SliceHelper) Every[T any](slice []T, predicate func(T) bool) bool {
    for _, v := range slice {
        if !predicate(v) {
            return false
        }
    }
    return true
}

func main() {
    helper := SliceHelper{}
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // Menggunakan Map untuk mengalikan setiap angka dengan 2
    doubled := helper.Map(numbers, func(n int) int {
        return n * 2
    })
    fmt.Println("Doubled:", doubled) // Output: Doubled: [2 4 6 8 10 12 14 16 18 20]
    
    // Menggunakan Filter untuk memilih angka genap
    evenNumbers := helper.Filter(numbers, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println("Even numbers:", evenNumbers) // Output: Even numbers: [2 4 6 8 10]
    
    // Menggunakan Reduce untuk menghitung total
    sum := helper.Reduce(numbers, 0, func(acc int, n int) int {
        return acc + n
    })
    fmt.Println("Sum:", sum) // Output: Sum: 55
    
    // Menggunakan Find untuk mencari angka pertama lebih besar dari 5
    greaterThanFive, found := helper.Find(numbers, func(n int) bool {
        return n > 5
    })
    if found {
        fmt.Printf("First number > 5: %d\n", greaterThanFive) // Output: First number > 5: 6
    }
    
    // Menggunakan Some untuk memeriksa apakah ada angka genap
    hasEven := helper.Some(numbers, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println("Has even number:", hasEven) // Output: Has even number: true
    
    // Menggunakan Every untuk memeriksa apakah semua angka positif
    allPositive := helper.Every(numbers, func(n int) bool {
        return n > 0
    })
    fmt.Println("All positive:", allPositive) // Output: All positive: true
}
```

Dengan menggunakan pendekatan ini, kita bisa membuat koleksi fungsi utilitas yang mirip dengan method-method array di JavaScript seperti map, filter, reduce, find, some, every, dll. Perbedaan utamanya adalah bahwa dalam Go, kita perlu mengimplementasikan fungsi-fungsi ini sendiri atau menggunakan package `slices` (untuk Go 1.21+), sedangkan dalam JavaScript method-method tersebut sudah tersedia sebagai bagian dari objek Array.