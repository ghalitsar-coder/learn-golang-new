# Struct Methods dan Receiver dalam Go

Dalam Go, tidak ada konsep class seperti dalam bahasa pemrograman berorientasi objek tradisional. Sebagai gantinya, Go menggunakan struct dan method dengan receiver untuk mencapai enkapsulasi dan perilaku yang mirip dengan class method di bahasa lain.

## 1. Dasar Struct dan Method

### Struct
Struct adalah tipe data komposit yang memungkinkan kita menggabungkan beberapa field dalam satu kesatuan:

```go
package main

import "fmt"

// Person adalah struct yang merepresentasikan seseorang
type Person struct {
    Name string
    Age  int
    Email string
}
```

### Method dengan Receiver Value
Method dengan receiver value membuat copy dari struct:

```go
// Greet adalah method dengan receiver value
func (p Person) Greet() string {
    return fmt.Sprintf("Hello, my name is %s", p.Name)
}
```

### Method dengan Receiver Pointer
Method dengan receiver pointer menerima pointer ke struct asli:

```go
// HaveBirthday adalah method dengan receiver pointer
func (p *Person) HaveBirthday() {
    p.Age++
}

// SetEmail adalah method dengan receiver pointer untuk mengubah field
func (p *Person) SetEmail(email string) {
    p.Email = email
}

func main() {
    // Membuat instance struct
    person := Person{
        Name: "Alice",
        Age:  30,
        Email: "alice@example.com",
    }
    
    // Memanggil method dengan receiver value
    fmt.Println(person.Greet()) // Output: Hello, my name is Alice
    
    // Memanggil method dengan receiver pointer
    person.HaveBirthday()
    fmt.Printf("After birthday: %s is %d years old\n", person.Name, person.Age)
    // Output: After birthday: Alice is 31 years old
    
    // Mengubah email
    person.SetEmail("alice.new@example.com")
    fmt.Printf("New email: %s\n", person.Email)
    // Output: New email: alice.new@example.com
}
```

## 2. Perbedaan Receiver Value dan Receiver Pointer

Pemilihan antara receiver value dan receiver pointer sangat penting dalam Go:

```go
package main

import "fmt"

type Counter struct {
    Value int
}

// IncrementValue menggunakan receiver value
func (c Counter) IncrementValue() {
    c.Value++ // Ini hanya mengubah copy dari struct
}

// IncrementPointer menggunakan receiver pointer
func (c *Counter) IncrementPointer() {
    c.Value++ // Ini mengubah struct asli
}

// GetValue menggunakan receiver value
func (c Counter) GetValue() int {
    return c.Value
}

// GetValuePointer menggunakan receiver pointer
func (c *Counter) GetValuePointer() int {
    return c.Value
}

func main() {
    counter := Counter{Value: 0}
    
    fmt.Println("Initial value:", counter.GetValue()) // Output: Initial value: 0
    
    // Memanggil method dengan receiver value
    counter.IncrementValue()
    fmt.Println("After IncrementValue:", counter.GetValue()) // Output: After IncrementValue: 0 (tidak berubah)
    
    // Memanggil method dengan receiver pointer
    counter.IncrementPointer()
    fmt.Println("After IncrementPointer:", counter.GetValue()) // Output: After IncrementPointer: 1 (berubah)
}
```

## 3. Method pada Struct Kompleks

Contoh dengan struct yang lebih kompleks:

```go
package main

import (
    "fmt"
    "strings"
)

// Address merepresentasikan alamat
type Address struct {
    Street string
    City   string
    Zip    string
}

// FullName mengembalikan nama lengkap
func (a Address) FullName() string {
    return fmt.Sprintf("%s, %s %s", a.Street, a.City, a.Zip)
}

// Person dengan nested struct
type Person struct {
    FirstName string
    LastName  string
    Age       int
    Address   Address
    Hobbies   []string
}

// FullName mengembalikan nama lengkap person
func (p Person) FullName() string {
    return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}

// AddHobby menambahkan hobi ke slice
func (p *Person) AddHobby(hobby string) {
    p.Hobbies = append(p.Hobbies, hobby)
}

// RemoveHobby menghapus hobi dari slice
func (p *Person) RemoveHobby(hobby string) {
    for i, h := range p.Hobbies {
        if strings.ToLower(h) == strings.ToLower(hobby) {
            p.Hobbies = append(p.Hobbies[:i], p.Hobbies[i+1:]...)
            return
        }
    }
}

// HasHobby memeriksa apakah person memiliki hobi tertentu
func (p Person) HasHobby(hobby string) bool {
    for _, h := range p.Hobbies {
        if strings.ToLower(h) == strings.ToLower(hobby) {
            return true
        }
    }
    return false
}

// GetHobbies mengembalikan copy dari slice hobbies
func (p Person) GetHobbies() []string {
    // Membuat copy untuk mencegah modifikasi langsung
    hobbies := make([]string, len(p.Hobbies))
    copy(hobbies, p.Hobbies)
    return hobbies
}

// UpdateAddress memperbarui alamat
func (p *Person) UpdateAddress(street, city, zip string) {
    p.Address = Address{
        Street: street,
        City:   city,
        Zip:    zip,
    }
}

func main() {
    person := Person{
        FirstName: "John",
        LastName:  "Doe",
        Age:       25,
        Address: Address{
            Street: "123 Main St",
            City:   "New York",
            Zip:    "10001",
        },
        Hobbies: []string{"reading", "swimming"},
    }
    
    // Menggunakan method-method
    fmt.Println("Full name:", person.FullName()) // Output: Full name: John Doe
    fmt.Println("Address:", person.Address.FullName()) // Output: Address: 123 Main St, New York 10001
    
    // Menambahkan hobi
    person.AddHobby("coding")
    fmt.Println("Hobbies:", person.GetHobbies()) // Output: Hobbies: [reading swimming coding]
    
    // Memeriksa hobi
    fmt.Println("Has 'reading' hobby:", person.HasHobby("reading")) // Output: Has 'reading' hobby: true
    fmt.Println("Has 'dancing' hobby:", person.HasHobby("dancing")) // Output: Has 'dancing' hobby: false
    
    // Menghapus hobi
    person.RemoveHobby("swimming")
    fmt.Println("Hobbies after removing swimming:", person.GetHobbies()) 
    // Output: Hobbies after removing swimming: [reading coding]
    
    // Memperbarui alamat
    person.UpdateAddress("456 Oak Ave", "Los Angeles", "90210")
    fmt.Println("New address:", person.Address.FullName()) 
    // Output: New address: 456 Oak Ave, Los Angeles 90210
}
```

