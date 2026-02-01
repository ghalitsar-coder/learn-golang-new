package main

import (
	"fmt"
	"math"
)

type Person struct {
	Name  string
	Age   int
	Email string
}

type Rectangle struct {
	Width  float64
	Height float64
}

type Shape interface {
	GetArea() float64
	Perimeter() float64
}
type Circle struct {
	Radius float64
}

func (c Circle) GetArea() float64 {
	return math.Pi * c.Radius * c.Radius
}
func (r Rectangle) GetArea() float64 {
	return r.Height * r.Width
	//Output:
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func main() {

	Method_4()
}

func PrintShapeInfo(s Shape) {

	fmt.Printf("Area dari shape : %.2f\n", s.GetArea())
	fmt.Printf("Perimter dari shape : %.2f\n", s.Perimeter())

}

func Interface_5() {

	rect := Rectangle{10, 5}
	circle := Circle{5}

	PrintShapeInfo(rect)
	PrintShapeInfo(circle)
	//Output:

}

func (p Person) Greet() {

	//Output:
	fmt.Printf("Hello %v\n", p.Name)

}

func (p *Person) SetAge(age int) {
	p.Age = age

}

func (p *Person) IncrementAge() {
	p.Age++
	fmt.Printf("You grow older by one year")
	//Output:
}



func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
	//Output:

}

func Method_4() {

	person1 := Person{"adam", 22, "adam@.com"}
	fmt.Printf("person1 : %+v\n ", person1)
	person1.Greet()
	person1.SetAge(55)
	fmt.Printf("after set age person1 : %+v\n ", person1)

	rec1 := Rectangle{Width: 10, Height: 5}
	rec1.GetArea()
	fmt.Printf("RECT AREA : %+v\n ", rec1.GetArea())

	rec1.Scale(2)

	fmt.Printf("after scale RECT AREA : %+v\n", rec1.GetArea())

}

func Struct_3_2_nested() {
	type structName struct {
	}

	type Address struct {
		Street string
		City   string
	}

	type Employee struct {
		PersonInfo  Person
		AddressInfo Address
		Position    string
	}

	emp := Employee{
		PersonInfo:  Person{"Adam", 22, "something@a.com"},
		AddressInfo: Address{"jl.cigadung", "bandung"},
		Position:    "Jobless",
	}

	fmt.Printf("Employees : %+v\n", emp)

}

func Struct_3() {

	mahasiswa := Person{Name: "ghal", Age: 24, Email: "ghal@a.com"}
	mahasiswa_q := &Person{Name: "ghal", Age: 24, Email: "ghal@a.com"}

	mahasiswa2 := Person{"adam", 22, "adam@a.com"}
	// fmt.Printf("Mahasiswa 1 %v\n", mahasiswa)
	// fmt.Printf("Mahasiswa 2 %v\n", mahasiswa2)
	mahasiswa3 := mahasiswa
	mahasiswa4 := mahasiswa2
	mahasiswa4.Age = 89
	mahasiswa3.Age = 45

	fmt.Printf("Mahasiswa 1 %v\n", mahasiswa)
	fmt.Printf("Mahasiswa 2 %v\n", mahasiswa2)
	fmt.Printf("Mahasiswa 3 %v\n", mahasiswa3)
	fmt.Printf("Mahasiswa 4 %v\n", mahasiswa4)
	fmt.Printf("Mahasiswa Q %v\n", mahasiswa_q)

}

