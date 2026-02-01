package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {

	pipelinePattern()

	//Output:

}

func workerpool() {

	maxJobs := 10
	jobs := make(chan int, maxJobs)
	results := make(chan string, maxJobs)

	var wg sync.WaitGroup
	for wk := range 3 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for job := range jobs {
				result := fmt.Sprintf("Job ke-%d diterima Worker ke-%d\n", job, i)
				results <- result
			}
		}(wk)

	}

	go func() {
		for v := range maxJobs {
			fmt.Printf("Job ke-%d dikirim!\n", v)
			jobs <- v

		}
		close(jobs)
	}()
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Printf("%s", result)

	}
	//Output:

}

type Pakaian struct {
	Pemilik string
	Warna   string
	Tipe    string
	Status  string
	Number  int
}

func pipelinePattern() {

	//* pipeline adalah pattern yang menggunakan stages/tahap untuk tiap tiap pengerjaan task ,
	//*  misal buat laundry ada 100 baju , tahap 1.memasukan baju ke msin , 2. beri air,sabun dll ,3. keringkan ,4. lipat
	// TODO: Buat skenario laundry

	// masukan baju/generator
	clothes := generator(10)
	washedClothes := washing(clothes)
	dryedCloths := dryer(washedClothes)
	foldedClothes := folded(dryedCloths)

	for fd := range foldedClothes {

		fmt.Printf("\n\nPakaian %s Nomor-%d telah selesai\n\n", fd.Pemilik, fd.Number)

	}
	//Output:

}

func generator(max int) <-chan Pakaian {
	out := make(chan Pakaian, 10)
	color := []string{"red", "blue", "green", "yellow"}
	clothType := []string{"Baju", "Celana", "Kemeja", "Jeans"}
	go func() {
		for i := 0; i < max; i++ {
			out <- Pakaian{
				"adam",
				color[rand.Intn(len(color))],
				clothType[rand.Intn(len(clothType))],
				"Dirty",
				i,
			}

		}
		close(out)
	}()

	return out

	//Output:

}

func washing(dirtyClothes <-chan Pakaian) <-chan Pakaian {
	out := make(chan Pakaian, 10)
	go func() {
		for cloth := range dirtyClothes {
			time.Sleep(300 * time.Millisecond)
			out <- Pakaian{
				Pemilik: cloth.Pemilik,
				Warna:   cloth.Warna,
				Tipe:    cloth.Tipe,
				Number:  cloth.Number,
				Status:  "washed",
			}
			fmt.Printf("Pakaian %s Nomor-%d telah di cuci\n", cloth.Pemilik, cloth.Number)
		}
		close(out)
	}()
	//Output:
	return out

}

func dryer(washedClothes <-chan Pakaian) <-chan Pakaian {

	out := make(chan Pakaian, 10)

	go func() {
		for wc := range washedClothes {
			time.Sleep(225 * time.Millisecond)
			out <- Pakaian{
				Pemilik: wc.Pemilik,
				Warna:   wc.Warna,
				Tipe:    wc.Tipe,
				Number:  wc.Number,
				Status:  "dryed",
			}
			fmt.Printf("Pakaian %s Nomor-%d telah di keringkan\n", wc.Pemilik, wc.Number)

		}
		close(out)
	}()
	return out

}
func folded(dryedClothes <-chan Pakaian) <-chan Pakaian {

	out := make(chan Pakaian, 10)

	go func() {
		for dc := range dryedClothes {
			time.Sleep(125 * time.Millisecond)
			out <- Pakaian{
				Pemilik: dc.Pemilik,
				Warna:   dc.Warna,
				Tipe:    dc.Tipe,
				Status:  "folded",
				Number:  dc.Number,
			}
			fmt.Printf("Pakaian %s Nomor-%d telah di lipat\n", dc.Pemilik, dc.Number)

		}
		close(out)
	}()
	return out

}
