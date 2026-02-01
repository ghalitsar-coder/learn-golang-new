package main

import (
	"fmt"
)

func main2() {
	// Contoh penggunaan pipeline
	// Kita kirim angka 1, 2, 3, 4, 5
	// Generator akan mengembalikan channel yang berisi angka-angka tersebut
	stream := generator(1, 2, 3, 4, 5)

	// Kita baca hasilnya dari channel
	for num := range stream {
		fmt.Printf("Menerima: %d\n", num)
	}
}

func generator(nums ...int) <-chan int {
	// 1. Buat channel output
	out := make(chan int)

	// 2. Jalankan goroutine (Closure)
	// Di sini kita TIDAK PERLU passing 'out' atau 'nums' lewat parameter.
	// Anonymous function ini bisa melihat variable 'out' dan 'nums' milik 'generator'.
	go func() {
		for _, n := range nums {
			// Kita bisa akses 'out' langsung karena dia ada di scope pembungkus (closure)
			out <- n
		}
		// Jangan lupa tutup channel kalau sudah selesai
		close(out)
	}()

	// 3. Kembalikan channel (Ingat: Goroutine di atas bekerja di background!)
	return out
}