func Struct_PointerVsValue() {
	type Person struct {
		Name  string
		Age   int
		Email string
	}

	fmt.Println("=== 1. PERBEDAAN DASAR ===")
	// VALUE: membuat salinan
	valueStruct := Person{Name: "Budi", Age: 25, Email: "budi@mail.com"}

	// POINTER: menyimpan alamat memori (referensi)
	pointerStruct := &Person{Name: "Siti", Age: 23, Email: "siti@mail.com"}

	fmt.Printf("valueStruct: %v (tipe: %T)\n", valueStruct, valueStruct)
	fmt.Printf("pointerStruct: %v (tipe: %T)\n", pointerStruct, pointerStruct)
	fmt.Printf("Alamat memori pointerStruct: %p\n\n", pointerStruct)

	// ===== COPY vs REFERENCE =====
	fmt.Println("=== 2. COPY vs REFERENCE ===")

	// VALUE: membuat salinan independen
	person1 := Person{Name: "Ali", Age: 30, Email: "ali@mail.com"}
	person2 := person1 // COPY - salinan baru
	person2.Age = 35

	fmt.Printf("person1: %v (Age: %d)\n", person1, person1.Age)   // Age tetap 30
	fmt.Printf("person2: %v (Age: %d)\n\n", person2, person2.Age) // Age berubah 35

	// POINTER: berbagi data yang sama
	person3 := &Person{Name: "Rina", Age: 28, Email: "rina@mail.com"}
	person4 := person3 // REFERENCE - menunjuk ke data yang sama
	person4.Age = 40

	fmt.Printf("person3: %v (Age: %d)\n", person3, person3.Age)   // Age ikut berubah 40!
	fmt.Printf("person4: %v (Age: %d)\n\n", person4, person4.Age) // Age 40

	// ===== FUNCTION PARAMETER =====
	fmt.Println("=== 3. FUNCTION PARAMETER ===")

	// Function dengan VALUE - tidak mengubah original
	updateAgeByValue := func(p Person, newAge int) {
		p.Age = newAge
		fmt.Printf("  Di dalam func (value): %v\n", p)
	}

	// Function dengan POINTER - mengubah original
	updateAgeByPointer := func(p *Person, newAge int) {
		p.Age = newAge // Go otomatis dereference
		fmt.Printf("  Di dalam func (pointer): %v\n", p)
	}

	testPerson := Person{Name: "Doni", Age: 22, Email: "doni@mail.com"}
	fmt.Printf("Sebelum: %v\n", testPerson)

	updateAgeByValue(testPerson, 99)
	fmt.Printf("Setelah updateAgeByValue: %v (TIDAK BERUBAH)\n\n", testPerson)

	updateAgeByPointer(&testPerson, 99)
	fmt.Printf("Setelah updateAgeByPointer: %v (BERUBAH!)\n\n", testPerson)

	// ===== AKSES FIELD =====
	fmt.Println("=== 4. AKSES FIELD (SAMA!) ===")
	val := Person{Name: "Andi", Age: 26, Email: "andi@mail.com"}
	ptr := &Person{Name: "Lia", Age: 24, Email: "lia@mail.com"}

	// Keduanya bisa diakses dengan cara yang sama
	fmt.Printf("Value access: %s, %d\n", val.Name, val.Age)
	fmt.Printf("Pointer access: %s, %d\n", ptr.Name, ptr.Age) // Go auto-dereference

	// Ubah field
	val.Age = 27
	ptr.Age = 25

	fmt.Printf("After update - val: %v\n", val)
	fmt.Printf("After update - ptr: %v\n\n", ptr)

	// ===== KAPAN PAKAI POINTER? =====
	fmt.Println("=== 5. KAPAN PAKAI POINTER? ===")

	type LargeStruct struct {
		Data [1000]int
		Name string
	}

	type SmallStruct struct {
		ID   int
		Name string
	}

	fmt.Println("GUNAKAN POINTER jika:")
	fmt.Println("✓ Struct besar (hemat memory, tidak copy)")
	large := &LargeStruct{Name: "Big Data"}
	fmt.Printf("  Large struct pointer: %p\n", large)

	fmt.Println("✓ Perlu memodifikasi nilai original")
	modifiable := &Person{Name: "John", Age: 30, Email: "john@mail.com"}
	modifiable.Age = 31 // langsung ubah original
	fmt.Printf("  Modified: %v\n", modifiable)

	fmt.Println("✓ Sharing state antar functions")
	sharedState := &Person{Name: "Shared", Age: 20, Email: "shared@mail.com"}
	func1 := func(p *Person) { p.Age++ }
	func2 := func(p *Person) { p.Age++ }
	func1(sharedState)
	func2(sharedState)
	fmt.Printf("  Shared state age: %d (kedua func berbagi state)\n", sharedState.Age)

	fmt.Println("\nGUNAKAN VALUE jika:")
	fmt.Println("✓ Struct kecil (lebih cepat di stack)")
	small := SmallStruct{ID: 1, Name: "Small"}
	fmt.Printf("  Small struct: %v\n", small)

	fmt.Println("✓ Ingin immutability (tidak mau diubah)")
	immutable := Person{Name: "Fixed", Age: 25, Email: "fixed@mail.com"}
	passByValue := func(p Person) { p.Age = 99 } // tidak akan ubah original
	passByValue(immutable)
	fmt.Printf("  Immutable tetap: %v\n", immutable)

	fmt.Println("✓ Data temporary/lokal scope")

	// ===== NIL POINTER =====
	fmt.Println("\n=== 6. HATI-HATI: NIL POINTER ===")
	var nilPerson *Person // pointer tanpa inisialisasi = nil
	fmt.Printf("nilPerson: %v\n", nilPerson)

	// BAHAYA: akses nil pointer = panic!
	// nilPerson.Name = "Test" // PANIC! runtime error

	// Solusi: cek nil dulu
	if nilPerson != nil {
		nilPerson.Name = "Safe"
	} else {
		fmt.Println("Pointer is nil, initializing...")
		nilPerson = &Person{Name: "New", Age: 20, Email: "new@mail.com"}
	}
	fmt.Printf("After nil check: %v\n", nilPerson)

	// ===== PASSING TO FUNCTIONS: SUMMARY =====
	fmt.Println("\n=== 7. SUMMARY: POINTER vs VALUE ===")

	type Counter struct {
		Count int
	}

	// Function dengan value parameter - tidak mengubah original
	incrementByValue := func(c Counter) {
		c.Count++
		fmt.Printf("  Inside incrementByValue: %d\n", c.Count)
	}

	// Function dengan pointer parameter - mengubah original
	incrementByPointer := func(c *Counter) {
		c.Count++
		fmt.Printf("  Inside incrementByPointer: %d\n", c.Count)
	}

	counter := Counter{Count: 0}
	fmt.Printf("Initial: %d\n", counter.Count)

	incrementByValue(counter)
	fmt.Printf("After incrementByValue: %d (TIDAK BERUBAH)\n", counter.Count)

	incrementByPointer(&counter)
	fmt.Printf("After incrementByPointer: %d (BERUBAH!)\n", counter.Count)

	// Dengan pointer variable
	counterPtr := &Counter{Count: 10}
	incrementByValue(*counterPtr) // dereference dulu
	fmt.Printf("After incrementByValue on ptr: %d (TIDAK BERUBAH)\n", counterPtr.Count)

	incrementByPointer(counterPtr) // langsung pass
	fmt.Printf("After incrementByPointer on ptr: %d (BERUBAH!)\n", counterPtr.Count)
}

