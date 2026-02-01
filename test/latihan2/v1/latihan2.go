package main

import "fmt"

func main() {
	section1_Array()
}

func section1_Array() {

	numbers := []int{1, 2, 3}
	animals := []string {"singa","rubah","gajah"}
	// numbersDeklarasiPanjang := [3]int{1, 2, 3, 4}
	fmt.Printf("numbers %v\n", numbers)
	fmt.Printf("Hewan %v\n", animals)
	// fmt.Printf("numbers %v", numbersDeklarasiPanjang)

	animals[1] = "ayam"

	fmt.Printf("Hewan setelah update %v\n", animals)

	for i, v := range animals {
		fmt.Printf("Hewan ke - %d %s\n", i + 1,v)
	}
}