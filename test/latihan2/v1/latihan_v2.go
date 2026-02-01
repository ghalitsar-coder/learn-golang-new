package main

import "fmt"

func main() {
	Method4()
}

func Slices1() {
	arr := []int{1, 2, 3, 4}
	arr2 := []string{"ayam", "ikan", "udang"}
	arr3 := []interface{}{1, false, "ikan"}

	valueAyam := arr2[0]
	fmt.Println(arr, arr2, arr3)
	fmt.Printf("makanan %s\n", valueAyam)


	// arr4 := make([]string , len(arr2))
	// copy(arr4,  arr2)
	// fmt.Printf("isi array 4 %v", arr4)
	// arr4[0] = ""
	// fmt.Printf("isi array 4 v2 %v", arr4)


	fmt.Printf("Isi dari arr2 %v", arr2)
	
	
	SliceReceiver(arr2)
	fmt.Printf("Isi dari arr2 AFTER UPDATE = %v", arr2)





}

func Map2()  {
	scores := map[string]int{
		"math":57,
		"comp":99,
		"datascience":84,
	}

	fmt.Printf("ini isi scores %v\n", scores)
	fmt.Printf("ini isi dari math %v\n", scores["math"])
	scores["math"] = 95
	fmt.Printf("ini isi updated scores %v\n", scores)

	if val , exist := scores["math"]; exist {
		fmt.Printf("nilai dari math adalah %d\n", val)
		
		}else {
			fmt.Println("tidak ada key math")

	}

	delete(scores ,"comp")

	fmt.Printf("updated scores %v\n", scores)

	jadwal := []string{"senin","selasa"}
	// score := 

	mhs := map[string]interface{}{
		"nama":"gal",
		"jadwal":jadwal,
		"scores" :map[string]map[string]int {
			"senin" :{"math":80,"ds":90},
			"selasa" :{"math":91,"ds":77},
		},
	}

	// iterasi

	for _, val := range mhs {
		switch v := val.(type){
			case string:
				fmt.Printf("ini value string %s\n", v)	
				case []string:	
				fmt.Printf("ini value array string %v\n", v)	
				case map[string]map[string]int:	
				for key, val := range v {
					
					fmt.Printf("key %s math %d ds %d \n",key,val["math"],val["ds"])	
				}
			default :
			fmt.Printf("unknown %T\n", v)
		}
		
	}


}

type Address struct {
	City string
	Village string
	Province string

}

type Mahasiswa struct {
	Name string
	Age int
	Hobbies []string
	Address *Address
} 

func NewAddress(city,village,province string ) *Address {
	return &Address{city,village,province}
}

func NewMahasiswa(name string ,age int, hobbies []string, address *Address) *Mahasiswa {
	return &Mahasiswa{name,age,hobbies,address}
}

func Struct3() {
	// hobbies := []string{}
	addr1 := NewAddress("bandung","cigadung","jawa barat")
	mhs := NewMahasiswa("adam",33,[]string{"renang","lari"}, addr1)

	fmt.Printf("isi dari mhs %+v\n", mhs)
	fmt.Printf("isi dari mhs address %+v\n", *mhs.Address)

}

func Method4()  {
	addr1 := NewAddress("bandung","cigadung","jawa barat")
	mhs := NewMahasiswa("megan",19,[]string{"lari","coding"},addr1)
	fmt.Printf("Mahasiswa - %+v", mhs)
	mhs.ChangeName("Dudung")
	fmt.Printf("Mahasiswa Updated - %+v", mhs)
}


func (m Mahasiswa) GreetingMhs() string {
	return fmt.Sprintf("hello %s ",m.Name)
}

func (m *Mahasiswa) ChangeName(name string)  {
	m.Name = name
}






func SliceReceiver(arr []string)  {
	arr[0] = "buaya"
}
