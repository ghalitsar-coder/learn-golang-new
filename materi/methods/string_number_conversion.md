# Konversi String & Number di Go

## Pengenalan

Di Go, konversi antara string dan number (angka) merupakan operasi umum dalam pengembangan aplikasi. Karena Go adalah bahasa yang ketat dalam hal tipe data, kita tidak bisa secara otomatis menggabungkan string dan number seperti di beberapa bahasa lain. Konversi ini harus dilakukan secara eksplisit menggunakan fungsi-fungsi yang disediakan oleh package `strconv`.

## Konversi Number ke String

### Menggunakan `strconv.Itoa` (Integer to ASCII)
Fungsi `strconv.Itoa` digunakan untuk mengkonversi integer ke string.

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    num := 42
    str := strconv.Itoa(num)
    fmt.Printf("Number: %d, String: %s\n", num, str)
}
```

### Menggunakan `strconv.FormatInt`
Fungsi `strconv.FormatInt` digunakan untuk mengkonversi `int64` ke string dalam basis tertentu (desimal, biner, oktal, heksadesimal).

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    num := int64(255)
    
    // Desimal (basis 10)
    decimal := strconv.FormatInt(num, 10)
    
    // Biner (basis 2)
    binary := strconv.FormatInt(num, 2)
    
    // Heksadesimal (basis 16)
    hex := strconv.FormatInt(num, 16)
    
    fmt.Printf("Desimal: %s, Biner: %s, Heksadesimal: %s\n", decimal, binary, hex)
}
```

### Menggunakan `strconv.FormatFloat`
Fungsi `strconv.FormatFloat` digunakan untuk mengkonversi float ke string.

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    f := 3.14159
    
    // Format float dengan 2 digit desimal
    str := strconv.FormatFloat(f, 'f', 2, 64)
    fmt.Printf("Float: %.5f, String: %s\n", f, str)
}
```

Format specifier untuk `strconv.FormatFloat`:
- `'f'`: Notasi titik desimal biasa
- `'e'`: Notasi ilmiah (misalnya 1.234e+05)
- `'E'`: Notasi ilmiah dengan huruf besar (misalnya 1.234E+05)
- `'g'`: Format otomatis (akan memilih 'e' atau 'f' berdasarkan nilai)

## Konversi String ke Number

### Menggunakan `strconv.Atoi` (ASCII to Integer)
Fungsi `strconv.Atoi` adalah jalan pintas untuk mengkonversi string ke integer (int).

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    str := "123"
    num, err := strconv.Atoi(str)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("String: %s, Number: %d\n", str, num)
    }
}
```

### Menggunakan `strconv.ParseInt`
Fungsi `strconv.ParseInt` digunakan untuk mengkonversi string ke integer dengan tipe data tertentu dan basis angka.

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    // Konversi dari string desimal ke int64
    str := "255"
    num, err := strconv.ParseInt(str, 10, 64)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("String: %s, Int64: %d\n", str, num)
    }
    
    // Konversi dari string heksadesimal ke int64
    hexStr := "ff"
    hexNum, err := strconv.ParseInt(hexStr, 16, 64)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Hex String: %s, Int64: %d\n", hexStr, hexNum)
    }
}
```

Parameter `strconv.ParseInt`:
- `str`: string yang akan dikonversi
- `base`: basis angka (0, 2, 8, 10, 16)
- `bitSize`: ukuran bit dari tipe integer (0, 8, 16, 32, 64)

### Menggunakan `strconv.ParseFloat`
Fungsi `strconv.ParseFloat` digunakan untuk mengkonversi string ke float64.

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    str := "3.14159"
    num, err := strconv.ParseFloat(str, 64)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("String: %s, Float64: %.5f\n", str, num)
    }
}
```

## Konversi Menggunakan `fmt.Sprintf`

Alternatif lain untuk mengkonversi number ke string adalah dengan menggunakan `fmt.Sprintf`:

```go
package main

import "fmt"

func main() {
    num := 123
    str := fmt.Sprintf("%d", num)
    fmt.Printf("Number: %d, String: %s\n", num, str)
    
    f := 3.14
    str2 := fmt.Sprintf("%.2f", f)
    fmt.Printf("Float: %.1f, String: %s\n", f, str2)
}
```

## Contoh Studi Kasus: Parsing Input Pengguna

Berikut adalah contoh implementasi konversi string ke number dalam kasus parsing input pengguna:

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    
    fmt.Print("Masukkan angka: ")
    input, _ := reader.ReadString('\n')
    input = strings.TrimSpace(input)
    
    // Konversi string ke integer
    num, err := strconv.Atoi(input)
    if err != nil {
        fmt.Printf("Input bukan angka: %v\n", err)
        return
    }
    
    fmt.Printf("Anda memasukkan angka: %d\n", num)
    
    // Contoh penggunaan angka untuk perhitungan
    result := num * 2
    fmt.Printf("Hasil dari %d * 2 = %d\n", num, result)
}
```

## Ringkasan

| Konversi | Fungsi | Contoh |
|----------|--------|--------|
| Integer ke String | `strconv.Itoa()` | `strconv.Itoa(42)` → `"42"` |
| Integer ke String (basis) | `strconv.FormatInt()` | `strconv.FormatInt(255, 16)` → `"ff"` |
| Float ke String | `strconv.FormatFloat()` | `strconv.FormatFloat(3.14, 'f', 2, 64)` → `"3.14"` |
| String ke Integer | `strconv.Atoi()` | `strconv.Atoi("123")` → `123` |
| String ke Integer (basis) | `strconv.ParseInt()` | `strconv.ParseInt("ff", 16, 64)` → `255` |
| String ke Float | `strconv.ParseFloat()` | `strconv.ParseFloat("3.14", 64)` → `3.14` |

## Tips dan Trik

1. Selalu periksa error saat melakukan parsing dari string ke number
2. Gunakan `fmt.Sprintf` untuk konversi ke string jika format output perlu disesuaikan
3. Untuk basis bilangan, pastikan basis yang digunakan sesuai dengan format string (misal: jika string dalam format heksadesimal, gunakan basis 16)
4. Ketika bekerja dengan bilangan besar, pertimbangkan untuk menggunakan tipe data yang sesuai (int64, float64)

## Latihan

1. Buat program yang mengonversi suhu dari Celcius ke Fahrenheit, dengan input berupa string
2. Buat fungsi yang mengembalikan jumlah digit dalam sebuah angka integer
3. Implementasikan program kalkulator sederhana yang menerima input string dan mengembalikan hasil perhitungan