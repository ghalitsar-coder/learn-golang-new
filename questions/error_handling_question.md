Q: Apakah ketika kita mengirim error dari suatu implementasi code, maka akan diterima di log atau di fmt.Errorf? Juga apakah mengirim error dengan fmt.Errorf cara yang lumrah?

Contoh kode:
```go
func main() {
    nums := []int{}
    result, err := sumVariadic(nums)
    if err != nil {
        log.Fatal("Error : ", err)
    } else {
        fmt.Printf("Result %d", result)
    }
}

func sumVariadic(nums []int) (int, error) {
    if len(nums) == 0 {
        return 0, fmt.Errorf("array kosong")
    }
    total := 0
    for _, v := range nums {
        total = total + v
    }
    return total, nil
}
```