# Context di Golang: Best Practices, Anti-Patterns, dan Contoh Lengkap

## Best Practices

### 1. Selalu Panggil cancel()

```go
// ✅ BENAR: Gunakan defer cancel() segera setelah membuat context
func goodPractice() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // SELALU panggil cancel!
    
    // ... gunakan ctx
}

// ❌ SALAH: Lupa memanggil cancel
func badPractice() {
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    // cancel tidak pernah dipanggil = RESOURCE LEAK!
    
    // Timer internal masih berjalan sampai timeout
}
```

**Mengapa penting?**
- `WithCancel`, `WithDeadline`, `WithTimeout` mengalokasikan resources
- Jika `cancel()` tidak dipanggil, resources tidak akan di-release sampai:
  - Timeout tercapai
  - Parent context di-cancel
- Ini bisa menyebabkan memory leak pada aplikasi long-running

### 2. Context Sebagai Parameter Pertama

```go
// ✅ BENAR: Context adalah parameter pertama, dengan nama "ctx"
func DoSomething(ctx context.Context, userID string, data []byte) error {
    // ...
}

// ❌ SALAH: Context di tengah atau akhir
func DoSomethingBad(userID string, ctx context.Context, data []byte) error {
    // ...
}

// ❌ SALAH: Context sebagai struct field
type Service struct {
    ctx context.Context // JANGAN lakukan ini!
    db  *sql.DB
}
```

**Alasan:**
- Konsistensi dengan standard library Go
- Mudah ditemukan saat code review
- Memudahkan tooling dan linting

### 3. Jangan Simpan Context di Struct

```go
// ❌ SALAH: Context disimpan di struct
type BadClient struct {
    ctx context.Context // JANGAN!
    httpClient *http.Client
}

// ✅ BENAR: Context di-pass ke setiap method
type GoodClient struct {
    httpClient *http.Client
}

func (c *GoodClient) Get(ctx context.Context, url string) (*Response, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    return c.httpClient.Do(req)
}
```

**Mengapa:**
- Context bersifat request-scoped, bukan application-scoped
- Menyimpan context di struct bisa menyebabkan stale context
- Setiap operasi harus mendapat context yang fresh

### 4. Propagate Context ke Seluruh Call Chain

```go
// ✅ BENAR: Context di-propagate ke semua level
func HandleRequest(ctx context.Context, req *Request) error {
    // Level 1: Handler
    if err := validateRequest(ctx, req); err != nil {
        return err
    }
    
    // Level 2: Service
    result, err := processData(ctx, req.Data)
    if err != nil {
        return err
    }
    
    // Level 3: Repository
    return saveResult(ctx, result)
}

func validateRequest(ctx context.Context, req *Request) error {
    // Gunakan ctx untuk timeout awareness
    return nil
}

func processData(ctx context.Context, data []byte) (*Result, error) {
    // Propagate context ke external call
    return callExternalAPI(ctx, data)
}

func saveResult(ctx context.Context, result *Result) error {
    // Context untuk database query
    return db.ExecContext(ctx, "INSERT INTO results VALUES (?)", result)
}
```

### 5. Check Context Sebelum Operasi Mahal

```go
func processLargeDataset(ctx context.Context, items []Item) error {
    for i, item := range items {
        // ✅ Check context periodik untuk operasi panjang
        if i%100 == 0 {
            select {
            case <-ctx.Done():
                return ctx.Err()
            default:
            }
        }
        
        // Proses item
        if err := processItem(item); err != nil {
            return err
        }
    }
    return nil
}
```

### 6. Gunakan Custom Type untuk Context Keys

```go
// ✅ BENAR: Custom type prevents collision
type contextKey string

const (
    userIDKey   contextKey = "userID"
    tenantIDKey contextKey = "tenantID"
)

func SetUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (string, bool) {
    userID, ok := ctx.Value(userIDKey).(string)
    return userID, ok
}

// ❌ SALAH: String key langsung
func BadSetUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, "userID", userID) // Bisa collision!
}
```