## 4. Method dengan Parameter dan Return Value

Method juga bisa menerima parameter dan mengembalikan nilai:

```go
package main

import "fmt"

// Rectangle merepresentasikan persegi panjang
type Rectangle struct {
    Width  float64
    Height float64
}

// Area menghitung luas persegi panjang
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Perimeter menghitung keliling persegi panjang
func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// IsSquare memeriksa apakah persegi panjang adalah persegi
func (r Rectangle) IsSquare() bool {
    return r.Width == r.Height
}

// Scale mengubah ukuran persegi panjang
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

// Contains memeriksa apakah titik (x,y) berada dalam persegi panjang
func (r Rectangle) Contains(x, y float64) bool {
    return x >= 0 && x <= r.Width && y >= 0 && y <= r.Height
}

// Intersects memeriksa apakah dua persegi panjang berpotongan
func (r Rectangle) Intersects(other Rectangle) bool {
    // Ini adalah implementasi sederhana, asumsi persegi panjang dimulai dari (0,0)
    return r.Width > 0 && r.Height > 0 && 
           other.Width > 0 && other.Height > 0 &&
           r.Width >= other.Width && r.Height >= other.Height
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    
    fmt.Printf("Rectangle: %.1f x %.1f\n", rect.Width, rect.Height)
    fmt.Printf("Area: %.1f\n", rect.Area())           // Output: Area: 50.0
    fmt.Printf("Perimeter: %.1f\n", rect.Perimeter()) // Output: Perimeter: 30.0
    fmt.Printf("Is square: %t\n", rect.IsSquare())    // Output: Is square: false
    
    // Mengubah ukuran
    rect.Scale(2)
    fmt.Printf("After scaling by 2: %.1f x %.1f\n", rect.Width, rect.Height)
    fmt.Printf("New area: %.1f\n", rect.Area())       // Output: New area: 200.0
    
    // Memeriksa titik
    fmt.Printf("Contains point (5, 3): %t\n", rect.Contains(5, 3)) // Output: Contains point (5, 3): true
    fmt.Printf("Contains point (15, 10): %t\n", rect.Contains(15, 10)) // Output: Contains point (15, 10): false
    
    // Memeriksa interseksi
    otherRect := Rectangle{Width: 5, Height: 3}
    fmt.Printf("Intersects with 5x3 rectangle: %t\n", rect.Intersects(otherRect)) 
    // Output: Intersects with 5x3 rectangle: true
}
```

## 5. Constructor dan Factory Methods

Dalam Go, tidak ada constructor eksplisit seperti dalam bahasa lain, tetapi kita bisa membuat fungsi factory:

