# Bagian 5: Lanjutan dan Studi Kasus

## 1. Reflection

Reflection adalah kemampuan program untuk memeriksa struktur tipe dan nilai pada saat runtime. Package `reflect` menyediakan fungsi untuk ini.

### a. Dasar Reflection
```go
import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	p := Person{Name: "Alice", Age: 30}
	
	// Mendapatkan reflect.Type
	t := reflect.TypeOf(p)
	fmt.Println("Type:", t.Name())        // Output: Person
	fmt.Println("Kind:", t.Kind())        // Output: struct
	
	// Mendapatkan reflect.Value
	v := reflect.ValueOf(p)
	fmt.Println("Value:", v)              // Output: {Alice 30}
	
	// Mengakses field
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fmt.Printf("Field: %s, Type: %s, Value: %v, Tag: %s\n", 
			field.Name, field.Type, value, field.Tag)
	}
}
```

### b. Penggunaan Praktis Reflection
Reflection sering digunakan dalam framework serialisasi (seperti `encoding/json`), ORM, dan alat bantu debugging.

Contoh sederhana untuk mencetak struct:
```go
func PrintStruct(s interface{}) {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	
	if t.Kind() != reflect.Struct {
		fmt.Println("Bukan struct")
		return
	}
	
	fmt.Printf("Struct: %s\n", t.Name())
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fmt.Printf("  %s (%s) = %v\n", field.Name, field.Type, value)
	}
}

func main() {
	p := Person{Name: "Bob", Age: 25}
	PrintStruct(p)
}
```

Catatan: Reflection bisa memperlambat program dan membuat kode kurang aman tipe, jadi gunakan dengan bijak.

## 2. Context

Package `context` digunakan untuk mengirim sinyal pembatalan, batas waktu, dan data lintas API dan antar proses.

### a. Membatalkan Goroutine
```go
import (
	"context"
	"fmt"
	"time"
)

func longRunningTask(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled:", ctx.Err())
			return
		default:
			fmt.Println("Working...")
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	// Membuat context dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // Pastikan cancel dipanggil
	
	go longRunningTask(ctx)
	
	// Menunggu context selesai
	<-ctx.Done()
	fmt.Println("Main: Context finished")
}
```

### b. Membawa Data dengan Context
```go
func main() {
	// Membuat context dengan value
	ctx := context.WithValue(context.Background(), "userID", 123)
	
	// Mengambil value
	userID := ctx.Value("userID")
	if id, ok := userID.(int); ok {
		fmt.Println("User ID:", id)
	}
}
```

## 3. HTTP Server/Client

### a. HTTP Server Dasar
```go
import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	
	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed to start:", err)
	}
}
```

### b. HTTP Server dengan `http.ServeMux`
```go
func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from mux!")
	})
	
	mux.HandleFunc("/goodbye", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Goodbye!")
	})
	
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	
	fmt.Println("Mux server starting on :8080...")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server failed to start:", err)
	}
}
```

### c. HTTP Client
```go
import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}
	
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}
```

## 4. Database Integration

### a. Menggunakan SQLite dengan `database/sql`
Pertama, install driver SQLite:
```bash
go get github.com/mattn/go-sqlite3
```

Contoh kode:
```go
import (
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	// Membuka koneksi database
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	// Membuat tabel
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL
	);`
	
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	
	// Menyisipkan data
	insertQuery := `INSERT INTO users(name, age) VALUES (?, ?)`
	_, err = db.Exec(insertQuery, "Alice", 25)
	if err != nil {
		log.Fatal(err)
	}
	
	_, err = db.Exec(insertQuery, "Bob", 30)
	if err != nil {
		log.Fatal(err)
	}
	
	// Mengambil data
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	
	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Age)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	
	fmt.Println("Users:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}
```

## 5. Studi Kasus: Membangun Aplikasi RESTful API Sederhana

Kita akan membuat aplikasi REST API untuk manajemen buku sederhana dengan operasi CRUD.

Struktur direktori:
```
bookstore/
├── go.mod
├── main.go
├── handlers/
│   └── book_handlers.go
├── models/
│   └── book.go
└── storage/
    └── book_storage.go
```

### Langkah 1: Inisialisasi Module
```bash
go mod init bookstore
```

### Langkah 2: Model (`models/book.go`)
```go
package models

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Price  float64 `json:"price"`
}
```

### Langkah 3: Storage (`storage/book_storage.go`)
```go
package storage

