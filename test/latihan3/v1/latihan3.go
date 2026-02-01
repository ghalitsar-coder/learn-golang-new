package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type IntORFloat interface {
	int | float64
}

func divide[T IntORFloat](a, b T) (T, error) {
	if b < 0 {
		return 0, fmt.Errorf("b = %d, Nilai B tidak boleh kurang dari 0", b)
	}
	return a / b, nil

}

func main() {
	// fmt.Println("========== PATTERN 1: Result Struct (RECOMMENDED) ==========")
	// goRoutine4()

	// fmt.Println("\n========== PATTERN 2: Two Separate Channels ==========")
	// goRoutine5()

	workerPoolPattern()
}

func workerPoolPattern() {
	// 2 channel 1 job ,1 result
	jobs := make(chan int, 10)
	results := make(chan string, 10)

	// kita bikin wg buat wait goroutine
	var wg sync.WaitGroup

	// buat worker 3
	for i := 0; i < 3; i++ {
		fmt.Printf("worker ke-%d dibuat", i)
		// buat go routine
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			// ambil tiap job
			for job := range jobs {
				// simulasi pake time.sleep
				time.Sleep(100 * time.Millisecond)
				result := fmt.Sprintf("Job ke-%d , Worker ke-%d\n", job, idx)
				results <- result
				// isi content dari channel results pake sprintf
			}
		}(i)
	}

	// kirim job untuk 10 job
	for i := 0; i < 10; i++ {
		jobs <- i
		fmt.Printf("mengirim job ke-%d\n", i)
	}

	close(jobs)

	wg.Wait()
	close(results)
	// wait
	// bisa cek results
	for result := range results {
		fmt.Printf("%s\n", result)
	}
	//Output:

}

// ===== PATTERN 1: Result Struct (RECOMMENDED) =====
// Struct untuk menampung hasil dan error sekaligus
type Result struct {
	Value int
	Error error
}

func goRoutine4() {
	// Channel yang mengirim Result (berisi value DAN error)
	resultChan := make(chan Result, 1)

	var wg sync.WaitGroup

	// Test case 1: Error (b <= 0)
	wg.Add(1)
	go func(a, b int, ch chan Result) {
		defer wg.Done()

		if b <= 0 {
			// Kirim Result dengan error, value = 0
			ch <- Result{
				Value: 0,
				Error: fmt.Errorf("Error: b tidak bisa kurang dari atau sama dengan 0 (b = %d)", b),
			}
			return
		}

		// Kirim Result dengan value, error = nil
		ch <- Result{
			Value: a / b,
			Error: nil,
		}
	}(10, 0, resultChan) // b = 0, akan error

	fmt.Println("Menunggu goroutine...")
	wg.Wait()

	// Terima result
	result := <-resultChan
	if result.Error != nil {
		fmt.Printf("❌ Terjadi error: %v\n", result.Error)
	} else {
		fmt.Printf("✅ Hasil: %d\n", result.Value)
	}

	// Test case 2: Success (b > 0)
	wg.Add(1)
	go func(a, b int, ch chan Result) {
		defer wg.Done()

		if b <= 0 {
			ch <- Result{
				Value: 0,
				Error: fmt.Errorf("Error: b tidak bisa kurang dari atau sama dengan 0 (b = %d)", b),
			}
			return
		}

		ch <- Result{
			Value: a / b,
			Error: nil,
		}
	}(20, 4, resultChan) // b = 4, akan sukses

	wg.Wait()
	result = <-resultChan
	if result.Error != nil {
		fmt.Printf("❌ Terjadi error: %v\n", result.Error)
	} else {
		fmt.Printf("✅ Hasil: %d\n", result.Value)
	}
}

