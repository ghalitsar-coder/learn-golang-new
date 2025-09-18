# Number/Math Methods dalam Go

Dalam Go, kita menggunakan package `math` untuk operasi matematika yang kompleks, berbeda dengan JavaScript yang memiliki objek Math dengan berbagai method bawaan. Di Go, kita juga memiliki operator aritmatika standar untuk operasi dasar.

## 1. Operasi Dasar

Operasi aritmatika dasar dalam Go menggunakan operator:

```go
package main

import "fmt"

func main() {
    a, b := 10.0, 3.0
    
    // Penjumlahan
    sum := a + b
    fmt.Println("Sum:", sum) // Output: Sum: 13
    
    // Pengurangan
    diff := a - b
    fmt.Println("Difference:", diff) // Output: Difference: 7
    
    // Perkalian
    product := a * b
    fmt.Println("Product:", product) // Output: Product: 30
    
    // Pembagian
    quotient := a / b
    fmt.Println("Quotient:", quotient) // Output: Quotient: 3.3333333333333335
    
    // Modulo (hanya untuk integer)
    intA, intB := 10, 3
    remainder := intA % intB
    fmt.Println("Remainder:", remainder) // Output: Remainder: 1
}
```

## 2. Fungsi Matematika dari Package math

Package `math` menyediakan berbagai fungsi matematika:

```go
package main

import (
    "fmt"
    "math"
)

func main() {
    x := 4.0
    y := 3.14159
    
    // Akar kuadrat
    sqrt := math.Sqrt(x)
    fmt.Println("Square root of 4:", sqrt) // Output: Square root of 4: 2
    
    // Pangkat
    power := math.Pow(x, 2)
    fmt.Println("4 to the power of 2:", power) // Output: 4 to the power of 2: 16
    
    // Pembulatan
    rounded := math.Round(y)
    fmt.Println("Rounded 3.14159:", rounded) // Output: Rounded 3.14159: 3
    
    // Pembulatan ke atas
    ceil := math.Ceil(3.14)
    fmt.Println("Ceiling of 3.14:", ceil) // Output: Ceiling of 3.14: 4
    
    // Pembulatan ke bawah
    floor := math.Floor(3.14)
    fmt.Println("Floor of 3.14:", floor) // Output: Floor of 3.14: 3
    
    // Nilai absolut
    abs := math.Abs(-5.5)
    fmt.Println("Absolute value of -5.5:", abs) // Output: Absolute value of -5.5: 5.5
    
    // Fungsi trigonometri
    angle := math.Pi / 4 // 45 derajat
    sin := math.Sin(angle)
    cos := math.Cos(angle)
    tan := math.Tan(angle)
    fmt.Printf("Sin(45°): %.3f, Cos(45°): %.3f, Tan(45°): %.3f\n", sin, cos, tan)
    // Output: Sin(45°): 0.707, Cos(45°): 0.707, Tan(45°): 1.000
    
    // Logaritma
    log := math.Log(math.E) // log e = 1
    log10 := math.Log10(100) // log10 100 = 2
    fmt.Println("ln(e):", log) // Output: ln(e): 1
    fmt.Println("log10(100):", log10) // Output: log10(100): 2
    
    // Eksponensial
    exp := math.Exp(1) // e^1
    fmt.Println("e^1:", exp) // Output: e^1: 2.718281828459045
}
```

## 3. Konstanta Matematika

Package `math` juga menyediakan konstanta matematika:

```go
package main

import (
    "fmt"
    "math"
)

func main() {
    fmt.Println("Pi:", math.Pi)           // Output: Pi: 3.141592653589793
    fmt.Println("E:", math.E)             // Output: E: 2.718281828459045
    fmt.Println("Phi (Golden Ratio):", math.Phi) // Output: Phi (Golden Ratio): 1.618033988749895
    
    // Konstanta lainnya
    fmt.Println("Square root of 2:", math.Sqrt2)     // Output: Square root of 2: 1.4142135623730951
    fmt.Println("Square root of E:", math.SqrtE)     // Output: Square root of E: 1.6487212707001282
    fmt.Println("Square root of Pi:", math.SqrtPi)    // Output: Square root of Pi: 1.7724538509055159
    fmt.Println("Ln(2):", math.Ln2)                  // Output: Ln(2): 0.6931471805599453
    fmt.Println("Ln(10):", math.Ln10)                // Output: Ln(10): 2.302585092994046
    fmt.Println("Log10(2):", math.Log10E)            // Output: Log10(2): 0.4342944819032518
    fmt.Println("Log10(E):", math.Log2E)             // Output: Log10(E): 1.4426950408889634
}
```

## 4. Fungsi Utilitas untuk Angka

Berikut adalah beberapa fungsi utilitas untuk memanipulasi angka:

