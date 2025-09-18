# Bagian 2: Struktur Data dan Interface

## 1. Array dan Slice

### a. Array
Array adalah kumpulan elemen dengan tipe data yang sama dan ukuran tetap. Ukuran array merupakan bagian dari tipenya, sehingga `[5]int` dan `[10]int` adalah tipe yang berbeda.

Deklarasi dan inisialisasi:
```go
// Deklarasi array dengan panjang 5, elemen default adalah 0
var numbers [5]int

// Inisialisasi langsung
var fruits [3]string = [3]string{"apple", "banana", "orange"}

// Inisialisasi dengan ... untuk menentukan panjang otomatis
animals := [...]string{"cat", "dog", "bird"}

// Mengakses elemen
fmt.Println(fruits[0]) // Output: apple

// Mengubah elemen
numbers[0] = 10
fmt.Println(numbers[0]) // Output: 10

// Panjang array
fmt.Println(len(fruits)) // Output: 3
```

Iterasi array:
```go
for i := 0; i < len(fruits); i++ {
	fmt.Printf("Fruit %d: %s\n", i, fruits[i])
}

// atau dengan range
for index, value := range animals {
	fmt.Printf("Animal %d: %s\n", index, value)
}
```

### b. Slice
Slice adalah struktur data yang lebih dinamis dan fleksibel dibanding array. Slice merupakan reference ke bagian dari array. Slice tidak menyimpan data itu sendiri, melainkan mereferensikan elemen-elemen dari array lain.

Membuat slice:
```go
// Membuat slice dengan make
slice1 := make([]int, 5) // Membuat slice dengan panjang 5, kapasitas 5

// Membuat slice dengan panjang dan kapasitas spesifik
slice2 := make([]int, 3, 5) // Panjang 3, kapasitas 5

// Membuat slice dari array
arr := [5]int{1, 2, 3, 4, 5}
slice3 := arr[1:4] // Elemen index 1 sampai 3 (4 tidak termasuk)

// Membuat slice literal
slice4 := []int{10, 20, 30}
```

Operasi pada slice:
```go
// Menambah elemen dengan append
slice4 = append(slice4, 40)
slice4 = append(slice4, 50, 60) // Bisa menambahkan beberapa elemen sekaligus

// Panjang dan kapasitas
fmt.Println("Length:", len(slice4))   // Output: 5
fmt.Println("Capacity:", cap(slice4)) // Output: 6 (kapasitas awal slice literal biasanya dua kali panjangnya)

// Copy slice
src := []int{1, 2, 3}
dst := make([]int, len(src))
copy(dst, src)
fmt.Println(dst) // Output: [1 2 3]
```

Perbedaan penting antara Array dan Slice:
- Array memiliki ukuran tetap, sedangkan slice bisa berubah ukuran.
- Saat mengirim array ke fungsi, itu adalah *pass by value* (copy seluruh array), sedangkan slice adalah *pass by reference* (hanya copy header slice-nya).

Contoh:
```go
func modifyArray(arr [3]int) {
	arr[0] = 100 // Tidak mempengaruhi array asli
}

func modifySlice(slc []int) {
	slc[0] = 100 // Memodifikasi slice asli karena reference
}

func main() {
	myArr := [3]int{1, 2, 3}
	mySlice := []int{1, 2, 3}
	
	modifyArray(myArr)
	fmt.Println(myArr) // Output: [1 2 3] (tidak berubah)
	
	modifySlice(mySlice)
	fmt.Println(mySlice) // Output: [100 2 3] (berubah)
}
```

## 2. Map

Map adalah koleksi key-value yang tidak terurut. Key harus memiliki tipe yang *comparable* (seperti string, int, bool), dan value bisa bertipe apa saja.

Deklarasi dan inisialisasi:
```go
// Membuat map kosong
var ages map[string]int // Nil map, akan panic jika diakses sebelum dibuat dengan make

// Membuat map dengan make
ages = make(map[string]int)
ages["Alice"] = 25
ages["Bob"] = 30

// Inisialisasi literal
scores := map[string]int{
	"math": 90,
	"english": 85,
	"science": 88,
}

// Mengakses elemen
fmt.Println(ages["Alice"]) // Output: 25

// Memeriksa keberadaan key
if score, exists := scores["math"]; exists {
	fmt.Printf("Math score: %d\n", score)
} else {
	fmt.Println("Math score not found")
}

// Menghapus elemen
delete(ages, "Bob")
fmt.Println(ages) // Output: map[Alice:25]
```

Iterasi map:
```go
for subject, score := range scores {
	fmt.Printf("%s: %d\n", subject, score)
}
```

## 3. Struct

