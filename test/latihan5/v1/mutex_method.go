package main

import (
	"fmt"
	"sync"
	"time"
)

type DataStore struct {
	data  map[string]string
	mutex sync.RWMutex
}

func NewDataStore() *DataStore {
	return &DataStore{
		data: make(map[string]string),
	}
}

func (ds *DataStore) Write(key, value string) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	fmt.Printf("Writing %s = %s\n", key, value)
	ds.data[key] = value
	time.Sleep(100 * time.Millisecond) // Simulasi penulisan
}

func (ds *DataStore) Read(key string) string {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	fmt.Printf("Reading %s\n", key)
	time.Sleep(50 * time.Millisecond) // Simulasi pembacaan
	return ds.data[key]
}

func main() {
	store := NewDataStore()

	// Goroutine writer
	go func() {

		durationName := time.Now()
		store.Write("name", "Alice")
		elapsedName := time.Since(durationName)
		fmt.Printf("durasi WRITE %v\n", elapsedName.Round(time.Since(durationName)))

		durationCity := time.Now()
		store.Write("city", "New York")
		elapsedCity := time.Since(durationCity)
		fmt.Printf("durasi WRITE %v\n", elapsedCity.Round(time.Since(durationCity)))
	}()

	// Goroutine readers
	go func() {
		duration := time.Now()
		for i := 0; i < 3; i++ {
			name := store.Read("name")
			fmt.Printf("Reader 1 got name: %s\n", name)
		}
		elapsed := time.Since(duration)
		fmt.Printf("durasi READ 1 %v\n", elapsed.Round(time.Since(duration)))
	}()

	go func() {
		duration := time.Now()

		for i := 0; i < 3; i++ {
			city := store.Read("city")
			fmt.Printf("Reader 2 got city: %s\n", city)
		}
		elapsed := time.Since(duration)
		fmt.Printf("durasi READ-2 %v\n", elapsed.Round(time.Since(duration)))
	}()

	time.Sleep(2 * time.Second)
}
