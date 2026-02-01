package specific

import (
	"fmt"
	"testing"
)

// =============================================================================
// BAGIAN 1: STRUCT UNTUK TESTING
// =============================================================================

// Person adalah struct dasar untuk testing pointer vs value
type Person struct {
	Name string
	Age  int
}

// Counter adalah struct sederhana untuk menghitung
type Counter struct {
	Value int
}

// =============================================================================
// BAGIAN 2: METODE DENGAN VALUE RECEIVER
// =============================================================================

// GetName mengembalikan nama (value receiver - tidak bisa modify)
func (p Person) GetName() string {
	return p.Name
}

// GrowOlderByValue mencoba menambah umur (VALUE RECEIVER - TIDAK AKAN BERPENGARUH!)
func (p Person) GrowOlderByValue() {
	p.Age++ // Ini hanya mengubah COPY, bukan asli!
	fmt.Printf("  [dalam fungsi] Age setelah ++: %d\n", p.Age)
}

// IncrementByValue menambah counter (VALUE RECEIVER - TIDAK AKAN BERPENGARUH!)
func (c Counter) IncrementByValue() {
	c.Value++ // Ini hanya mengubah COPY!
	fmt.Printf("  [dalam fungsi] Value setelah ++: %d\n", c.Value)
}

// =============================================================================
// BAGIAN 3: METODE DENGAN POINTER RECEIVER
// =============================================================================

// GrowOlderByPointer menambah umur (POINTER RECEIVER - AKAN BERPENGARUH!)
func (p *Person) GrowOlderByPointer() {
	p.Age++ // Ini mengubah data ASLI!
	fmt.Printf("  [dalam fungsi] Age setelah ++: %d\n", p.Age)
}

// IncrementByPointer menambah counter (POINTER RECEIVER - AKAN BERPENGARUH!)
func (c *Counter) IncrementByPointer() {
	c.Value++ // Ini mengubah data ASLI!
	fmt.Printf("  [dalam fungsi] Value setelah ++: %d\n", c.Value)
}

// =============================================================================
// BAGIAN 4: FUNGSI DENGAN PARAMETER VALUE VS POINTER
// =============================================================================

// updatePersonByValue menerima COPY dari Person
func updatePersonByValue(p Person, newAge int) {
	p.Age = newAge
	fmt.Printf("  [dalam fungsi] Age: %d\n", p.Age)
}

// updatePersonByPointer menerima POINTER ke Person
func updatePersonByPointer(p *Person, newAge int) {
	p.Age = newAge
	fmt.Printf("  [dalam fungsi] Age: %d\n", p.Age)
}

// =============================================================================
// BAGIAN 5: TEST CASES
// =============================================================================

// TestValueReceiverNotModifyOriginal membuktikan bahwa value receiver TIDAK mengubah data asli
func TestValueReceiverNotModifyOriginal(t *testing.T) {
	fmt.Println("\n========== TEST: Value Receiver TIDAK Mengubah Data Asli ==========")

	person := Person{Name: "Budi", Age: 25}
	fmt.Printf("SEBELUM: %s, Age: %d\n", person.Name, person.Age)

	// Panggil method dengan value receiver
	person.GrowOlderByValue()

	fmt.Printf("SETELAH: %s, Age: %d\n", person.Name, person.Age)

	// Cek apakah umur berubah
	if person.Age != 25 {
		t.Errorf("UNEXPECTED: Age berubah menjadi %d, padahal seharusnya tetap 25", person.Age)
	} else {
		fmt.Println("✅ SESUAI EKSPEKTASI: Age tetap 25 karena value receiver hanya mengubah copy")
	}
}