Struct adalah tipe data komposit yang memungkinkan kita menggabungkan beberapa field (variabel) dalam satu kesatuan. Ini mirip dengan konsep "class" di bahasa pemrograman lain, namun Go tidak memiliki konsep class secara eksplisit.

Deklarasi struct:
```go
type Person struct {
	Name string
	Age  int
	Email string
}

type Rectangle struct {
	Width  float64
	Height float64
}
```

Membuat instance struct:
```go
// Cara 1: Menggunakan literal
person1 := Person{Name: "Alice", Age: 25, Email: "alice@example.com"}

// Cara 2: Menggunakan new (menghasilkan pointer)
person2 := new(Person)
person2.Name = "Bob"
person2.Age = 30
person2.Email = "bob@example.com"

// Cara 3: Literal tanpa nama field (urutan harus sesuai)
person3 := Person{"Charlie", 35, "charlie@example.com"}

// Mengakses field
fmt.Println(person1.Name) // Output: Alice
person1.Age = 26 // Mengubah field
fmt.Println(person1.Age) // Output: 26
```

Nested struct:
```go
type Address struct {
	Street string
	City   string
}

type Employee struct {
	PersonInfo Person
	AddressInfo Address
	Position string
}

func main() {
	emp := Employee{
		PersonInfo: Person{Name: "David", Age: 28, Email: "david@example.com"},
		AddressInfo: Address{Street: "123 Main St", City: "Gotham"},
		Position: "Developer",
	}
	
	fmt.Printf("Employee: %+v\n", emp)
	fmt.Println("City:", emp.AddressInfo.City)
}
```

## 4. Method

Method adalah fungsi yang memiliki *receiver*. Receiver adalah argumen yang mendefinisikan objek mana yang akan memanggil method tersebut. Ini memungkinkan kita untuk menetapkan perilaku ke struct (mirip dengan metode dalam class OOP).

Mendefinisikan method:
```go
// Method dengan receiver value
func (p Person) Greet() string {
	return fmt.Sprintf("Hi, I'm %s", p.Name)
}

// Method dengan receiver pointer (lebih umum untuk mengubah data)
func (p *Person) SetAge(age int) {
	p.Age = age
}

// Method untuk struct Rectangle
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

func main() {
	person := Person{Name: "Eve", Age: 22, Email: "eve@example.com"}
	fmt.Println(person.Greet()) // Memanggil method
	
	// Memanggil method dengan receiver pointer
	person.SetAge(23)
	fmt.Println("New age:", person.Age)
	
	rect := Rectangle{Width: 10, Height: 5}
	fmt.Println("Area:", rect.Area()) // Output: 50
	
	// Scale rectangle
	rect.Scale(2)
	fmt.Println("Scaled area:", rect.Area()) // Output: 200
}
```

Perbedaan antara receiver value dan pointer:
- Receiver value (`func (p Person) ...`): Membuat *copy* dari struct saat method dipanggil. Perubahan dalam method tidak mempengaruhi struct asli.
- Receiver pointer (`func (p *Person) ...`): Menerima *pointer* ke struct. Perubahan dalam method mempengaruhi struct asli. Lebih efisien karena tidak perlu menyalin struct.

## 5. Interface

Interface adalah kumpulan definisi method (tanpa implementasi). Tipe apa pun yang memiliki semua method yang didefinisikan dalam interface secara otomatis mengimplementasikannya. Ini adalah cara Go menerapkan konsep *polymorphism*.

Mendefinisikan interface:
```go
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Struct Rectangle (sudah didefinisikan sebelumnya)
// Kita akan menambahkan method Perimeter
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Struct Circle
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}
```

Menggunakan interface:
```go
import "math"

// Fungsi yang menerima interface
func PrintShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f\n", s.Area())
	fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
}

func main() {
	rect := Rectangle{Width: 10, Height: 5}
	circ := Circle{Radius: 3}
	
	// Karena Rectangle dan Circle mengimplementasikan Shape,
	// mereka bisa dikirim ke fungsi yang menerima Shape
	PrintShapeInfo(rect)
	PrintShapeInfo(circ)
}
```

Interface kosong (`interface{}`):
Interface kosong adalah interface tanpa method apapun. Semua tipe mengimplementasikan interface kosong secara otomatis. Ini digunakan ketika kita ingin membuat fungsi yang bisa menerima parameter dengan tipe apapun.

```go
func PrintAnything(v interface{}) {
	fmt.Println(v)
}

func main() {
	PrintAnything("Hello")
	PrintAnything(42)
	PrintAnything(true)
	PrintAnything(Rectangle{Width: 1, Height: 2})
}
```

Namun, sejak Go 1.18, konsep *generics* telah diperkenalkan yang seringkali lebih disarankan daripada `interface{}` untuk fleksibilitas tipe yang aman.