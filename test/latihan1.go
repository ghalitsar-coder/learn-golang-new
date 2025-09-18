package main

import (
	"fmt"
	"log"
)

func main() {

	nums:= []int{}
	result ,err:= sumVariadic(nums)
	 if err != nil {
		log.Fatal("Error : ",err)
	 }else {

		 fmt.Printf("Result %d", result)
		}


}


func sumVariadic(nums []int) (int, error) {
	if len(nums) == 0 {
		return 0,fmt.Errorf("array kosong")
	}
	total :=0 
	for _, v := range nums {
		total = total + v
		
	}
	return total, nil
	
}


func sum(num1 int , num2 int) int {
	return num1 + num2
}