### 7. Jangan Pass nil Context

```go
// ❌ SALAH: nil context
func badCall() {
    result, err := someFunction(nil, data) // PANIC potential!
}

// ✅ BENAR: Gunakan context.TODO() jika tidak yakin
func goodCall() {
    result, err := someFunction(context.TODO(), data)
}

// ✅ LEBIH BAIK: Gunakan context.Background() jika memang top-level
func bestCall() {
    ctx := context.Background()
    result, err := someFunction(ctx, data)
}
```

---

## Anti-Patterns (Pola yang Harus Dihindari)

### Anti-Pattern 1: Context di Struct untuk Long-lived Objects

```go
// ❌ ANTI-PATTERN
type UserService struct {
    ctx context.Context // SALAH!
    db  *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
    return &UserService{
        ctx: context.Background(), // Context ini akan "stale"
        db:  db,
    }
}

// ✅ SOLUSI: Pass context ke setiap method
type UserService struct {
    db *sql.DB
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    return s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)
}
```

### Anti-Pattern 2: Menggunakan Context untuk Optional Values

```go
// ❌ ANTI-PATTERN: Menggunakan context untuk configuration
func processData(ctx context.Context, data []byte) error {
    // SALAH: Debug mode seharusnya bukan context value
    debug := ctx.Value("debug").(bool)
    if debug {
        log.Println("Debug mode enabled")
    }
    // ...
}

// ✅ SOLUSI: Gunakan parameter eksplisit atau config struct
type ProcessConfig struct {
    Debug bool
}

func processData(ctx context.Context, data []byte, config ProcessConfig) error {
    if config.Debug {
        log.Println("Debug mode enabled")
    }
    // ...
}
```

### Anti-Pattern 3: Ignoring Context Cancellation

```go
// ❌ ANTI-PATTERN: Tidak check ctx.Done()
func slowOperation(ctx context.Context) error {
    for i := 0; i < 1000000; i++ {
        // Operasi ini tidak bisa di-cancel!
        heavyComputation()
    }
    return nil
}

// ✅ SOLUSI: Periodic check
func slowOperationFixed(ctx context.Context) error {
    for i := 0; i < 1000000; i++ {
        // Check setiap 1000 iterasi
        if i%1000 == 0 {
            select {
            case <-ctx.Done():
                return ctx.Err()
            default:
            }
        }
        heavyComputation()
    }
    return nil
}
```

### Anti-Pattern 4: Creating Too Many Context Layers

```go
// ❌ ANTI-PATTERN: Terlalu banyak layer WithValue
func badMiddleware(ctx context.Context) context.Context {
    ctx = context.WithValue(ctx, "a", 1)
    ctx = context.WithValue(ctx, "b", 2)
    ctx = context.WithValue(ctx, "c", 3)
    ctx = context.WithValue(ctx, "d", 4)
    ctx = context.WithValue(ctx, "e", 5)
    // 5 layers! Value lookup menjadi O(5)
    return ctx
}

// ✅ SOLUSI: Gunakan single struct untuk grouped values
type RequestInfo struct {
    A, B, C, D, E int
}

func goodMiddleware(ctx context.Context) context.Context {
    info := RequestInfo{1, 2, 3, 4, 5}
    return context.WithValue(ctx, requestInfoKey, info)
    // Hanya 1 layer!
}
```

### Anti-Pattern 5: Using context.Background() Everywhere

```go
// ❌ ANTI-PATTERN: Menggunakan Background() di dalam function
func processRequest(req *Request) error {
    ctx := context.Background() // SALAH! Ini membuat context baru
    return db.QueryContext(ctx, "SELECT...")
}

// ✅ SOLUSI: Terima context dari caller
func processRequest(ctx context.Context, req *Request) error {
    return db.QueryContext(ctx, "SELECT...")
}
```

