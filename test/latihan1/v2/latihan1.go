package main

import "fmt"

// Filter memilih elemen-elemen yang memenuhi predicate
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func FindIndex[T any](slice []T, predicate func(T) bool) int {
	for i, v := range slice {
		if predicate(v) {
			return i
		}
	}
	return -1
}

func main() {
	functions()
}

 

// functions
func functions() {
	sum := func(arr []interface{}) int {
		// Filter hanya angka (int) dari slice
		numbersOnly := Filter(arr, func(v interface{}) bool {
			_, ok := v.(int)
			return ok
		})

		var total int
		for _, v := range numbersOnly {
			if num, ok := v.(int); ok {
				total = total + num
			}
		}
		return total
	}

	numbers := []interface{}{1, 4, "halo", false, 5}

	result := sum(numbers)

	fmt.Printf("hasil dari result adalah %d\n", result)

	// variadic function
	mahasiswa := []string{"andi", "bob", "wisnu"}
	findIndexFromName := func(names []string, name string) int {
		return FindIndex(names, func(v string) bool {
			return v == name
		})

	}
	result2 := findIndexFromName(mahasiswa, "wisnu")
	fmt.Printf("content %v", result2)

	func() {
		fmt.Println("This is an anonymous function")
	}()


}

func Looping_4() {

	for i := range 5 {
		if i != 2 {
			fmt.Printf("Ini perulangan ke - %d\n", i+1)
		}
	}

	fmt.Println("")

	numbers := []int{1, 2, 3, 4, 5}
	for _, v := range numbers {
		fmt.Printf("Ini perulangan ke - %d\n", v)

	}

	for index, value := range numbers {
		fmt.Printf("Ini index-%d , ini value-%d\n", index, value)

	}

}

func VariabletTipeData_1() {
	var nama string = "andi"
	var umur int = 20

	nama2 := "budi"
	umur2 := 34

	// int bilangan bulat 32bit / 64 bit
	// int8, int16, int32, int64
	// uint unsigned int (bilangan bulat positif)
	// float32, float64 bilangan desimal
	// complex64, complex128 bilangan kompleks
	// rune = int32 merepresentasikan karakter unicode
	// byte = uint8 merepresentasikan 1 karakter ASCII
	// string
	// bool
	fmt.Println(nama, umur, nama2, umur2)

	var isReady bool = true
	count := 100
	var price float32 = 10.5
	var message string = "hello world"

	isLoggedIn := true
	userName := "staff"

	fmt.Println("helo")
	fmt.Println(isReady, count, price, message, isLoggedIn, userName)

	isUsernameValid := false
	isValid := isLoggedIn && isUsernameValid
	fmt.Printf("Isvalid is %v\n", isValid)
	isAdmin := userName != "admin"
	fmt.Printf("Is admin %v\n", isAdmin)

}

func Operator_2() {
	// 	### a. Operator Aritmatika
	// - `+` (Penjumlahan)
	// - `-` (Pengurangan)
	// - `*` (Perkalian)
	// - `/` (Pembagian)
	// - `%` (Modulo/Sisa bagi)

	// ### b. Operator Perbandingan
	// - `==` (Sama dengan)
	// - `!=` (Tidak sama dengan)
	// - `<` (Kurang dari)
	// - `<=` (Kurang dari atau sama dengan)
	// - `>` (Lebih dari)
	// - `>=` (Lebih dari atau sama dengan)

	// ### c. Operator Logika
	// - `&&` (AND)
	// - `||` (OR)
	// - `!` (NOT)

	// ### d. Operator Lainnya
	// - `=` (Assignment)
	// - `:=` (Short variable declaration)
	// - `+=`, `-=`, `*=`, `/=`, `%=` (Compound assignment)

	a := 10
	b := 15
	sum := a + b
	product := b * a
	quotient := a / b
	reminder := a % b

	fmt.Println("Hasil penjumlahan a + b =", sum)
	fmt.Println("Hasil perkalian a * b =", product)
	fmt.Println("Hasil pembagian a / b =", quotient)
	fmt.Println("Hasil sisa bagi a % b =", reminder)

	// Operator Perbandingan
	isEqual := a == b
	isNotEqual := a != b
	isLessThan := a < b
	isGreaterThanOrEqual := a >= b

	fmt.Println("Apakah a sama dengan b?", isEqual)
	fmt.Println("Apakah a tidak sama dengan b?", isNotEqual)
	fmt.Println("Apakah a kurang dari b?", isLessThan)
	fmt.Println("Apakah a lebih dari atau sama dengan b?", isGreaterThanOrEqual)

	isLoggedIn := true
	userName := "staff"

	isUsernameValid := false
	isValid := isLoggedIn && isUsernameValid
	fmt.Printf("Isvalid is %v\n", isValid)
	isStaff := userName != "" && userName != "admin"
	fmt.Printf("Is staff %v\n", isStaff)

}

func Branching_3() {
	age := 26
	hasMoney := true
	if age > 25 && hasMoney {
		fmt.Println("Anda dapat membeli motor")
	} else if age > 17 && hasMoney {
		fmt.Println("Anda dapat membeli sepeda")
	} else {
		fmt.Println("anda tidak bisa membeli motor atau sepeda")
	}

	switch age {
	case 20:
		fmt.Println("Umur anda 20 tahun")
	case 21:
		fmt.Println("Umur anda 21")
	default:
		fmt.Printf("umur anda dibawah 20 tahun atau diatas 21 tahun")
	}

}