func MapMethod_2() {

	person := make(map[string]int)
	person["andre"] = 30
	person["siti"] = 25

	fmt.Printf("isi dari person %v\n", person)

	mahasiswa := map[string]int{
		"ghal":  24,
		"fikar": 23,
		"felix": 21,
	}

	fmt.Printf("isi dari mahasiswa %v\n", mahasiswa)

	if _, exist := mahasiswa["ghal"]; exist {
		fmt.Printf("umur gal %v\n", mahasiswa["ghal"])

	} else {
		fmt.Printf("mahasiswa tidak ada")
	}

}

func Slices_1_1() {
	slice := make([]int, 3, 5)
	fmt.Printf("isi slice %v\n", slice)

	numbers := []int{1, 2, 3}
	additionNumber := []int{4, 5, 6}

	numbers = append(numbers, additionNumber...)

	dst := make([]int, len(numbers), len(numbers)+1)
	copy(dst, numbers)
	// fmt.Printf("isis dari %v\n", dst)

	dst = append(dst, 99)
	dst = append(dst, 92)
	dst = append(dst, 91)
	dst = append(dst, 95)
	// fmt.Printf("isi dari dst new %v", dst)

	changeValue := func(arr []int) {

		arr[0] = 40

	}
	fmt.Printf("isi sebelum diganti %v\n", dst)
	changeValue(dst)
	fmt.Printf("isi sebelum diganti %v \n", dst)

}