import (
	"sync"
	"bookstore/models"
)

type BookStorage struct {
	books map[int]models.Book
	nextID int
	mu    sync.RWMutex
}

func NewBookStorage() *BookStorage {
	return &BookStorage{
		books: make(map[int]models.Book),
		nextID: 1,
	}
}

func (s *BookStorage) Create(book models.Book) models.Book {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	book.ID = s.nextID
	s.nextID++
	s.books[book.ID] = book
	return book
}

func (s *BookStorage) GetAll() []models.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	books := make([]models.Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}
	return books
}

func (s *BookStorage) GetByID(id int) (models.Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	book, exists := s.books[id]
	return book, exists
}

func (s *BookStorage) Update(id int, updatedBook models.Book) (models.Book, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.books[id]; !exists {
		return models.Book{}, false
	}
	
	updatedBook.ID = id
	s.books[id] = updatedBook
	return updatedBook, true
}

func (s *BookStorage) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.books[id]; !exists {
		return false
	}
	
	delete(s.books, id)
	return true
}
```

### Langkah 4: Handlers (`handlers/book_handlers.go`)
```go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	
	"bookstore/models"
	"bookstore/storage"
)

type BookHandler struct {
	storage *storage.BookStorage
}

func NewBookHandler(storage *storage.BookStorage) *BookHandler {
	return &BookHandler{storage: storage}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	createdBook := h.storage.Create(book)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdBook)
}

func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books := h.storage.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID dari URL path, misal: /books/1
	idStr := r.PathValue("id") // Go 1.22+
	// Untuk versi sebelum 1.22, gunakan: idStr := mux.Vars(r)["id"] dengan gorilla/mux
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}
	
	book, exists := h.storage.GetByID(id)
	if !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}
	
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	updatedBook, exists := h.storage.Update(id, book)
	if !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}
	
	if !h.storage.Delete(id) {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}
```

### Langkah 5: Main (`main.go`)
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	
	"bookstore/handlers"
	"bookstore/storage"
)

func main() {
	bookStorage := storage.NewBookStorage()
	bookHandler := handlers.NewBookHandler(bookStorage)
	
	mux := http.NewServeMux()
	
	// Routes
	mux.HandleFunc("POST /books", bookHandler.CreateBook)
	mux.HandleFunc("GET /books", bookHandler.GetAllBooks)
	mux.HandleFunc("GET /books/{id}", bookHandler.GetBookByID)
	mux.HandleFunc("PUT /books/{id}", bookHandler.UpdateBook)
	mux.HandleFunc("DELETE /books/{id}", bookHandler.DeleteBook)
	
	// Simple middleware untuk logging
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("[%s] %s\n", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
	
	server := &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(mux),
	}
	
	fmt.Println("Bookstore API server starting on :8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
```

### Langkah 6: Menjalankan Aplikasi
1.  Simpan semua file sesuai struktur di atas.
2.  Jalankan aplikasi:
    ```bash
    go run main.go
    ```
3.  Gunakan `curl` atau alat seperti Postman untuk menguji API:

    - **Create Book (POST /books)**:
      ```bash
      curl -X POST http://localhost:8080/books \
        -H "Content-Type: application/json" \
        -d '{"title": "Go Programming", "author": "John Doe", "price": 29.99}'
      ```

    - **Get All Books (GET /books)**:
      ```bash
      curl http://localhost:8080/books
      ```

    - **Get Book by ID (GET /books/1)**:
      ```bash
      curl http://localhost:8080/books/1
      ```

    - **Update Book (PUT /books/1)**:
      ```bash
      curl -X PUT http://localhost:8080/books/1 \
        -H "Content-Type: application/json" \
        -d '{"title": "Advanced Go Programming", "author": "John Doe", "price": 39.99}'
      ```

    - **Delete Book (DELETE /books/1)**:
      ```bash
      curl -X DELETE http://localhost:8080/books/1
      ```

### Kesimpulan Proyek
Aplikasi ini menunjukkan:
- Penggunaan package untuk organisasi kode
- Pembuatan dan penggunaan module
- Implementasi RESTful API dengan HTTP server bawaan Go
- Struktur aplikasi yang bersih dan terpisah (models, storage, handlers)
- Operasi CRUD dasar
- Penanganan error HTTP
- Penggunaan middleware untuk logging

Dengan membangun proyek ini, Anda telah menerapkan sebagian besar konsep penting dalam pemrograman Go, dari dasar hingga topik yang lebih lanjutan.