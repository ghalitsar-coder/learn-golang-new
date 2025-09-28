package main

import "fmt"


func main() {

	CreateFunctionNo5()
}

func TipeDataNo1()  {
	fmt.Println("Hello world")	

	name := "John Doe"
	age := 33
	isLoggedIn := false
	price := 10.24

	fmt.Println("My name is %s and my age is %d , logged in : %t , price : %.2f\n", name,age,isLoggedIn,price)
}

func  OperatorNo2()  {
	a := 10
	b:= 20
	sum := a + b
	diff := a - b
	product := a * b
	quatient := a / b
	reminder := a % b

	fmt.Println("Aritmatika", sum,diff,product,quatient,reminder)


	isEqual := a == b
	isNotEqual := a != b
	isGreater := a > b
	fmt.Println("Perbandingan",isEqual,isNotEqual,isGreater)


	x := 1
	notResult := (x)
	fmt.Println("Isi Not Result " ,notResult)

}

func PercabanganNo3()  {
	x := 5
	if x > 5 && x < 10 {
		fmt.Println("Nilai x lebih dari 5 dan kurang dari 10,  x = ",x)
	}else if x >= 10 {
		fmt.Println("Nilai x merupakan 10 atau lebih, x = ", x)
	}else {
		fmt.Println("Nilai X = 5 atau x kurang dari 5 , X = ", x)
	}

	hari := "senin"
	switch hari {
	 
	case "jumat":
		fmt.Println("ketemu hari ini weekend ,hari ini hari ", hari)
	case "sabtu":
		fmt.Println("ketemu hari ini weekend ,hari ini hari ", hari)
	default :
		fmt.Println("ketemu hari ini weekday ,hari ini hari ", hari)
	}
}


func PerulanganNo4()  {
	for i := 0; i < 5; i++ {
		fmt.Println("Pengulangan 5 kali , sekarang index ke - ",i + 1)
	}

	arr := []string{"ayam","ikan","udang"}

	for _, v := range arr {
		fmt.Printf("bisa %s\n",v)
	}

	arr2 := []map[string]interface{}{
		{"nama":"adam","umur":30},
		{"nama":"rusdi","umur":24},
		{"nama":"ahmad","umur":24},
	}
	for _, person := range arr2 {
		if person["nama"] == "rusdi" {
			person["nama"] = "mamang"
		}
		fmt.Printf("Nama %s , Umur %v\n", person["nama"],person["umur"])
	}

}

func CreateFunctionNo5()  {
	name := "John"
	greet := func () string {
		return fmt.Sprintf("Hello %s",name)
	}

	fmt.Println("Greet ", greet())
	fmt.Println("Greet person ", greetPerson("gibran"))


	fmt.Printf("Hasil dari variadic function adalah %d", variadicFunction(1,2,3,4,5))
	
}

func greetPerson(name string) string {
	return fmt.Sprintf("Hello %s",name)
	
}

func variadicFunction(numbers ...int) int {
	total :=0
	for _, v := range numbers {
		total = total + v
	}
	return total
}