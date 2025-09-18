# Map Methods dalam Go

Dalam Go, map adalah struktur data key-value yang tidak terurut. Berbeda dengan JavaScript yang memiliki built-in methods untuk object, dalam Go kita perlu mengimplementasikan sendiri fungsi-fungsi untuk memanipulasi map.

## 1. Keys dan Values

### Keys
Mendapatkan semua keys dari map:

```go
package main

import "fmt"

// Keys mengembalikan semua key dari map
func Keys[K comparable, V any](m map[K]V) []K {
    keys := make([]K, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}
```

### Values
Mendapatkan semua values dari map:

```go
// Values mengembalikan semua value dari map
func Values[K comparable, V any](m map[K]V) []V {
    values := make([]V, 0, len(m))
    for _, v := range m {
        values = append(values, v)
    }
    return values
}

func main() {
    studentScores := map[string]int{
        "Alice": 90,
        "Bob":   85,
        "Carol": 95,
    }
    
    // Mendapatkan semua keys
    names := Keys(studentScores)
    fmt.Println("Names:", names) // Output: Names: [Alice Bob Carol] (urutan bisa berbeda)
    
    // Mendapatkan semua values
    scores := Values(studentScores)
    fmt.Println("Scores:", scores) // Output: Scores: [90 85 95] (urutan bisa berbeda)
}
```

## 2. Map (Transformasi)

Mengaplikasikan transformasi ke setiap value dalam map:

```go
package main

import "fmt"

// MapValues mengaplikasikan fungsi transformasi ke setiap value dalam map
func MapValues[K comparable, V any, R any](m map[K]V, transform func(V) R) map[K]R {
    result := make(map[K]R)
    for k, v := range m {
        result[k] = transform(v)
    }
    return result
}

func main() {
    prices := map[string]float64{
        "apple":  1.2,
        "banana": 0.8,
        "orange": 1.5,
    }
    
    // Menaikkan harga 10%
    increasedPrices := MapValues(prices, func(price float64) float64 {
        return price * 1.1
    })
    
    fmt.Println("Original prices:", prices)
    fmt.Println("Increased prices:", increasedPrices)
    // Output:
    // Original prices: map[apple:1.2 banana:0.8 orange:1.5]
    // Increased prices: map[apple:1.3200000000000004 banana:0.88 orange:1.6500000000000001]
}
```

## 3. Filter

Memfilter pasangan key-value berdasarkan predicate:

```go
package main

import "fmt"

// FilterMap memfilter pasangan key-value berdasarkan predicate
func FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
    result := make(map[K]V)
    for k, v := range m {
        if predicate(k, v) {
            result[k] = v
        }
    }
    return result
}

func main() {
    studentScores := map[string]int{
        "Alice": 90,
        "Bob":   85,
        "Carol": 95,
        "David": 78,
        "Eve":   92,
    }
    
    // Memfilter siswa dengan nilai >= 90
    topStudents := FilterMap(studentScores, func(name string, score int) bool {
        return score >= 90
    })
    
    fmt.Println("Top students:", topStudents)
    // Output: Top students: map[Alice:90 Carol:95 Eve:92] (urutan bisa berbeda)
}
```

## 4. Reduce

Mengakumulasi nilai dari pasangan key-value dalam map:

```go
package main

import "fmt"

// ReduceMap mengakumulasi nilai dari pasangan key-value dalam map
func ReduceMap[K comparable, V any, R any](m map[K]V, initial R, accumulator func(R, K, V) R) R {
    result := initial
    for k, v := range m {
        result = accumulator(result, k, v)
    }
    return result
}

func main() {
    sales := map[string]int{
        "ProductA": 100,
        "ProductB": 150,
        "ProductC": 200,
    }
    
    // Menghitung total penjualan
    totalSales := ReduceMap(sales, 0, func(acc int, product string, quantity int) int {
        return acc + quantity
    })
    
    fmt.Println("Total sales:", totalSales) // Output: Total sales: 450
    
    // Membuat string deskriptif dari penjualan
    salesReport := ReduceMap(sales, "", func(acc string, product string, quantity int) string {
        if acc == "" {
            return fmt.Sprintf("%s: %d", product, quantity)
        }
        return fmt.Sprintf("%s, %s: %d", acc, product, quantity)
    })
    
    fmt.Println("Sales report:", salesReport)
    // Output: Sales report: ProductA: 100, ProductB: 150, ProductC: 200 (urutan bisa berbeda)
}
```

