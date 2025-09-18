# String Methods dalam Go

Dalam Go, string adalah tipe data primitif yang immutable (tidak dapat diubah). Untuk memanipulasi string, kita menggunakan package `strings` yang menyediakan berbagai fungsi untuk operasi string.

## 1. Split dan Join

### Split
Memecah string menjadi slice berdasarkan pemisah.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    text := "apple,banana,orange"
    parts := strings.Split(text, ",")
    fmt.Println(parts) // Output: [apple banana orange]
    
    // Split dengan pemisah yang tidak ada
    text2 := "hello world"
    parts2 := strings.Split(text2, ",")
    fmt.Println(parts2) // Output: [hello world]
    
    // Split dengan string kosong
    text3 := "hello"
    chars := strings.Split(text3, "")
    fmt.Println(chars) // Output: [h e l l o]
}
```

### Join
Menggabungkan slice string menjadi satu string dengan pemisah tertentu.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    fruits := []string{"apple", "banana", "orange"}
    result := strings.Join(fruits, ", ")
    fmt.Println(result) // Output: apple, banana, orange
    
    // Join dengan pemisah kosong
    chars := []string{"h", "e", "l", "l", "o"}
    word := strings.Join(chars, "")
    fmt.Println(word) // Output: hello
}
```

## 2. Case Manipulation

### ToUpper dan ToLower
Mengubah huruf menjadi kapital atau kecil.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    text := "Hello World"
    
    upper := strings.ToUpper(text)
    fmt.Println(upper) // Output: HELLO WORLD
    
    lower := strings.ToLower(text)
    fmt.Println(lower) // Output: hello world
}
```

### Title
Mengubah huruf pertama setiap kata menjadi kapital.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    text := "hello world from go"
    title := strings.Title(text)
    fmt.Println(title) // Output: Hello World From Go
}
```

## 3. Trimming

### Trim
Menghapus karakter dari awal dan akhir string.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    text := "  hello world  "
    trimmed := strings.TrimSpace(text)
    fmt.Println(trimmed) // Output: hello world
    
    // Trim dengan karakter tertentu
    text2 := "---hello---"
    trimmed2 := strings.Trim(text2, "-")
    fmt.Println(trimmed2) // Output: hello
    
    // TrimPrefix dan TrimSuffix
    text3 := "hello world"
    prefixRemoved := strings.TrimPrefix(text3, "hello")
    fmt.Println(prefixRemoved) // Output:  world
    
    suffixRemoved := strings.TrimSuffix(text3, "world")
    fmt.Println(suffixRemoved) // Output: hello 
}
```

## 4. Searching dan Replacing

### Contains
Memeriksa apakah string mengandung substring tertentu.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    text := "hello world"
    
    hasHello := strings.Contains(text, "hello")
    fmt.Println(hasHello) // Output: true
    
    hasGo := strings.Contains(text, "go")
    fmt.Println(hasGo) // Output: false
}
```

### Index
Mencari posisi pertama dari substring.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    text := "hello world"
    
    index := strings.Index(text, "world")
    fmt.Println(index) // Output: 6
    
    notFound := strings.Index(text, "go")
    fmt.Println(notFound) // Output: -1
}
```

### Replace
Mengganti substring dengan string lain.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    text := "hello world, hello go"
    
    // Replace hanya satu kali
    replaced := strings.Replace(text, "hello", "hi", 1)
    fmt.Println(replaced) // Output: hi world, hello go
    
    // Replace semua
    replacedAll := strings.Replace(text, "hello", "hi", -1)
    fmt.Println(replacedAll) // Output: hi world, hi go
    
    // Atau menggunakan ReplaceAll
    replacedAll2 := strings.ReplaceAll(text, "hello", "hi")
    fmt.Println(replacedAll2) // Output: hi world, hi go
}
```

## 5. Repeating dan Padding

### Repeat
Mengulang string sebanyak n kali.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    text := "go"
    repeated := strings.Repeat(text, 3)
    fmt.Println(repeated) // Output: gogogo
    
    // Untuk membuat separator
    separator := strings.Repeat("-", 20)
    fmt.Println(separator) // Output: --------------------
}
```

## 6. Comparison

### Compare dan EqualFold
Membandingkan dua string.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    str1 := "Hello"
    str2 := "hello"
    
    // Compare case sensitive
    result := strings.Compare(str1, str2)
    fmt.Println(result) // Output: -1 (karena "Hello" < "hello" dalam ASCII)
    
    // EqualFold case insensitive
    equal := strings.EqualFold(str1, str2)
    fmt.Println(equal) // Output: true
}
```

## 7. Fields
Memecah string berdasarkan whitespace.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    text := "  hello   world  from   go  "
    fields := strings.Fields(text)
    fmt.Println(fields) // Output: [hello world from go]
}
```

## 8. HasPrefix dan HasSuffix
Memeriksa awalan dan akhiran string.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    url := "https://example.com"
    
    hasHTTP := strings.HasPrefix(url, "https://")
    fmt.Println(hasHTTP) // Output: true
    
    hasCom := strings.HasSuffix(url, ".com")
    fmt.Println(hasCom) // Output: true
}
```

## Contoh Implementasi Handler/Helper

Berikut contoh bagaimana kita bisa membuat handler/helper untuk string:

```go
package main

import (
    "fmt"
    "strings"
)

// StringHelper adalah struct untuk mengelompokkan method-method string
type StringHelper struct{}

// Split memecah string berdasarkan pemisah
func (sh StringHelper) Split(s, sep string) []string {
    return strings.Split(s, sep)
}

// Map mengaplikasikan fungsi transformasi ke setiap bagian string setelah di-split
func (sh StringHelper) Map(s, sep string, transform func(string) string) []string {
    parts := sh.Split(s, sep)
    result := make([]string, len(parts))
    
    for i, part := range parts {
        result[i] = transform(part)
    }
    
    return result
}

// Filter memfilter bagian string berdasarkan predicate
func (sh StringHelper) Filter(s, sep string, predicate func(string) bool) []string {
    parts := sh.Split(s, sep)
    result := make([]string, 0)
    
    for _, part := range parts {
        if predicate(part) {
            result = append(result, part)
        }
    }
    
    return result
}

// Join menggabungkan slice string dengan pemisah
func (sh StringHelper) Join(elements []string, sep string) string {
    return strings.Join(elements, sep)
}

func main() {
    helper := StringHelper{}
    text := "apple,banana,orange,grape"
    
    // Menggunakan Map untuk mengubah semua menjadi kapital
    capitalized := helper.Map(text, ",", strings.ToUpper)
    fmt.Println(capitalized) // Output: [APPLE BANANA ORANGE GRAPE]
    
    // Menggunakan Filter untuk memilih hanya yang mengandung "a"
    withA := helper.Filter(text, ",", func(s string) bool {
        return strings.Contains(s, "a")
    })
    fmt.Println(withA) // Output: [apple banana orange grape]
    
    // Menggabungkan kembali dengan pemisah berbeda
    result := helper.Join(withA, " | ")
    fmt.Println(result) // Output: apple | banana | orange | grape
}
```

Ini adalah beberapa method dasar untuk manipulasi string dalam Go. Berbeda dengan JavaScript di mana method-method ini adalah bagian dari objek string itu sendiri, dalam Go kita menggunakan fungsi-fungsi dari package `strings` atau membuat method-method kita sendiri dengan receiver.