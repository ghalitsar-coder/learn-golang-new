package main

import "fmt"

// =============================================================================
// DEMO: Perbedaan Struct Literal vs new()
// =============================================================================

type Person struct {
	Name  string
	Age   int
	Email string
}

func main() {
	fmt.Println("========== CARA 1: Literal (Menghasilkan VALUE) ==========")

	// Literal: menghasilkan struct VALUE
	person1 := Person{Name: "Alice", Age: 25, Email: "alice@example.com"}

	fmt.Printf("Type: %T\n", person1)       // Type: main.Person (VALUE)
	fmt.Printf("Value: %+v\n", person1)     // Value: {Name:Alice Age:25 Email:alice@example.com}
	fmt.Printf("Address: %p\n\n", &person1) // Address: 0x... (ambil addressnya pakai &)

	// Akses field langsung
	fmt.Println("Name:", person1.Name)
	person1.Age = 26
	fmt.Println("Age after update:", person1.Age)

	fmt.Println("\n========== CARA 2: new() (Menghasilkan POINTER) ==========")

	// new(): menghasilkan POINTER ke struct (zero value untuk semua field)
	person2 := new(Person)

	fmt.Printf("Type: %T\n", person2)      // Type: *main.Person (POINTER)
	fmt.Printf("Value: %+v\n", person2)    // Value: &{Name: Age:0 Email:} (zero values)
	fmt.Printf("Address: %p\n\n", person2) // Address: 0x... (sudah pointer, tidak perlu &)

	// Set field (Go otomatis dereference, tidak perlu (*person2).Name)
	person2.Name = "Bob"
	person2.Age = 30
	person2.Email = "bob@example.com"

	fmt.Printf("After setting values: %+v\n", person2)

	fmt.Println("\n========== CARA 3: Literal POINTER (& di depan literal) ==========")

	// Cara alternatif: literal dengan & untuk dapat pointer
	person3 := &Person{Name: "Charlie", Age: 35, Email: "charlie@example.com"}

	fmt.Printf("Type: %T\n", person3)      // Type: *main.Person (POINTER)
	fmt.Printf("Value: %+v\n", person3)    // Value: &{Name:Charlie Age:35 Email:charlie@example.com}
	fmt.Printf("Address: %p\n\n", person3) // Address: 0x...

	fmt.Println("\n========== PERBANDINGAN DALAM FUNGSI ==========")

	// Test dengan value
	modifyPersonValue(person1)
	fmt.Println("After modifyPersonValue, person1.Age:", person1.Age) // Tetap 26 (tidak berubah)

	// Test dengan pointer (person2 sudah pointer)
	modifyPersonPointer(person2)
	fmt.Println("After modifyPersonPointer, person2.Age:", person2.Age) // Berubah jadi 100

	// Test dengan pointer (person3 juga pointer)
	modifyPersonPointer(person3)
	fmt.Println("After modifyPersonPointer, person3.Age:", person3.Age) // Berubah jadi 100

	// Kalau mau person1 (value) bisa diubah, kirim addressnya
	modifyPersonPointer(&person1)
	fmt.Println("After modifyPersonPointer(&person1), person1.Age:", person1.Age) // Berubah jadi 100

	fmt.Println("\n========== KAPAN PAKAI APA? ==========")
	printUsageGuidelines()
}

// Fungsi dengan parameter value (terima copy)
func modifyPersonValue(p Person) {
	p.Age = 100
	fmt.Printf("  [dalam fungsi value] Age = %d\n", p.Age)
}

// Fungsi dengan parameter pointer (terima reference)
func modifyPersonPointer(p *Person) {
	p.Age = 100
	fmt.Printf("  [dalam fungsi pointer] Age = %d\n", p.Age)
}

func printUsageGuidelines() {
	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────┐
│ CARA 1: LITERAL (VALUE)                                             │
│ person := Person{Name: "Alice", Age: 25, Email: "..."}              │
├─────────────────────────────────────────────────────────────────────┤
│ ✅ Struct kecil (beberapa field sederhana)                          │
│ ✅ Tidak perlu modify dari fungsi lain                              │
│ ✅ Immutability (functional programming style)                      │
│ ✅ Temporary objects (tidak perlu bertahan lama)                    │
│ ❌ Struct besar (banyak field) - overhead copy                      │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│ CARA 2: new()                                                       │
│ person := new(Person) // Zero values: "", 0, ""                     │
├─────────────────────────────────────────────────────────────────────┤
│ ✅ Butuh pointer tapi inisialisasi bertahap                         │
│ ✅ Struct dengan zero values yang valid                             │
│ ⚠️  Jarang dipakai (lebih sering pakai cara 3)                      │
│ ❌ Kalau langsung tahu nilai field-nya (pakai cara 3)               │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│ CARA 3: LITERAL POINTER (PALING UMUM untuk pointer)                │
│ person := &Person{Name: "Alice", Age: 25, Email: "..."}             │
├─────────────────────────────────────────────────────────────────────┤
│ ✅ Butuh pointer DAN langsung tahu nilai field                      │
│ ✅ Struct besar (efisien pass ke fungsi)                            │
│ ✅ Method perlu modify struct (pointer receiver)                    │
│ ✅ Perlu share state antar fungsi                                   │
│ ✅ Struct sebagai part dari data structure (linked list, tree)      │
│ ✅ Bisa return nil (pointer bisa nil, value tidak)                  │
└─────────────────────────────────────────────────────────────────────┘
`)
}