## 5. Merge

Menggabungkan dua atau lebih map:

```go
package main

import "fmt"

// MergeMap menggabungkan dua atau lebih map
func MergeMap[K comparable, V any](maps ...map[K]V) map[K]V {
    result := make(map[K]V)
    for _, m := range maps {
        for k, v := range m {
            result[k] = v
        }
    }
    return result
}

func main() {
    defaultConfig := map[string]string{
        "host": "localhost",
        "port": "8080",
    }
    
    userConfig := map[string]string{
        "port": "3000",
        "ssl":  "true",
    }
    
    // Menggabungkan konfigurasi, userConfig akan menimpa defaultConfig jika ada key yang sama
    finalConfig := MergeMap(defaultConfig, userConfig)
    
    fmt.Println("Final config:", finalConfig)
    // Output: Final config: map[host:localhost port:3000 ssl:true] (urutan bisa berbeda)
}
```

## 6. HasKey

Memeriksa apakah map memiliki key tertentu:

```go
package main

import "fmt"

// HasKey memeriksa apakah map memiliki key tertentu
func HasKey[K comparable, V any](m map[K]V, key K) bool {
    _, exists := m[key]
    return exists
}

func main() {
    inventory := map[string]int{
        "apple":  10,
        "banana": 5,
        "orange": 8,
    }
    
    // Memeriksa apakah ada stok apple
    hasApples := HasKey(inventory, "apple")
    fmt.Println("Has apples:", hasApples) // Output: Has apples: true
    
    // Memeriksa apakah ada stok grapes
    hasGrapes := HasKey(inventory, "grapes")
    fmt.Println("Has grapes:", hasGrapes) // Output: Has grapes: false
}
```

## 7. Pick dan Omit

Memilih atau mengabaikan key tertentu dari map:

```go
package main

import "fmt"

// PickMap memilih hanya key-key tertentu dari map
func PickMap[K comparable, V any](m map[K]V, keys []K) map[K]V {
    result := make(map[K]V)
    for _, key := range keys {
        if value, exists := m[key]; exists {
            result[key] = value
        }
    }
    return result
}

// OmitMap mengabaikan key-key tertentu dari map
func OmitMap[K comparable, V any](m map[K]V, keys []K) map[K]V {
    // Membuat map untuk pencarian cepat
    omitKeys := make(map[K]bool)
    for _, key := range keys {
        omitKeys[key] = true
    }
    
    result := make(map[K]V)
    for k, v := range m {
        if !omitKeys[k] {
            result[k] = v
        }
    }
    return result
}

func main() {
    user := map[string]interface{}{
        "name":    "Alice",
        "email":   "alice@example.com",
        "age":     25,
        "address": "123 Main St",
        "phone":   "555-1234",
    }
    
    // Memilih hanya field tertentu
    publicInfo := PickMap(user, []string{"name", "email"})
    fmt.Println("Public info:", publicInfo)
    // Output: Public info: map[email:alice@example.com name:Alice] (urutan bisa berbeda)
    
    // Mengabaikan field tertentu
    privateInfo := OmitMap(user, []string{"phone", "address"})
    fmt.Println("Private info:", privateInfo)
    // Output: Private info: map[age:25 email:alice@example.com name:Alice] (urutan bisa berbeda)
}
```

## Contoh Implementasi Handler/Helper

Berikut contoh bagaimana kita bisa membuat handler/helper untuk map:

```go
package main

import "fmt"

// MapHelper adalah struct untuk mengelompokkan method-method map
type MapHelper struct{}

// Keys mengembalikan semua key dari map
func (mh MapHelper) Keys[K comparable, V any](m map[K]V) []K {
    keys := make([]K, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}

// Values mengembalikan semua value dari map
func (mh MapHelper) Values[K comparable, V any](m map[K]V) []V {
    values := make([]V, 0, len(m))
    for _, v := range m {
        values = append(values, v)
    }
    return values
}

// MapValues mengaplikasikan fungsi transformasi ke setiap value dalam map
func (mh MapHelper) MapValues[K comparable, V any, R any](m map[K]V, transform func(V) R) map[K]R {
    result := make(map[K]R)
    for k, v := range m {
        result[k] = transform(v)
    }
    return result
}

// FilterMap memfilter pasangan key-value berdasarkan predicate
func (mh MapHelper) FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
    result := make(map[K]V)
    for k, v := range m {
        if predicate(k, v) {
            result[k] = v
        }
    }
    return result
}

// ReduceMap mengakumulasi nilai dari pasangan key-value dalam map
func (mh MapHelper) ReduceMap[K comparable, V any, R any](m map[K]V, initial R, accumulator func(R, K, V) R) R {
    result := initial
    for k, v := range m {
        result = accumulator(result, k, v)
    }
    return result
}

// MergeMap menggabungkan dua atau lebih map
func (mh MapHelper) MergeMap[K comparable, V any](maps ...map[K]V) map[K]V {
    result := make(map[K]V)
    for _, m := range maps {
        for k, v := range m {
            result[k] = v
        }
    }
    return result
}

// HasKey memeriksa apakah map memiliki key tertentu
func (mh MapHelper) HasKey[K comparable, V any](m map[K]V, key K) bool {
    _, exists := m[key]
    return exists
}

// PickMap memilih hanya key-key tertentu dari map
func (mh MapHelper) PickMap[K comparable, V any](m map[K]V, keys []K) map[K]V {
    result := make(map[K]V)
    for _, key := range keys {
        if value, exists := m[key]; exists {
            result[key] = value
        }
    }
    return result
}

// OmitMap mengabaikan key-key tertentu dari map
func (mh MapHelper) OmitMap[K comparable, V any](m map[K]V, keys []K) map[K]V {
    // Membuat map untuk pencarian cepat
    omitKeys := make(map[K]bool)
    for _, key := range keys {
        omitKeys[key] = true
    }
    
    result := make(map[K]V)
    for k, v := range m {
        if !omitKeys[k] {
            result[k] = v
        }
    }
    return result
}

func main() {
    helper := MapHelper{}
    
    // Contoh penggunaan dengan data mahasiswa
    students := map[string]map[string]interface{}{
        "alice": {
            "name":  "Alice",
            "age":   20,
            "score": 85.5,
        },
        "bob": {
            "name":  "Bob",
            "age":   22,
            "score": 92.0,
        },
        "carol": {
            "name":  "Carol",
            "age":   19,
            "score": 78.5,
        },
    }
    
    // Mendapatkan semua nama mahasiswa
    names := helper.MapValues(students, func(student map[string]interface{}) interface{} {
        return student["name"]
    })
    fmt.Println("Student names:", names)
    
    // Memfilter mahasiswa dengan nilai di atas 80
    highAchievers := helper.FilterMap(students, func(id string, student map[string]interface{}) bool {
        score := student["score"].(float64)
        return score > 80
    })
    fmt.Println("High achievers:", highAchievers)
    
    // Menghitung rata-rata nilai
    totalScore := helper.ReduceMap(students, 0.0, func(acc float64, id string, student map[string]interface{}) float64 {
        return acc + student["score"].(float64)
    })
    averageScore := totalScore / float64(len(students))
    fmt.Printf("Average score: %.2f\n", averageScore)
    
    // Menambahkan field baru ke setiap mahasiswa
    studentsWithGrade := helper.MapValues(students, func(student map[string]interface{}) map[string]interface{} {
        score := student["score"].(float64)
        grade := "C"
        if score >= 90 {
            grade = "A"
        } else if score >= 80 {
            grade = "B"
        }
        
        // Membuat copy dari student dan menambahkan grade
        studentWithGrade := make(map[string]interface{})
        for k, v := range student {
            studentWithGrade[k] = v
        }
        studentWithGrade["grade"] = grade
        
        return studentWithGrade
    })
    
    fmt.Println("Students with grades:", studentsWithGrade)
}
```

Dengan menggunakan pendekatan ini, kita bisa membuat koleksi fungsi utilitas untuk map yang mirip dengan method-method object di JavaScript. Perbedaan utamanya adalah bahwa dalam Go, kita perlu mengimplementasikan fungsi-fungsi ini sendiri, sedangkan dalam JavaScript method-method tersebut sudah tersedia sebagai bagian dari objek Object atau melalui library seperti Lodash.