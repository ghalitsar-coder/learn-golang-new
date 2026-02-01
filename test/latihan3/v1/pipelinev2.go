package main

import (
	"fmt"
	"time"
)

func main() {

	//Output:

	nums := generator2([]int{1, 2, 3, 4, 5})
	filteredResult := filtered(nums)

	for v := range filteredResult {
		fmt.Printf("Data : %d\n", v)
	}

}

//  generator

func generator2(nums []int) <-chan int {
	out := make(chan int)

	go func() {
		for _, num := range nums {
			out <- num
			time.Sleep(200 * time.Millisecond)

		}
		close(out)
	}()
	//Output:
	return out

}

func filtered(values <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for val := range values {

			out <- val * 2
			time.Sleep(300 * time.Millisecond)

		}
		close(out)

	}()

	return out

	//Output:

}