// TestPointerReceiverModifiesOriginal membuktikan bahwa pointer receiver MENGUBAH data asli
func TestPointerReceiverModifiesOriginal(t *testing.T) {
	fmt.Println("\n========== TEST: Pointer Receiver MENGUBAH Data Asli ==========")

	person := Person{Name: "Ani", Age: 30}
	fmt.Printf("SEBELUM: %s, Age: %d\n", person.Name, person.Age)

	// Panggil method dengan pointer receiver
	person.GrowOlderByPointer()

	fmt.Printf("SETELAH: %s, Age: %d\n", person.Name, person.Age)

	// Cek apakah umur berubah
	if person.Age != 31 {
		t.Errorf("GAGAL: Age seharusnya berubah menjadi 31, tapi masih %d", person.Age)
	} else {
		fmt.Println("✅ SESUAI EKSPEKTASI: Age berubah jadi 31 karena pointer receiver mengubah data asli")
	}
}

// TestCounterValueVsPointer membandingkan perilaku value vs pointer receiver
func TestCounterValueVsPointer(t *testing.T) {
	fmt.Println("\n========== TEST: Perbandingan Counter Value vs Pointer ==========")

	// Test dengan Value Receiver
	fmt.Println("\n--- Value Receiver ---")
	counter1 := Counter{Value: 0}
	fmt.Printf("SEBELUM: Value = %d\n", counter1.Value)

	counter1.IncrementByValue()
	counter1.IncrementByValue()
	counter1.IncrementByValue()

	fmt.Printf("SETELAH 3x increment: Value = %d\n", counter1.Value)

	if counter1.Value != 0 {
		t.Errorf("Value seharusnya tetap 0, tapi %d", counter1.Value)
	} else {
		fmt.Println("✅ Value tetap 0 karena value receiver tidak mengubah asli")
	}

	// Test dengan Pointer Receiver
	fmt.Println("\n--- Pointer Receiver ---")
	counter2 := Counter{Value: 0}
	fmt.Printf("SEBELUM: Value = %d\n", counter2.Value)

	counter2.IncrementByPointer()
	counter2.IncrementByPointer()
	counter2.IncrementByPointer()

	fmt.Printf("SETELAH 3x increment: Value = %d\n", counter2.Value)

	if counter2.Value != 3 {
		t.Errorf("Value seharusnya 3, tapi %d", counter2.Value)
	} else {
		fmt.Println("✅ Value berubah jadi 3 karena pointer receiver mengubah asli")
	}
}

// TestFunctionParameterValueVsPointer membandingkan parameter fungsi value vs pointer
func TestFunctionParameterValueVsPointer(t *testing.T) {
	fmt.Println("\n========== TEST: Parameter Fungsi Value vs Pointer ==========")

	// Test dengan parameter value
	fmt.Println("\n--- Parameter Value ---")
	person1 := Person{Name: "Dina", Age: 20}
	fmt.Printf("SEBELUM: Age = %d\n", person1.Age)

	updatePersonByValue(person1, 100)

	fmt.Printf("SETELAH: Age = %d\n", person1.Age)

	if person1.Age != 20 {
		t.Errorf("Age seharusnya tetap 20, tapi %d", person1.Age)
	} else {
		fmt.Println("✅ Age tetap 20 karena fungsi menerima copy")
	}

	// Test dengan parameter pointer
	fmt.Println("\n--- Parameter Pointer ---")
	person2 := Person{Name: "Erik", Age: 20}
	fmt.Printf("SEBELUM: Age = %d\n", person2.Age)

	updatePersonByPointer(&person2, 100) // Perhatikan &person2

	fmt.Printf("SETELAH: Age = %d\n", person2.Age)

	if person2.Age != 100 {
		t.Errorf("Age seharusnya 100, tapi %d", person2.Age)
	} else {
		fmt.Println("✅ Age berubah jadi 100 karena fungsi menerima pointer")
	}
}