---

## Contoh Implementasi Real-World Lengkap

### Contoh 1: REST API Server dengan Full Context Support

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/google/uuid"
)

// ============ Context Keys ============
type contextKey string

const (
    requestIDKey contextKey = "requestID"
    userIDKey    contextKey = "userID"
    startTimeKey contextKey = "startTime"
)

// ============ Middleware ============

// RequestIDMiddleware menambahkan unique request ID ke setiap request
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        ctx := context.WithValue(r.Context(), requestIDKey, requestID)
        ctx = context.WithValue(ctx, startTimeKey, time.Now())
        
        w.Header().Set("X-Request-ID", requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// TimeoutMiddleware menambahkan timeout ke setiap request
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), timeout)
            defer cancel()
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// LoggingMiddleware logs request details
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        requestID, _ := ctx.Value(requestIDKey).(string)
        
        log.Printf("[%s] Started %s %s", requestID, r.Method, r.URL.Path)
        
        // Wrap response writer untuk capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
        
        next.ServeHTTP(wrapped, r)
        
        startTime, _ := ctx.Value(startTimeKey).(time.Time)
        duration := time.Since(startTime)
        
        log.Printf("[%s] Completed %d in %v",
            requestID, wrapped.statusCode, duration)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// ============ Service Layer ============

type UserService struct {
    // Dependencies would go here (db, cache, etc.)
}

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    // Check if context is already cancelled
    if err := ctx.Err(); err != nil {
        return nil, err
    }
    
    requestID, _ := ctx.Value(requestIDKey).(string)
    log.Printf("[%s] UserService.GetUser called for id=%s", requestID, id)
    
    // Simulate database call
    select {
    case <-time.After(100 * time.Millisecond):
        return &User{
            ID:    id,
            Name:  "John Doe",
            Email: "john@example.com",
        }, nil
    case <-ctx.Done():
        log.Printf("[%s] GetUser cancelled: %v", requestID, ctx.Err())
        return nil, ctx.Err()
    }
}

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// ============ Handlers ============

func makeUserHandler(userService *UserService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        requestID, _ := ctx.Value(requestIDKey).(string)
        
        userID := r.URL.Query().Get("id")
        if userID == "" {
            writeJSON(w, http.StatusBadRequest, map[string]string{
                "error":      "id parameter is required",
                "request_id": requestID,
            })
            return
        }
        
        user, err := userService.GetUser(ctx, userID)
        if err != nil {
            if err == context.DeadlineExceeded {
                writeJSON(w, http.StatusGatewayTimeout, map[string]string{
                    "error":      "request timeout",
                    "request_id": requestID,
                })
                return
            }
            if err == context.Canceled {
                // Client disconnected
                log.Printf("[%s] Client disconnected", requestID)
                return
            }
            
            writeJSON(w, http.StatusInternalServerError, map[string]string{
                "error":      err.Error(),
                "request_id": requestID,
            })
            return
        }
        
        writeJSON(w, http.StatusOK, user)
    }
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, map[string]string{
        "status": "healthy",
        "time":   time.Now().Format(time.RFC3339),
    })
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

// ============ Main ============

func main() {
    // Initialize services
    userService := &UserService{}
    
    // Setup router
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/user", makeUserHandler(userService))
    
    // Chain middleware
    handler := RequestIDMiddleware(
        LoggingMiddleware(
            TimeoutMiddleware(5 * time.Second)(mux),
        ),
    )
    
    // Create server
    server := &http.Server{
        Addr:         ":8080",
        Handler:      handler,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }
    
    // Graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan
        
        log.Println("Shutting down server...")
        
        // Create shutdown context with timeout
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        if err := server.Shutdown(shutdownCtx); err != nil {
            log.Printf("Shutdown error: %v", err)
        }
    }()
    
    log.Printf("Server starting on %s", server.Addr)
    if err := server.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("Server error: %v", err)
    }
    
    log.Println("Server stopped gracefully")
}
```

### Contoh 2: Background Job Worker dengan Context

```go
package main

