package main

import "fmt"

// =============================================================================
// DEMO: Map dengan dan tanpa Pointer
// =============================================================================

// modifyMapWithoutPointer - cara normal (tanpa pointer)
func modifyMapWithoutPointer(m map[string]int) {
	m["score"] = 100
	m["level"] = 5
	fmt.Println("  [dalam fungsi tanpa pointer] map =", m)
}

// modifyMapWithPointer - dengan pointer (redundant tapi VALID)
func modifyMapWithPointer(m *map[string]int) {
	(*m)["score"] = 200
	(*m)["level"] = 10
	fmt.Println("  [dalam fungsi dengan pointer] map =", *m)
}

func main() {
	fmt.Println("========== TEST 1: Map TANPA Pointer (Normal) ==========")
	map1 := map[string]int{"score": 0, "level": 1}
	fmt.Println("SEBELUM:", map1)
	modifyMapWithoutPointer(map1)
	fmt.Println("SETELAH:", map1)
	fmt.Println("✅ BERHASIL - Map berubah!")

	fmt.Println("\n========== TEST 2: Map DENGAN Pointer (Redundant tapi Valid) ==========")
	map2 := map[string]int{"score": 0, "level": 1}
	fmt.Println("SEBELUM:", map2)
	modifyMapWithPointer(&map2) // Perhatikan: &map2
	fmt.Println("SETELAH:", map2)
	fmt.Println("✅ BERHASIL - Map juga berubah!")

	fmt.Println("\n========== KESIMPULAN ==========")
	fmt.Println("Kedua cara SAMA-SAMA BERHASIL!")
	fmt.Println("Tapi pakai pointer untuk map itu REDUNDANT (tidak perlu)")
	fmt.Println("Karena map sudah secara internal adalah pointer")
}