// ===== PATTERN 2: Two Separate Channels =====
func goRoutine5() {
	valueChan := make(chan int, 1)
	errorChan := make(chan error, 1)

	var wg sync.WaitGroup

	// Test dengan error
	wg.Add(1)
	go func(a, b int, valCh chan int, errCh chan error) {
		defer wg.Done()

		if b <= 0 {
			// Kirim error ke error channel
			errCh <- fmt.Errorf("Error: b tidak bisa <= 0 (b = %d)", b)
			return
		}

		// Kirim value ke value channel
		valCh <- a / b
	}(10, 0, valueChan, errorChan)

	wg.Wait()

	// Cek channel mana yang dapat data
	select {
	case err := <-errorChan:
		fmt.Printf("❌ Terjadi error: %v\n", err)
	case val := <-valueChan:
		fmt.Printf("✅ Hasil: %d\n", val)
	default:
		fmt.Println("Tidak ada data")
	}

	// Test dengan success
	wg.Add(1)
	go func(a, b int, valCh chan int, errCh chan error) {
		defer wg.Done()

		if b <= 0 {
			errCh <- fmt.Errorf("Error: b tidak bisa <= 0 (b = %d)", b)
			return
		}

		valCh <- a / b
	}(20, 5, valueChan, errorChan)

	wg.Wait()

	select {
	case err := <-errorChan:
		fmt.Printf("❌ Terjadi error: %v\n", err)
	case val := <-valueChan:
		fmt.Printf("✅ Hasil: %d\n", val)
	default:
		fmt.Println("Tidak ada data")
	}
}

func goRoutine6() {

	var wg sync.WaitGroup
	data := make(chan Result, 1)

	wg.Add(1)
	go func(a, b int, data chan Result) {
		defer wg.Done()
		if b <= 0 {
			data <- Result{
				Value: 0,
				Error: fmt.Errorf("Nilai B = %d , B tidak boleh <= 0", b),
			}
			return
		}

		data <- Result{
			Value: a / b,
			Error: nil,
		}

	}(10, 0, data)

	wg.Wait()

	result := <-data

	if result.Error != nil {
		fmt.Printf("Error : `%+v`", result.Error)

	} else {
		fmt.Printf("Hasilnya adalah %d", result.Value)
	}

	//Output:

}

func goRoutine3() {

	var wg sync.WaitGroup

	wg.Add(1)
	go printNumbers(&wg)
	fmt.Println("main finished")
	wg.Wait()
	fmt.Println("after 5 ?")

	//Output:

}
func printNumbers(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 5; i++ {
		fmt.Println(i)
		time.Sleep(100 * time.Millisecond)
	}
}

func goRoutineTest() {
	done := make(chan bool)

	go printNumbers2(done)

	fmt.Println("Main finished")

	// Menunggu sinyal dari goroutine
	<-done
}

func printNumbers2(done chan bool) {

	for i := range 5 {
		fmt.Printf("%d\n", i+1)
		time.Sleep(100 * time.Millisecond)
	}

	done <- true
	//Output:

}

func goRoutine2() {

	done := make(chan bool)

	go printNumbers2(done)

	fmt.Println("Main finished")

	<-done

}

func goroutine_3() {
	go sayHello("arjit")
	go sayHello("jajang")

	time.Sleep(1 * time.Second)
	fmt.Println("main function finished")
}
func sayHello(name string) {

	for index := range 3 {
		fmt.Printf("Hello %s , index-%d\n", name, index)
		time.Sleep(100 * time.Millisecond)

	}

	//Output:

}

func readFile(filename string) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		fmt.Printf("Error : %+v", err)
	}

	//Output:

}

func mightPanic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("program dilanjutkan")

		}
	}()

	fmt.Println("About to panic")
	panic("This is panic!")
	fmt.Println("this line will not be executed")

}

func defer_recover_panic_error_2() {

	mightPanic()
	defer fmt.Println("akan di deferred print last/3")
	fmt.Println("print ke-1")
	defer fmt.Println("deferred print ke-2")
	defer fmt.Println("deferred print ke-1")
	fmt.Println("print ke-2")

	readFile("zero.text")

}

func error_handling_1() {

	result, err := divide(10, -5)
	if err != nil {
		fmt.Printf("Error : %v", err)
	} else {
		fmt.Printf("Hasilnya :  %v", result)
	}
	//Output:

}
