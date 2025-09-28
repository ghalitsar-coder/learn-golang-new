# Zero Value di Go

## Pengertian

Di bahasa pemrograman Go, **Zero Value** adalah nilai bawaan (default) yang secara otomatis diberikan kepada variabel ketika variabel tersebut dideklarasikan tanpa diberi nilai awal (inisialisasi). Ini berlaku untuk semua tipe data dasar (primitive types) maupun tipe data kompleks seperti struct, array, slice, map, pointer, dan channel.

## Tujuan Zero Value

Zero value memastikan bahwa **setiap variabel di Go selalu memiliki nilai yang valid** sejak dideklarasikan, mencegah terjadinya nilai tak terdefinisi (undefined) atau nilai acak seperti yang bisa terjadi di bahasa pemrograman lain (misalnya C/C++). Ini mendukung keamanan dan keandalan program Go.

## Zero Value untuk Tipe Data Dasar

Berikut adalah zero value untuk beberapa tipe data dasar di Go:

| Tipe Data      | Zero Value |
|----------------|------------|
| `bool`         | `false`    |
| `int`          | `0`        |
| `int8`         | `0`        |
| `int16`        | `0`        |
| `int32`        | `0`        |
| `int64`        | `0`        |
| `uint`         | `0`        |
| `uint8`        | `0`        |
| `uint16`       | `0`        |
| `uint32`       | `0`        |
| `uint64`       | `0`        |
| `uintptr`      | `0`        |
| `float32`      | `0.0`      |
| `float64`      | `0.0`      |
| `complex64`    | `0+0i`     |
| `complex128`   | `0+0i`     |
| `string`       | `""` (string kosong) |
| `rune`         | `0` (sama dengan `int32`) |
| `byte`         | `0` (sama dengan `uint8`) |
| `uintptr`      | `0`        |

## Contoh Penggunaan

Contoh kode Go yang menunjukkan zero value:

```go
package main

import "fmt"

func main() {
	var (
		boolVal   bool
		intVal    int
		floatVal  float64
		stringVal string
	)

	fmt.Printf("Zero value untuk bool: %t\n", boolVal)        // Output: false
	fmt.Printf("Zero value untuk int: %d\n", intVal)          // Output: 0
	fmt.Printf("Zero value untuk float64: %f\n", floatVal)    // Output: 0.000000
	fmt.Printf("Zero value untuk string: '%s'\n", stringVal)  // Output: ''
}
```

## Zero Value untuk Tipe Data Kompleks

### Array
Zero value dari array adalah array dengan setiap elemen diinisialisasi ke zero value dari tipe elemen tersebut.

```go
var arr [3]int // Zero value: [0, 0, 0]
```

### Slice
Zero value dari slice adalah `nil`. Slice yang nil berbeda dari slice kosong (`[]int{}`), meskipun keduanya bisa digunakan secara fungsional mirip dalam banyak kasus.

```go
var slice []int // Zero value: nil
```

### Map
Zero value dari map adalah `nil`.

```go
var m map[string]int // Zero value: nil
```

### Struct
Zero value dari struct adalah struct di mana setiap field-nya diinisialisasi ke zero value dari tipe field tersebut.

```go
type Person struct {
	Name string
	Age  int
}

var p Person // Zero value: Person{Name: "", Age: 0}
```

### Pointer
Zero value dari pointer adalah `nil`.

```go
var ptr *int // Zero value: nil
```

### Channel
Zero value dari channel adalah `nil`.

```go
var ch chan int // Zero value: nil
```

## Perbedaan Zero Value dan Nilai Falsy

Penting dicatat bahwa zero value **tidak** sama dengan konsep *falsy* di bahasa seperti JavaScript. Zero value adalah nilai *nyata* yang diberikan kepada variabel, sedangkan *falsy* adalah nilai-nilai yang dianggap `false` dalam konteks boolean.

Dalam Go, tidak ada konsep *falsy*. Jika kamu ingin mengecek apakah suatu nilai adalah zero value (misalnya string kosong atau angka nol), kamu harus membandingkannya secara eksplisit:

Contoh:
```go
var s string // zero value: ""
var i int    // zero value: 0

// Tidak bisa: if s { ... } // Error: non-boolean condition
 // Harus: if s != "" { ... } // Cek eksplisit
 // Harus: if i == 0 { ... } // Cek eksplisit
```

## Kesimpulan

Zero value adalah nilai bawaan yang diberikan Go kepada variabel yang dideklarasikan tanpa nilai awal. Ini adalah fitur penting untuk menjaga keamanan dan keandalan kode Go, karena menjamin bahwa tidak ada variabel yang tidak memiliki nilai. Pemahaman tentang zero value sangat penting untuk menghindari kebingungan, terutama saat bekerja dengan pointer, slice, map, dan struct, di mana zero value bisa berupa `nil` atau kombinasi nilai-nilai lain.