```go
package main

import (
    "fmt"
    "time"
)

// User merepresentasikan pengguna
type User struct {
    ID        int
    Username  string
    Email     string
    CreatedAt time.Time
    IsActive  bool
}

// NewUser adalah factory function untuk membuat User baru
func NewUser(username, email string) *User {
    return &User{
        Username:  username,
        Email:     email,
        CreatedAt: time.Now(),
        IsActive:  true,
    }
}

// NewUserWithID adalah factory function untuk membuat User dengan ID spesifik
func NewUserWithID(id int, username, email string) *User {
    return &User{
        ID:        id,
        Username:  username,
        Email:     email,
        CreatedAt: time.Now(),
        IsActive:  true,
    }
}

// Activate mengaktifkan user
func (u *User) Activate() {
    u.IsActive = true
}

// Deactivate menonaktifkan user
func (u *User) Deactivate() {
    u.IsActive = false
}

// UpdateEmail memperbarui email user
func (u *User) UpdateEmail(email string) error {
    if email == "" {
        return fmt.Errorf("email cannot be empty")
    }
    u.Email = email
    return nil
}

// GetUserInfo mengembalikan informasi user dalam format string
func (u User) GetUserInfo() string {
    status := "active"
    if !u.IsActive {
        status = "inactive"
    }
    return fmt.Sprintf("User: %s (%s) - %s - Created: %s", 
        u.Username, u.Email, status, u.CreatedAt.Format("2006-01-02"))
}

func main() {
    // Membuat user dengan factory function
    user1 := NewUser("alice", "alice@example.com")
    fmt.Println(user1.GetUserInfo())
    
    user2 := NewUserWithID(100, "bob", "bob@example.com")
    fmt.Println(user2.GetUserInfo())
    
    // Mengubah status user
    user1.Deactivate()
    fmt.Printf("After deactivation: %s\n", user1.GetUserInfo())
    
    // Memperbarui email
    if err := user2.UpdateEmail("bob.new@example.com"); err != nil {
        fmt.Printf("Error updating email: %v\n", err)
    } else {
        fmt.Printf("After email update: %s\n", user2.GetUserInfo())
    }
}
```

## 6. Method Chaining

Kita juga bisa mengimplementasikan method chaining seperti dalam JavaScript:

```go
package main

import (
    "fmt"
    "strings"
)

// StringBuilder adalah struct untuk membangun string
type StringBuilder struct {
    value string
}

// NewStringBuilder adalah factory function untuk StringBuilder
func NewStringBuilder() *StringBuilder {
    return &StringBuilder{}
}

// Append menambahkan string dan mengembalikan pointer ke struct untuk chaining
func (sb *StringBuilder) Append(s string) *StringBuilder {
    sb.value += s
    return sb
}

// AppendLine menambahkan string diikuti newline
func (sb *StringBuilder) AppendLine(s string) *StringBuilder {
    sb.value += s + "\n"
    return sb
}

// Repeat menambahkan string diulang n kali
func (sb *StringBuilder) Repeat(s string, n int) *StringBuilder {
    for i := 0; i < n; i++ {
        sb.value += s
    }
    return sb
}

// Upper mengubah semua karakter menjadi huruf kapital
func (sb *StringBuilder) Upper() *StringBuilder {
    sb.value = strings.ToUpper(sb.value)
    return sb
}

// Lower mengubah semua karakter menjadi huruf kecil
func (sb *StringBuilder) Lower() *StringBuilder {
    sb.value = strings.ToLower(sb.value)
    return sb
}

// String mengembalikan nilai string
func (sb StringBuilder) String() string {
    return sb.value
}

// Reset mengosongkan nilai string
func (sb *StringBuilder) Reset() *StringBuilder {
    sb.value = ""
    return sb
}

func main() {
    // Menggunakan method chaining
    builder := NewStringBuilder()
    
    result := builder.
        Append("Hello").
        Append(" ").
        Append("World").
        AppendLine("!").
        Repeat("-", 10).
        AppendLine("").
        Upper().
        String()
    
    fmt.Print(result)
    // Output:
    // HELLO WORLD!
    // ----------
    
    // Menggunakan kembali dengan Reset
    result2 := builder.
        Reset().
        AppendLine("New text").
        Lower().
        String()
    
    fmt.Print(result2)
    // Output:
    // new text
}
```

## 7. Interface dan Method

Method dalam Go juga memungkinkan kita mengimplementasikan konsep interface:

```go
package main

import (
    "fmt"
    "math"
)

// Shape adalah interface yang mendefinisikan method-method untuk bentuk geometris
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Rectangle merepresentasikan persegi panjang
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Circle merepresentasikan lingkaran
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Calculate menghitung area dan perimeter dari bentuk
func Calculate(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 3}
    
    fmt.Println("Rectangle:")
    Calculate(rect) // Output: Area: 50.00, Perimeter: 30.00
    
    fmt.Println("Circle:")
    Calculate(circle) // Output: Area: 28.27, Perimeter: 18.85
    
    // Kita juga bisa menggunakan slice dari interface
    shapes := []Shape{rect, circle}
    
    for i, shape := range shapes {
        fmt.Printf("Shape %d:\n", i+1)
        Calculate(shape)
    }
}
```

## Kesimpulan

Dalam Go, method dengan receiver adalah cara utama untuk menambahkan perilaku ke struct. Perbedaan antara receiver value dan receiver pointer penting untuk dipahami:

1. **Receiver Value**:
   - Membuat copy dari struct
   - Tidak dapat mengubah field struct asli
   - Cocok untuk method yang hanya membaca data

2. **Receiver Pointer**:
   - Menerima pointer ke struct asli
   - Dapat mengubah field struct
   - Lebih efisien karena tidak perlu menyalin struct
   - Cocok untuk method yang mengubah data

Dengan menggunakan method dan receiver, kita bisa membuat API yang mirip dengan class-based OOP, meskipun dengan pendekatan yang berbeda dan seringkali lebih sederhana.