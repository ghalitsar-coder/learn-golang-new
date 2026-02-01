package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// =============================================================================
// DEMO 1: Basic Context Cancellation
// =============================================================================

func demo1_BasicCancellation() {
	fmt.Println("\n" + "============================================================")
	fmt.Println("DEMO 1: Basic Context Cancellation")
	fmt.Println("============================================================\n")

	// Buat context dengan cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// Worker goroutine
	go func() {
		for i := 1; ; i++ {
			select {
			case <-ctx.Done():
				fmt.Printf("Worker: Berhenti karena %v\n", ctx.Err())
				return
			default:
				fmt.Printf("Worker: Iterasi ke-%d\n", i)
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Biarkan worker berjalan beberapa saat
	time.Sleep(1 * time.Second)

	// Cancel context!
	fmt.Println("Main: Memanggil cancel()...")
	cancel()

	// Tunggu sebentar untuk melihat output
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Main: Demo selesai")
}

// =============================================================================
// DEMO 2: Context with Timeout
// =============================================================================

func demo2_Timeout() {
	fmt.Println("\n" + "============================================================")
	fmt.Println("DEMO 2: Context with Timeout")
	fmt.Println("============================================================\n")

	// Context dengan timeout 2 detik
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Penting: selalu panggil cancel!

	fmt.Println("Mulai operasi dengan timeout 2 detik...")

	// Simulasi operasi yang memakan waktu
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Operasi selesai!") // Tidak akan tercapai
	case <-ctx.Done():
		fmt.Printf("Timeout! Error: %v\n", ctx.Err())
	}
}

// =============================================================================
// DEMO 3: Context with Deadline
// =============================================================================

func demo3_Deadline() {
	fmt.Println("\n" + "============================================================")
	fmt.Println("DEMO 3: Context with Deadline")
	fmt.Println("============================================================\n")

	// Set deadline 1.5 detik dari sekarang
	deadline := time.Now().Add(1500 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// Tampilkan informasi deadline
	d, ok := ctx.Deadline()
	if ok {
		fmt.Printf("Deadline: %v\n", d.Format(time.RFC3339Nano))
		fmt.Printf("Waktu tersisa: %v\n", time.Until(d))
	}

	fmt.Println("\nMemulai operasi...")

	// Simulasi beberapa tahap operasi
	for i := 1; i <= 3; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("Tahap %d: Dibatalkan - %v\n", i, ctx.Err())
			return
		case <-time.After(700 * time.Millisecond):
			fmt.Printf("Tahap %d: Selesai (tersisa: %v)\n", i, time.Until(d))
		}
	}

	fmt.Println("Semua tahap selesai!")
}

// =============================================================================
// DEMO 4: Context with Values
// =============================================================================

// Custom key types untuk menghindari collision
type contextKey string

const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
	roleKey      contextKey = "role"
)

func demo4_Values() {
	fmt.Println("\n" + "============================================================")
	fmt.Println("DEMO 4: Context with Values")
	fmt.Println("============================================================\n")

	// Buat context dengan values
	ctx := context.Background()
	ctx = context.WithValue(ctx, userIDKey, "user-123")
	ctx = context.WithValue(ctx, requestIDKey, "req-abc-456")
	ctx = context.WithValue(ctx, roleKey, "admin")

	// Pass ke function chain
	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	// Ambil values dari context
	userID, _ := ctx.Value(userIDKey).(string)
	requestID, _ := ctx.Value(requestIDKey).(string)
	role, _ := ctx.Value(roleKey).(string)

	fmt.Printf("Processing request:\n")
	fmt.Printf("  - Request ID: %s\n", requestID)
	fmt.Printf("  - User ID: %s\n", userID)
	fmt.Printf("  - Role: %s\n", role)

	// Panggil function berikutnya
	saveData(ctx)
}

func saveData(ctx context.Context) {
	requestID, _ := ctx.Value(requestIDKey).(string)
	fmt.Printf("\nSaving data (request: %s)...\n", requestID)
}

// =============================================================================
// DEMO 5: Parent-Child Context Cancellation
// =============================================================================