```go
package main

import (
    "fmt"
    "math"
    "strconv"
)

// Max mengembalikan nilai terbesar dari dua angka
func Max[T int | int32 | int64 | float32 | float64](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// Min mengembalikan nilai terkecil dari dua angka
func Min[T int | int32 | int64 | float32 | float64](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// Clamp membatasi nilai dalam rentang tertentu
func Clamp[T int | int32 | int64 | float32 | float64](value, min, max T) T {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

// Abs mengembalikan nilai absolut dari angka
func Abs[T int | int32 | int64 | float32 | float64](x T) T {
    if x < 0 {
        return -x
    }
    return x
}

// IsEven memeriksa apakah angka genap
func IsEven[T int | int32 | int64](n T) bool {
    return n%2 == 0
}

// IsOdd memeriksa apakah angka ganjil
func IsOdd[T int | int32 | int64](n T) bool {
    return n%2 != 0
}

// ToFixed membatasi angka desimal hingga n digit
func ToFixed(value float64, precision int) float64 {
    multiplier := math.Pow(10, float64(precision))
    return math.Round(value*multiplier) / multiplier
}

// FormatNumber memformat angka dengan pemisah ribuan
func FormatNumber(n int) string {
    in := strconv.FormatInt(int64(n), 10)
    numOfDigits := len(in)
    if n < 0 {
        numOfDigits-- // First character is the minus sign (don't want to add comma before minus sign)
    }
    numOfCommas := (numOfDigits - 1) / 3

    out := make([]byte, len(in)+numOfCommas)
    if n < 0 {
        in, out[0] = in[1:], '-'
    }

    for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
        out[j] = in[i]
        if i == 0 {
            return string(out)
        }
        if k++; k == 3 {
            j, k = j-1, 0
            out[j] = ','
        }
    }
}

func main() {
    // Menggunakan fungsi Max dan Min
    fmt.Println("Max of 5 and 3:", Max(5, 3)) // Output: Max of 5 and 3: 5
    fmt.Println("Min of 5 and 3:", Min(5, 3)) // Output: Min of 5 and 3: 3
    
    // Menggunakan fungsi Clamp
    fmt.Println("Clamp 15 between 0 and 10:", Clamp(15, 0, 10)) // Output: Clamp 15 between 0 and 10: 10
    fmt.Println("Clamp -5 between 0 and 10:", Clamp(-5, 0, 10)) // Output: Clamp -5 between 0 and 10: 0
    fmt.Println("Clamp 7 between 0 and 10:", Clamp(7, 0, 10))   // Output: Clamp 7 between 0 and 10: 7
    
    // Menggunakan fungsi Abs
    fmt.Println("Absolute value of -10:", Abs(-10)) // Output: Absolute value of -10: 10
    
    // Menggunakan fungsi IsEven dan IsOdd
    fmt.Println("Is 4 even?", IsEven(4)) // Output: Is 4 even? true
    fmt.Println("Is 5 odd?", IsOdd(5))   // Output: Is 5 odd? true
    
    // Menggunakan fungsi ToFixed
    fmt.Println("Pi fixed to 2 decimals:", ToFixed(math.Pi, 2)) // Output: Pi fixed to 2 decimals: 3.14
    
    // Menggunakan fungsi FormatNumber
    fmt.Println("Formatted number:", FormatNumber(1234567)) // Output: Formatted number: 1,234,567
}
```

## 5. Random Number Generation

Menghasilkan angka acak dalam Go:

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    // Menginisialisasi seed untuk generator angka acak
    rand.Seed(time.Now().UnixNano())
    
    // Menghasilkan angka acak antara 0 dan 100
    randomInt := rand.Intn(101) // 101 karena Intn menghasilkan angka dari 0 hingga n-1
    fmt.Println("Random integer (0-100):", randomInt)
    
    // Menghasilkan angka float64 antara 0.0 dan 1.0
    randomFloat := rand.Float64()
    fmt.Println("Random float (0.0-1.0):", randomFloat)
    
    // Menghasilkan angka dalam rentang tertentu
    min, max := 10.0, 20.0
    randomRange := min + rand.Float64()*(max-min)
    fmt.Printf("Random float (%.1f-%.1f): %.2f\n", min, max, randomRange)
    
    // Menghasilkan angka acak dari slice
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    randomElement := numbers[rand.Intn(len(numbers))]
    fmt.Println("Random element from slice:", randomElement)
}
```

## Contoh Implementasi Handler/Helper

Berikut contoh bagaimana kita bisa membuat handler/helper untuk number/math:

```go
package main

import (
    "fmt"
    "math"
    "math/rand"
    "strconv"
    "time"
)

// NumberHelper adalah struct untuk mengelompokkan method-method number
type NumberHelper struct{}