// TestSliceBehavior menunjukkan perilaku unik slice
func TestSliceBehavior(t *testing.T) {
	fmt.Println("\n========== TEST: Perilaku Slice (Special Case!) ==========")

	// Fungsi yang memodifikasi elemen slice
	modifySliceElement := func(s []int) {
		s[0] = 999
		fmt.Printf("  [dalam fungsi] s[0] = %d\n", s[0])
	}

	// Fungsi yang mencoba append ke slice
	appendToSlice := func(s []int) {
		s = append(s, 4)
		fmt.Printf("  [dalam fungsi] slice = %v\n", s)
	}

	// Fungsi yang append dengan pointer
	appendToSliceWithPointer := func(s *[]int) {
		*s = append(*s, 4)
		fmt.Printf("  [dalam fungsi] slice = %v\n", *s)
	}

	// Test 1: Modify elemen - AKAN berubah!
	fmt.Println("\n--- Modify Elemen Slice (tanpa pointer) ---")
	slice1 := []int{1, 2, 3}
	fmt.Printf("SEBELUM: %v\n", slice1)
	modifySliceElement(slice1)
	fmt.Printf("SETELAH: %v\n", slice1)

	if slice1[0] != 999 {
		t.Errorf("Elemen pertama seharusnya 999, tapi %d", slice1[0])
	} else {
		fmt.Println("✅ Elemen berubah karena slice adalah reference type (header di-copy, data sama)")
	}

	// Test 2: Append tanpa pointer - TIDAK akan berubah di luar fungsi
	fmt.Println("\n--- Append ke Slice (tanpa pointer) ---")
	slice2 := []int{1, 2, 3}
	fmt.Printf("SEBELUM: %v (len=%d)\n", slice2, len(slice2))
	appendToSlice(slice2)
	fmt.Printf("SETELAH: %v (len=%d)\n", slice2, len(slice2))

	if len(slice2) != 3 {
		t.Errorf("Length seharusnya tetap 3, tapi %d", len(slice2))
	} else {
		fmt.Println("✅ Slice tetap 3 elemen karena append membuat slice baru di dalam fungsi")
	}

	// Test 3: Append dengan pointer - AKAN berubah
	fmt.Println("\n--- Append ke Slice (dengan pointer) ---")
	slice3 := []int{1, 2, 3}
	fmt.Printf("SEBELUM: %v (len=%d)\n", slice3, len(slice3))
	appendToSliceWithPointer(&slice3)
	fmt.Printf("SETELAH: %v (len=%d)\n", slice3, len(slice3))

	if len(slice3) != 4 {
		t.Errorf("Length seharusnya 4, tapi %d", len(slice3))
	} else {
		fmt.Println("✅ Slice jadi 4 elemen karena kita pakai pointer ke slice")
	}
}

// TestMapBehavior menunjukkan perilaku map
func TestMapBehavior(t *testing.T) {
	fmt.Println("\n========== TEST: Perilaku Map (Reference Type) ==========")

	// Fungsi yang memodifikasi map
	modifyMap := func(m map[string]int) {
		m["score"] = 100
		fmt.Printf("  [dalam fungsi] score = %d\n", m["score"])
	}

	m := map[string]int{"score": 0}
	fmt.Printf("SEBELUM: score = %d\n", m["score"])

	modifyMap(m) // Tanpa pointer!

	fmt.Printf("SETELAH: score = %d\n", m["score"])

	if m["score"] != 100 {
		t.Errorf("Score seharusnya 100, tapi %d", m["score"])
	} else {
		fmt.Println("✅ Score berubah karena map adalah reference type (tidak perlu pointer)")
	}
}

// TestPrimitiveTypeBehavior menunjukkan perilaku tipe primitif
func TestPrimitiveTypeBehavior(t *testing.T) {
	fmt.Println("\n========== TEST: Perilaku Tipe Primitif ==========")

	// Fungsi tanpa pointer
	addTenByValue := func(x int) {
		x += 10
		fmt.Printf("  [dalam fungsi] x = %d\n", x)
	}

	// Fungsi dengan pointer
	addTenByPointer := func(x *int) {
		*x += 10
		fmt.Printf("  [dalam fungsi] x = %d\n", *x)
	}

	// Test value
	fmt.Println("\n--- Primitif tanpa Pointer ---")
	num1 := 5
	fmt.Printf("SEBELUM: %d\n", num1)
	addTenByValue(num1)
	fmt.Printf("SETELAH: %d\n", num1)

	if num1 != 5 {
		t.Errorf("Num seharusnya tetap 5, tapi %d", num1)
	} else {
		fmt.Println("✅ Tetap 5 karena integer adalah value type")
	}

	// Test pointer
	fmt.Println("\n--- Primitif dengan Pointer ---")
	num2 := 5
	fmt.Printf("SEBELUM: %d\n", num2)
	addTenByPointer(&num2)
	fmt.Printf("SETELAH: %d\n", num2)

	if num2 != 15 {
		t.Errorf("Num seharusnya 15, tapi %d", num2)
	} else {
		fmt.Println("✅ Berubah jadi 15 karena kita pakai pointer")
	}
}

