# Perbandingan Concurrency Patterns di Golang

File ini mendemonstrasikan 3 pola concurrency utama dalam satu skenario E-Commerce.

## 📌 Ringkasan Pola

| Pola | Konsep | Skenario Penggunaan |
| :--- | :--- | :--- |
| **1. Pipeline** | `Stage A` -> `Stage B` -> `Stage C` | Proses **Urut/Sequential**. Contoh: Validasi -> Packing -> Kirim. |
| **2. Worker Pool** | `Queue` -> `Limited Workers` | Membatasi Resource. Contoh: Max 3 Payment Processor. |
| **3. Fan-Out / Fan-In** | `Scatter` -> `Gather` | Agregasi Data Cepat. Contoh: Search di 3 DB sekaligus. |

---

## 💻 Source Code Implementasi

Berikut adalah kode lengkap yang bisa dijalankan (`go run test/specific/pattern_concurrency.md` jika di-rename ke .go):

```go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ========================================
// SKENARIO: SISTEM PEMROSESAN PESANAN E-COMMERCE
// ========================================

// Order represents pesanan dari customer
type Order struct {
	ID       int
	Customer string
	Items    []string
	Total    float64
}

// ProcessedOrder represents pesanan yang sudah diproses
type ProcessedOrder struct {
	Order
	ValidatedAt time.Time
	PackedAt    time.Time
	ShippedAt   time.Time
}

// ========================================
// 1. PIPELINE PATTERN
// ========================================
// Cocok untuk: Proses sequential yang harus berurutan
// Contoh: Validasi → Packing → Shipping (harus berurutan)

func runPipeline() {
	fmt.Println("\n=== PIPELINE PATTERN ===")
	fmt.Println("Skenario: Order processing yang harus sequential")
	fmt.Println("Flow: Validate → Pack → Ship (harus berurutan)\n")

	// Stage 1: Generate orders
	orders := generateOrders(5)

	// Stage 2: Validate orders
	validated := validateOrders(orders)

	// Stage 3: Pack orders
	packed := packOrders(validated)

	// Stage 4: Ship orders
	shipped := shipOrders(packed)

	// Collect results
	for order := range shipped {
		fmt.Printf("✓ Order #%d selesai diproses untuk %s\n", order.ID, order.Customer)
		fmt.Printf("  Timeline: Validated(%s) → Packed(%s) → Shipped(%s)\n",
			order.ValidatedAt.Format("15:04:05"),
			order.PackedAt.Format("15:04:05"),
			order.ShippedAt.Format("15:04:05"))
	}
}

func generateOrders(n int) <-chan Order {
	out := make(chan Order)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- Order{
				ID:       i,
				Customer: fmt.Sprintf("Customer-%d", i),
				Items:    []string{"Item-1", "Item-2"},
				Total:    float64(rand.Intn(1000) + 100),
			}
		}
	}()
	return out
}

func validateOrders(in <-chan Order) <-chan ProcessedOrder {
	out := make(chan ProcessedOrder)
	go func() {
		defer close(out)
		for order := range in {
			time.Sleep(200 * time.Millisecond) // Simulasi validasi
			processed := ProcessedOrder{
				Order:       order,
				ValidatedAt: time.Now(),
			}
			fmt.Printf("  [Validate] Order #%d validated\n", order.ID)
			out <- processed
		}
	}()
	return out
}

func packOrders(in <-chan ProcessedOrder) <-chan ProcessedOrder {
	out := make(chan ProcessedOrder)
	go func() {
		defer close(out)
		for order := range in {
			time.Sleep(300 * time.Millisecond) // Simulasi packing
			order.PackedAt = time.Now()
			fmt.Printf("  [Pack] Order #%d packed\n", order.ID)
			out <- order
		}
	}()
	return out
}

func shipOrders(in <-chan ProcessedOrder) <-chan ProcessedOrder {
	out := make(chan ProcessedOrder)
	go func() {
		defer close(out)
		for order := range in {
			time.Sleep(150 * time.Millisecond) // Simulasi shipping
			order.ShippedAt = time.Now()
			fmt.Printf("  [Ship] Order #%d shipped\n", order.ID)
			out <- order
		}
	}()
	return out
}

// ========================================
// 2. WORKER POOL PATTERN
// ========================================
// Cocok untuk: Banyak task independent yang perlu dibatasi concurrency-nya
// Contoh: Process payment dari banyak order secara parallel (max 3 workers)

type PaymentTask struct {
	OrderID int
	Amount  float64
}

type PaymentResult struct {
	OrderID   int
	Success   bool
	ProcessAt time.Time
}

func runWorkerPool() {
	fmt.Println("\n\n=== WORKER POOL PATTERN ===")
	fmt.Println("Skenario: Process payment dengan batasan max 3 concurrent workers")
	fmt.Println("Kenapa butuh pool? Mencegah overload payment gateway\n")

	numWorkers := 3
	tasks := make(chan PaymentTask, 10)
	results := make(chan PaymentResult, 10)

	// Start worker pool
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go paymentWorker(i, tasks, results, &wg)
	}

	// Send tasks
	go func() {
		for i := 1; i <= 10; i++ {
			tasks <- PaymentTask{
				OrderID: i,
				Amount:  float64(rand.Intn(500) + 100),
			}
		}
		close(tasks)
	}()

	// Close results when all workers done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	successCount := 0
	for result := range results {
		if result.Success {
			successCount++
			fmt.Printf("✓ Payment Order #%d berhasil diproses pada %s\n",
				result.OrderID, result.ProcessAt.Format("15:04:05"))
		}
	}
	fmt.Printf("\nTotal successful payments: %d/10\n", successCount)
}

func paymentWorker(id int, tasks <-chan PaymentTask, results chan<- PaymentResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("  [Worker-%d] Processing payment for Order #%d\n", id, task.OrderID)
		time.Sleep(500 * time.Millisecond) // Simulasi payment gateway call
		
		results <- PaymentResult{
			OrderID:   task.OrderID,
			Success:   true,
			ProcessAt: time.Now(),
		}
	}
	fmt.Printf("  [Worker-%d] selesai bekerja\n", id)
}

// ========================================
// 3. FAN-OUT / FAN-IN PATTERN
// ========================================
// Cocok untuk: Distribute work ke multiple goroutines, lalu merge hasilnya
// Contoh: Search produk di multiple databases, merge hasil jadi satu

type SearchQuery struct {
	Keyword string
}

type SearchResult struct {
	Source   string
	Products []string
	Duration time.Duration
}

func runFanOutFanIn() {
	fmt.Println("\n\n=== FAN-OUT / FAN-IN PATTERN ===")
	fmt.Println("Skenario: Search produk di 3 database berbeda secara parallel")
	fmt.Println("Fan-Out: Kirim query ke 3 DB sekaligus")
	fmt.Println("Fan-In: Merge semua hasil jadi satu stream\n")

	query := SearchQuery{Keyword: "laptop"}

	// FAN-OUT: Distribute query ke multiple sources
	mysqlResults := searchMySQL(query)
	postgresResults := searchPostgres(query)
	mongoResults := searchMongo(query)

	// FAN-IN: Merge semua results
	allResults := mergeResults(mysqlResults, postgresResults, mongoResults)

	// Process merged results
	totalProducts := 0
	for result := range allResults {
		fmt.Printf("✓ Hasil dari %s (waktu: %v):\n", result.Source, result.Duration)
		for _, product := range result.Products {
			fmt.Printf("  - %s\n", product)
			totalProducts++
		}
	}
	fmt.Printf("\nTotal produk ditemukan: %d\n", totalProducts)
}

func searchMySQL(query SearchQuery) <-chan SearchResult {
	out := make(chan SearchResult)
	go func() {
		defer close(out)
		start := time.Now()
		time.Sleep(300 * time.Millisecond) // Simulasi query MySQL
		
		out <- SearchResult{
			Source:   "MySQL",
			Products: []string{"Dell Laptop XPS", "HP Pavilion"},
			Duration: time.Since(start),
		}
	}()
	return out
}

func searchPostgres(query SearchQuery) <-chan SearchResult {
	out := make(chan SearchResult)
	go func() {
		defer close(out)
		start := time.Now()
		time.Sleep(450 * time.Millisecond) // Simulasi query Postgres
		
		out <- SearchResult{
			Source:   "PostgreSQL",
			Products: []string{"Lenovo ThinkPad", "Asus ROG"},
			Duration: time.Since(start),
		}
	}()
	return out
}

func searchMongo(query SearchQuery) <-chan SearchResult {
	out := make(chan SearchResult)
	go func() {
		defer close(out)
		start := time.Now()
		time.Sleep(250 * time.Millisecond) // Simulasi query MongoDB
		
		out <- SearchResult{
			Source:   "MongoDB",
			Products: []string{"Macbook Pro", "Acer Aspire"},
			Duration: time.Since(start),
		}
	}()
	return out
}

// Fan-In: Merge multiple channels into one
func mergeResults(channels ...<-chan SearchResult) <-chan SearchResult {
	out := make(chan SearchResult)
	var wg sync.WaitGroup

	// Function to copy from channel to output
	output := func(c <-chan SearchResult) {
		defer wg.Done()
		for result := range c {
			out <- result
		}
	}

	// Fan-in: Start goroutine for each input channel
	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	// Close output when all inputs are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ========================================
// COMPARISON & MAIN
// ========================================

func main() {
	// rand.Seed(time.Now().UnixNano()) // Deprecated in new Go, but fine for demo

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║  PERBANDINGAN CONCURRENCY PATTERNS DI GOLANG               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	// 1. Pipeline Pattern
	start := time.Now()
	runPipeline()
	fmt.Printf("\nWaktu total: %v\n", time.Since(start))

	// 2. Worker Pool Pattern
	start = time.Now()
	runWorkerPool()
	fmt.Printf("\nWaktu total: %v\n", time.Since(start))

	// 3. Fan-Out/Fan-In Pattern
	start = time.Now()
	runFanOutFanIn()
	fmt.Printf("\nWaktu total: %v\n", time.Since(start))
}
```