// Max mengembalikan nilai terbesar dari dua angka
func (nh NumberHelper) Max[T int | int32 | int64 | float32 | float64](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// Min mengembalikan nilai terkecil dari dua angka
func (nh NumberHelper) Min[T int | int32 | int64 | float32 | float64](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// Clamp membatasi nilai dalam rentang tertentu
func (nh NumberHelper) Clamp[T int | int32 | int64 | float32 | float64](value, min, max T) T {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

// Abs mengembalikan nilai absolut dari angka
func (nh NumberHelper) Abs[T int | int32 | int64 | float32 | float64](x T) T {
    if x < 0 {
        return -x
    }
    return x
}

// Round membulatkan angka ke desimal tertentu
func (nh NumberHelper) Round(value float64, precision int) float64 {
    multiplier := math.Pow(10, float64(precision))
    return math.Round(value*multiplier) / multiplier
}

// IsEven memeriksa apakah angka genap
func (nh NumberHelper) IsEven[T int | int32 | int64](n T) bool {
    return n%2 == 0
}

// IsOdd memeriksa apakah angka ganjil
func (nh NumberHelper) IsOdd[T int | int32 | int64](n T) bool {
    return n%2 != 0
}

// FormatWithCommas memformat angka dengan pemisah ribuan
func (nh NumberHelper) FormatWithCommas(n int) string {
    in := strconv.FormatInt(int64(n), 10)
    numOfDigits := len(in)
    if n < 0 {
        numOfDigits-- // First character is the minus sign
    }
    numOfCommas := (numOfDigits - 1) / 3

    out := make([]byte, len(in)+numOfCommas)
    if n < 0 {
        in, out[0] = in[1:], '-'
    }

    for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
        out[j] = in[i]
        if i == 0 {
            return string(out)
        }
        if k++; k == 3 {
            j, k = j-1, 0
            out[j] = ','
        }
    }
}

// RandomInt menghasilkan angka acak antara min dan max (inklusif)
func (nh NumberHelper) RandomInt(min, max int) int {
    return rand.Intn(max-min+1) + min
}

// RandomFloat menghasilkan angka float acak antara min dan max
func (nh NumberHelper) RandomFloat(min, max float64) float64 {
    return min + rand.Float64()*(max-min)
}

// Sum menghitung total dari slice angka
func (nh NumberHelper) Sum[T int | int32 | int64 | float32 | float64](numbers []T) T {
    var sum T
    for _, n := range numbers {
        sum += n
    }
    return sum
}

// Average menghitung rata-rata dari slice angka
func (nh NumberHelper) Average[T int | int32 | int64 | float32 | float64](numbers []T) float64 {
    if len(numbers) == 0 {
        return 0
    }
    sum := nh.Sum(numbers)
    return float64(sum) / float64(len(numbers))
}

func main() {
    // Menginisialisasi seed untuk generator angka acak
    rand.Seed(time.Now().UnixNano())
    
    helper := NumberHelper{}
    
    // Demonstrasi penggunaan method-method NumberHelper
    fmt.Println("=== Number Helper Demo ===")
    
    // Max dan Min
    fmt.Println("Max of 15 and 23:", helper.Max(15, 23)) // Output: Max of 15 and 23: 23
    fmt.Println("Min of 15 and 23:", helper.Min(15, 23)) // Output: Min of 15 and 23: 15
    
    // Clamp
    fmt.Println("Clamp 150 between 0 and 100:", helper.Clamp(150, 0, 100)) // Output: Clamp 150 between 0 and 100: 100
    
    // Abs
    fmt.Println("Absolute value of -42:", helper.Abs(-42)) // Output: Absolute value of -42: 42
    
    // Round
    fmt.Println("Round 3.14159 to 2 decimals:", helper.Round(3.14159, 2)) // Output: Round 3.14159 to 2 decimals: 3.14
    
    // IsEven dan IsOdd
    fmt.Println("Is 42 even?", helper.IsEven(42)) // Output: Is 42 even? true
    fmt.Println("Is 42 odd?", helper.IsOdd(42))   // Output: Is 42 odd? false
    
    // FormatWithCommas
    fmt.Println("Format 1234567:", helper.FormatWithCommas(1234567)) // Output: Format 1234567: 1,234,567
    
    // Random numbers
    fmt.Println("Random int (1-10):", helper.RandomInt(1, 10))
    fmt.Println("Random float (0.0-1.0):", helper.RandomFloat(0.0, 1.0))
    
    // Sum dan Average
    numbers := []int{10, 20, 30, 40, 50}
    fmt.Println("Numbers:", numbers)
    fmt.Println("Sum:", helper.Sum(numbers))           // Output: Sum: 150
    fmt.Println("Average:", helper.Average(numbers))   // Output: Average: 30
}
```

Dengan menggunakan pendekatan ini, kita bisa membuat koleksi fungsi utilitas untuk angka yang mirip dengan Math object di JavaScript. Perbedaan utamanya adalah bahwa dalam Go, kita menggunakan package `math` untuk fungsi matematika kompleks dan mengimplementasikan sendiri fungsi utilitas, sedangkan dalam JavaScript fungsi-fungsi tersebut sudah tersedia sebagai bagian dari objek Math.