// TestArrayVsSlice menunjukkan perbedaan array dan slice
func TestArrayVsSlice(t *testing.T) {
	fmt.Println("\n========== TEST: Array vs Slice ==========")

	// Fungsi modify array
	modifyArray := func(arr [3]int) {
		arr[0] = 999
		fmt.Printf("  [dalam fungsi] arr = %v\n", arr)
	}

	// Fungsi modify slice
	modifySlice := func(s []int) {
		s[0] = 999
		fmt.Printf("  [dalam fungsi] s = %v\n", s)
	}

	// Test Array - TIDAK berubah
	fmt.Println("\n--- Array (fixed size) ---")
	arr := [3]int{1, 2, 3}
	fmt.Printf("SEBELUM: %v\n", arr)
	modifyArray(arr)
	fmt.Printf("SETELAH: %v\n", arr)

	if arr[0] != 1 {
		t.Errorf("Array[0] seharusnya tetap 1, tapi %d", arr[0])
	} else {
		fmt.Println("✅ Array tidak berubah karena di-copy sepenuhnya")
	}

	// Test Slice - BERUBAH
	fmt.Println("\n--- Slice (dynamic size) ---")
	slice := []int{1, 2, 3}
	fmt.Printf("SEBELUM: %v\n", slice)
	modifySlice(slice)
	fmt.Printf("SETELAH: %v\n", slice)

	if slice[0] != 999 {
		t.Errorf("Slice[0] seharusnya 999, tapi %d", slice[0])
	} else {
		fmt.Println("✅ Slice berubah karena underlying data sama")
	}
}

// =============================================================================
// RINGKASAN HASIL YANG DIHARAPKAN
// =============================================================================
/*
Setelah menjalankan semua test dengan: go test -v -run "Test.*" ./test/specific/...

Anda akan melihat:

1. VALUE RECEIVER: Method dengan value receiver TIDAK mengubah data asli
   - Person.GrowOlderByValue() tidak mengubah Age asli
   - Counter.IncrementByValue() tidak mengubah Value asli

2. POINTER RECEIVER: Method dengan pointer receiver MENGUBAH data asli
   - Person.GrowOlderByPointer() mengubah Age asli
   - Counter.IncrementByPointer() mengubah Value asli

3. PARAMETER VALUE: Fungsi dengan parameter value TIDAK mengubah asli
   - updatePersonByValue(person, 100) tidak mengubah person.Age asli

4. PARAMETER POINTER: Fungsi dengan parameter pointer MENGUBAH asli
   - updatePersonByPointer(&person, 100) mengubah person.Age asli

5. SLICE (special case):
   - Modifikasi elemen: BERUBAH (karena reference ke underlying array sama)
   - Append tanpa pointer: TIDAK BERUBAH di luar fungsi
   - Append dengan pointer: BERUBAH

6. MAP: Selalu berubah karena map adalah reference type

7. ARRAY vs SLICE:
   - Array: TIDAK berubah (full copy)
   - Slice: BERUBAH (header copy, data sama)

KAPAN GUNAKAN POINTER?
- Ketika perlu mengubah data asli
- Ketika struct besar (hindari copy overhead)
- Ketika perlu nil value

KAPAN GUNAKAN VALUE?
- Ketika tidak perlu mengubah data asli
- Ketika struct kecil (seperti Point{X, Y})
- Ketika ingin immutability
*/
