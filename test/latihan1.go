package main

import (
	"fmt"
	"log"
)

func main() {
	nums := []int{}
	result, err := sumVariadic(nums)
	if err != nil {
		// Menggunakan log.Println daripada log.Fatal untuk hanya mencatat error tanpa menghentikan program
		log.Println("Error:", err)
		// Alternatif: fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %d\n", result)
	}

	// Contoh kedua dengan array yang memiliki elemen
	nums2 := []int{1, 2, 3, 4, 5}
	result2, err2 := sumVariadic(nums2)
	if err2 != nil {
		log.Println("Error:", err2)
	} else {
		fmt.Printf("Result: %d\n", result2)
	}
}

func sumVariadic(nums []int) (int, error) {
	if len(nums) == 0 {
		// Menggunakan fmt.Errorf untuk membuat error dengan pesan yang jelas
		return 0, fmt.Errorf("array kosong")
	}
	total := 0
	for _, v := range nums {
		total = total + v
	}
	return total, nil
}