func demo5_ParentChildCancellation() {
	fmt.Println("\n" + "============================================================")
	fmt.Println("DEMO 5: Parent-Child Context Cancellation")
	fmt.Println("============================================================\n")

	// Parent context
	parentCtx, parentCancel := context.WithCancel(context.Background())

	// Child context 1
	child1Ctx, _ := context.WithCancel(parentCtx)

	// Child context 2 (dengan timeout sendiri)
	child2Ctx, child2Cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer child2Cancel()

	var wg sync.WaitGroup

	// Worker untuk child 1
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; ; i++ {
			select {
			case <-child1Ctx.Done():
				fmt.Printf("Child 1 Worker: Berhenti - %v\n", child1Ctx.Err())
				return
			default:
				fmt.Printf("Child 1 Worker: Running... (iterasi %d)\n", i)
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	// Worker untuk child 2
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; ; i++ {
			select {
			case <-child2Ctx.Done():
				fmt.Printf("Child 2 Worker: Berhenti - %v\n", child2Ctx.Err())
				return
			default:
				fmt.Printf("Child 2 Worker: Running... (iterasi %d)\n", i)
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	// Tunggu sebentar
	time.Sleep(1 * time.Second)

	// Cancel parent - semua child juga akan di-cancel!
	fmt.Println("\n>>> Cancelling PARENT context <<<\n")
	parentCancel()

	wg.Wait()
	fmt.Println("\nSemua worker sudah berhenti")
}

// =============================================================================
// DEMO 6: Worker Pool dengan Context
// =============================================================================

type Job struct {
	ID   int
	Data string
}

type Result struct {
	JobID  int
	Output string
	Error  error
}

func demo6_WorkerPool() {
	fmt.Println("\n" + "============================================================")
	fmt.Println("DEMO 6: Worker Pool dengan Context")
	fmt.Println("============================================================\n")

	// Context dengan timeout 3 detik
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	jobs := make(chan Job, 20)
	results := make(chan Result, 20)

	// Start workers
	var wg sync.WaitGroup
	numWorkers := 3

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go workerPoolWorker(ctx, i, jobs, results, &wg)
	}

	// Send jobs
	go func() {
		for i := 1; i <= 15; i++ {
			jobs <- Job{ID: i, Data: fmt.Sprintf("Task-%d", i)}
		}
		close(jobs)
	}()

	// Wait for workers dan close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	successCount := 0
	errorCount := 0

	for result := range results {
		if result.Error != nil {
			fmt.Printf("❌ Job %d: Error - %v\n", result.JobID, result.Error)
			errorCount++
		} else {
			fmt.Printf("✅ Job %d: %s\n", result.JobID, result.Output)
			successCount++
		}
	}

	fmt.Printf("\n📊 Summary: %d sukses, %d error\n", successCount, errorCount)
}

func workerPoolWorker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		select {
		case <-ctx.Done():
			results <- Result{JobID: job.ID, Error: ctx.Err()}
			continue
		default:
		}

		// Simulasi processing dengan waktu random
		processingTime := time.Duration(rand.Intn(500)+100) * time.Millisecond

		select {
		case <-time.After(processingTime):
			results <- Result{
				JobID:  job.ID,
				Output: fmt.Sprintf("Worker %d processed in %v", id, processingTime),
			}
		case <-ctx.Done():
			results <- Result{JobID: job.ID, Error: ctx.Err()}
		}
	}
}

// =============================================================================
// DEMO 7: Multiple Select Patterns
// =============================================================================

func demo7_SelectPatterns() {
	fmt.Println("\n============================================================")
	fmt.Println("DEMO 7: Multiple Select Patterns")
	fmt.Println("============================================================\n")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dataCh := make(chan string)
	errorCh := make(chan error)

	// Producer goroutine
	go func() {
		for i := 1; i <= 5; i++ {
			time.Sleep(300 * time.Millisecond)

			// Random error simulation
			if rand.Intn(5) == 0 {
				errorCh <- fmt.Errorf("random error at message %d", i)
				continue
			}
			dataCh <- fmt.Sprintf("Message %d", i)
		}
		close(dataCh)
	}()

	// Consumer dengan multiple select
	fmt.Println("Menunggu data dengan timeout 2 detik...\n")

	messageCount := 0
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\n⏰ Context done: %v\n", ctx.Err())
			break loop

		case err := <-errorCh:
			fmt.Printf("⚠️ Error received: %v\n", err)

		case msg, ok := <-dataCh:
			if !ok {
				fmt.Println("\n📭 Channel closed, semua data sudah diterima")
				break loop
			}
			messageCount++
			fmt.Printf("📨 Received: %s\n", msg)
		}
	}

	fmt.Printf("\n📊 Total messages received: %d\n", messageCount)
}

// =============================================================================
// MAIN
// =============================================================================

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     GOLANG CONTEXT DEMO - Belajar Context dengan Contoh    ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	// Jalankan semua demo
	demo1_BasicCancellation()
	demo2_Timeout()
	demo3_Deadline()
	demo4_Values()
	demo5_ParentChildCancellation()
	demo6_WorkerPool()
	demo7_SelectPatterns()

	fmt.Println("\n============================================================")
	fmt.Println("SEMUA DEMO SELESAI!")
	fmt.Println("============================================================\n")
}
