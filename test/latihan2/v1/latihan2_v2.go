package main

import "fmt"


func main() {
	MapNo2()	
}


func SlicesNo1()  {
	arr := []int{1,2,3,4,5}
	arr2 := make([]int , 5)
	arrCopy := make([]int,len(arr))

	copy(arrCopy,arr)

	// append menambah value di akhir
	
	arr2 = append(arr2, 999)
	arr2 = append(arr2, 999)
	arr2 = append(arr2, 999)
	arr2 = append(arr2, 999)
	arr2 = append(arr2, 999)
	arr2 = append(arr2,arr...)
	fmt.Printf("Ini isi arr1 = %v, arr2 = %v , arrCopy = %v", arr,arr2,arrCopy)
}


func MapNo2()  {
	data := map[string]int{}
	data["adam"] = 24
	fmt.Printf("Data : %v", data)

	nilaiAdam := map[string]int{
		"math":90,
		"eng":92,
	}

	if score,exists := nilaiAdam["math"]; exists{
		fmt.Printf("Keynya  ada , nilainya %v", score)
		if score < 90 {
			fmt.Printf("SCORE kurang dari 90, score %d", score)
			}else {
			fmt.Printf("SCORE lebih dari 90, score %d", score)

		}
	}else {

		fmt.Printf("tidak ada key tersebut ", )
	}


}