import (
    "context"
    "fmt"
    "log"
    "math/rand"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

// Job represents a unit of work
type Job struct {
    ID       string
    Type     string
    Payload  map[string]interface{}
    Priority int
}

// JobResult represents the outcome of processing a job
type JobResult struct {
    JobID   string
    Success bool
    Error   error
    Output  string
}

// Worker processes jobs from a queue
type Worker struct {
    id         int
    jobQueue   <-chan Job
    resultChan chan<- JobResult
    wg         *sync.WaitGroup
}

func NewWorker(id int, jobs <-chan Job, results chan<- JobResult, wg *sync.WaitGroup) *Worker {
    return &Worker{
        id:         id,
        jobQueue:   jobs,
        resultChan: results,
        wg:         wg,
    }
}

func (w *Worker) Start(ctx context.Context) {
    defer w.wg.Done()
    
    log.Printf("Worker %d: Starting", w.id)
    
    for {
        select {
        case <-ctx.Done():
            log.Printf("Worker %d: Received shutdown signal, stopping...", w.id)
            return
            
        case job, ok := <-w.jobQueue:
            if !ok {
                log.Printf("Worker %d: Job queue closed, stopping...", w.id)
                return
            }
            
            // Create per-job context with timeout
            jobCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
            result := w.processJob(jobCtx, job)
            cancel() // Always cancel after job completes
            
            // Send result if parent context still active
            select {
            case w.resultChan <- result:
            case <-ctx.Done():
                return
            }
        }
    }
}

func (w *Worker) processJob(ctx context.Context, job Job) JobResult {
    log.Printf("Worker %d: Processing job %s (type: %s)", w.id, job.ID, job.Type)
    
    // Simulate varying processing times
    processingTime := time.Duration(rand.Intn(2000)+500) * time.Millisecond
    
    select {
    case <-time.After(processingTime):
        // Job completed successfully
        output := fmt.Sprintf("Processed by worker %d in %v", w.id, processingTime)
        log.Printf("Worker %d: Job %s completed successfully", w.id, job.ID)
        
        return JobResult{
            JobID:   job.ID,
            Success: true,
            Output:  output,
        }
        
    case <-ctx.Done():
        // Context cancelled (timeout or shutdown)
        log.Printf("Worker %d: Job %s cancelled: %v", w.id, job.ID, ctx.Err())
        
        return JobResult{
            JobID:   job.ID,
            Success: false,
            Error:   ctx.Err(),
        }
    }
}

// JobProducer generates jobs
type JobProducer struct {
    jobQueue chan<- Job
}

func NewJobProducer(jobs chan<- Job) *JobProducer {
    return &JobProducer{jobQueue: jobs}
}

func (p *JobProducer) Start(ctx context.Context) {
    defer close(p.jobQueue)
    
    jobTypes := []string{"email", "sms", "push", "webhook"}
    jobCounter := 0
    
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            log.Println("Producer: Stopping job generation")
            return
            
        case <-ticker.C:
            jobCounter++
            job := Job{
                ID:       fmt.Sprintf("job-%d", jobCounter),
                Type:     jobTypes[rand.Intn(len(jobTypes))],
                Priority: rand.Intn(10),
                Payload: map[string]interface{}{
                    "data": fmt.Sprintf("payload-%d", jobCounter),
                },
            }
            
            select {
            case p.jobQueue <- job:
                log.Printf("Producer: Created job %s", job.ID)
            case <-ctx.Done():
                return
            }
        }
    }
}

// ResultCollector collects and processes results
type ResultCollector struct {
    resultChan <-chan JobResult
    mu         sync.Mutex
    stats      struct {
        total    int
        success  int
        failures int
    }
}

func NewResultCollector(results <-chan JobResult) *ResultCollector {
    return &ResultCollector{resultChan: results}
}

