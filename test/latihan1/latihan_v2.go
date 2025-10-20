package main

import "fmt"

func main() {
	FunctionFn5()
}

func VariableTypeData1() {
	nama := "adam"
	umur := 23
	isLoggedIn := true
	rupiah := 124.444
	fmt.Printf("nama saya %s , umur saya %d, apakah sudah login %v , rupiah %d", nama,umur,isLoggedIn,rupiah )
}


func Operator2()  {
	// buatkana aritmetika , logika ,perbandingan

	a := 10
	b := 20
	c := a + b
	d := c / 2
	e := d * 2
	f := e % c
	fmt.Printf("isi f %v", f)

	if a % 2 == 0 {
		fmt.Printf("A adalah bilangan genap , a = %d\n", a)
	}else {
		fmt.Printf("A adalah bilangan ganjil , a = %d\n", a)
	}

	and := a >  5 && b < 100
	or := c > 10 || c < 50
	not := a != 10000000000

	fmt.Printf("ini adalah and %v ,or %v dan not %v\n", and,or,not )

	a2 := c == e
	a3 := a  != f
	a4 := c> a
	a5 := d <= a
	a6 := f >= a
	fmt.Println(a2,a3,a4,a5,a6)

}


func Percabangan3()  {
	a := 10

	if a == 10 {
		fmt.Printf("a == 10  %d \n", a)
	}

	if a > 20 {
		fmt.Printf("a > 20 %d \n", a)

	}else {
		fmt.Printf("a < 20 %d \n", a)
	}


	hari := "selasa"
	switch hari {
		case "senin",
		 "selasa",
		 "rabu",
		 "kamis",
		 "jumat":
			fmt.Printf("weekday %v \n", hari)
		default:
			fmt.Printf("weekend %v \n", hari)
			
		}


	n:=10
	switch  {
	case n > 10:
		fmt.Printf("n > 10 , n = %d\n", n)
	case n == 10:
		fmt.Printf("n == 10 , n = %d\n", n)
	default:
		fmt.Printf("n < 10 , n = %d\n", n)
	}

}

func Perulangan4()  {
	for i := 0; i < 5; i++ {
		fmt.Printf("i ke %d\n", i + 1)
	}

	arr := []int{1,2,3,4,5}
	for _, v := range arr {
		fmt.Printf("value %v\n", v)
	}

}

func FunctionFn5()  {

// getString := helloFunction("hello",55,true)

	fmt.Printf(" %v", helloFunction("hello",55,true))
	fmt.Printf("summarizeNumber %d\n", summarizeNumber(1,2,2,2,22,2,2))
	greeting := func () string  {
		return fmt.Sprintf("hello world")
	}
	fmt.Printf("greeting %s", greeting())


	
}






func helloFunction(a string , b int ,c bool) string {
	return fmt.Sprintf("isi dari string %s, isi int %d int , isi boolean %v",a,b,c)
}


func summarizeNumber(numbers ...int) int {
	total := 0
	for _, v := range numbers {
		total = total + v
	}
	return total
	
}