// slices addition
func Slices_1() {
	// ===== 1. SLICE PRIMITIF (int) - VALUE COPY =====
	fmt.Println("=== 1. Slice Primitif (int) ===")
	numbers := []int{1, 2, 3, 4, 5}
	additionNumbers := []int{6, 7, 8}
	numbers = append(numbers, additionNumbers...) // spread operator
	fmt.Printf("numbers: %v\n", numbers)

	// Ubah additionNumbers, tidak pengaruhi numbers
	additionNumbers[0] = 999
	fmt.Printf("Setelah ubah additionNumbers[0]=999\n")
	fmt.Printf("additionNumbers: %v\n", additionNumbers)
	fmt.Printf("numbers tetap: %v\n\n", numbers) // Tetap [1 2 3 4 5 6 7 8]

	// ===== 2. SLICE OF SLICES (2 level) - SHALLOW COPY =====
	fmt.Println("=== 2. Slice of Slices (2 level) ===")
	slice2D := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}

	addition2D := [][]int{
		{7, 8, 9},
		{10, 11, 12},
	}

	slice2D = append(slice2D, addition2D...)
	fmt.Printf("slice2D: %v\n", slice2D)

	// HATI-HATI: Ubah addition2D MASIH mempengaruhi slice2D!
	addition2D[0][0] = 777
	fmt.Printf("Setelah ubah addition2D[0][0]=777\n")
	fmt.Printf("addition2D: %v\n", addition2D)
	fmt.Printf("slice2D IKUT berubah: %v\n\n", slice2D) // [7 8 9] jadi [777 8 9]!

	// ===== 3. SLICE OF SLICE OF SLICES (3 level) - MASIH SHALLOW =====
	fmt.Println("=== 3. Slice 3 Level ===")
	slice3D := [][][]int{
		{
			{1, 2},
			{3, 4},
		},
		{
			{5, 6},
			{7, 8},
		},
	}

	addition3D := [][][]int{
		{
			{9, 10},
			{11, 12},
		},
	}

	slice3D = append(slice3D, addition3D...)
	fmt.Printf("slice3D: %v\n", slice3D)

	// Ubah addition3D TETAP mempengaruhi slice3D
	addition3D[0][0][0] = 9999
	fmt.Printf("Setelah ubah addition3D[0][0][0]=9999\n")
	fmt.Printf("addition3D: %v\n", addition3D)
	fmt.Printf("slice3D IKUT berubah: %v\n\n", slice3D)

	// ===== 4. SOLUSI: DEEP COPY MANUAL =====
	fmt.Println("=== 4. Solusi: Deep Copy ===")
	original := [][]int{{1, 2}, {3, 4}}

	// Deep copy manual
	deepCopy := make([][]int, len(original))
	for i := range original {
		deepCopy[i] = make([]int, len(original[i]))
		copy(deepCopy[i], original[i])
	}

	deepCopy = append(deepCopy, []int{5, 6})
	original[0][0] = 999

	fmt.Printf("original (diubah): %v\n", original)
	fmt.Printf("deepCopy (tetap): %v\n", deepCopy)

}