func (c *ResultCollector) Start(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            log.Println("Collector: Stopping result collection")
            return
            
        case result, ok := <-c.resultChan:
            if !ok {
                log.Println("Collector: Result channel closed")
                return
            }
            
            c.mu.Lock()
            c.stats.total++
            if result.Success {
                c.stats.success++
            } else {
                c.stats.failures++
            }
            c.mu.Unlock()
        }
    }
}

func (c *ResultCollector) GetStats() (total, success, failures int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.stats.total, c.stats.success, c.stats.failures
}

func main() {
    // Root context that will be cancelled on shutdown
    ctx, cancel := context.WithCancel(context.Background())
    
    // Setup signal handling for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Channels
    jobQueue := make(chan Job, 100)
    resultChan := make(chan JobResult, 100)
    
    // WaitGroup for workers
    var workerWg sync.WaitGroup
    
    // Start workers
    numWorkers := 3
    for i := 1; i <= numWorkers; i++ {
        workerWg.Add(1)
        worker := NewWorker(i, jobQueue, resultChan, &workerWg)
        go worker.Start(ctx)
    }
    
    // Start producer
    producer := NewJobProducer(jobQueue)
    go producer.Start(ctx)
    
    // Start result collector
    collector := NewResultCollector(resultChan)
    go collector.Start(ctx)
    
    // Print stats periodically
    statsTicker := time.NewTicker(5 * time.Second)
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case <-statsTicker.C:
                total, success, failures := collector.GetStats()
                log.Printf("Stats: Total=%d, Success=%d, Failures=%d",
                    total, success, failures)
            }
        }
    }()
    
    log.Println("Job processing system started. Press Ctrl+C to stop.")
    
    // Wait for shutdown signal
    <-sigChan
    log.Println("Shutdown signal received")
    
    // Cancel context to stop all components
    cancel()
    
    // Wait for workers to finish current jobs
    log.Println("Waiting for workers to finish...")
    workerWg.Wait()
    
    // Close result channel
    close(resultChan)
    
    // Final stats
    total, success, failures := collector.GetStats()
    log.Printf("Final Stats: Total=%d, Success=%d, Failures=%d",
        total, success, failures)
    
    log.Println("Shutdown complete")
}
```

---

## Ringkasan

Context di Go adalah tool yang sangat powerful untuk:

1. **Cancellation Propagation** - Membatalkan operasi secara cascade
2. **Timeout Management** - Mengontrol durasi maksimal operasi
3. **Request-Scoped Values** - Membawa data terkait request

### Kapan Menggunakan Context?

| Situasi | Gunakan Context? |
|---------|------------------|
| HTTP handlers | ✅ Ya (sudah ada dari `r.Context()`) |
| Database queries | ✅ Ya (`QueryContext`, `ExecContext`) |
| External API calls | ✅ Ya (`http.NewRequestWithContext`) |
| Background jobs | ✅ Ya (untuk graceful shutdown) |
| Unit tests | ✅ Ya (`context.Background()` atau dengan timeout) |
| Passing config | ❌ Tidak (gunakan parameter/struct) |
| Passing optional values | ❌ Tidak (gunakan functional options) |

### Key Takeaways

1. **Selalu panggil `cancel()`** dengan `defer`
2. **Context sebagai parameter pertama**
3. **Jangan simpan context di struct**
4. **Propagate context ke seluruh call chain**
5. **Check `ctx.Done()` untuk operasi panjang**
6. **Gunakan custom type untuk context keys**
7. **Jangan pass `nil` context**

Context adalah fondasi penting untuk membangun aplikasi Go yang robust, scalable, dan mudah di-maintain. Memahami cara kerjanya akan membantu Anda menulis kode concurrent yang lebih aman dan efisien.

---

*Materi ini adalah bagian dari seri belajar Golang. Happy Coding! 🚀*
