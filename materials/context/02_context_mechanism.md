# Context di Golang: Cara Kerja dan Mekanisme Internal

## Cara Kerja Internal Context

### Struktur Data Internal

Context di Go diimplementasikan menggunakan beberapa struct internal. Mari kita lihat bagaimana masing-masing bekerja:

```go
// emptyCtx - implementasi untuk Background() dan TODO()
type emptyCtx int

var (
    background = new(emptyCtx)
    todo       = new(emptyCtx)
)

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) { return }
func (*emptyCtx) Done() <-chan struct{}                   { return nil }
func (*emptyCtx) Err() error                              { return nil }
func (*emptyCtx) Value(key any) any                       { return nil }
```

### Mekanisme Cancellation

#### Bagaimana Cancel Bekerja?

```
┌────────────────────────────────────────────────────────────────────┐
│                     CANCELLATION FLOW                               │
├────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   1. cancel() dipanggil                                            │
│          │                                                          │
│          ▼                                                          │
│   2. c.done channel di-CLOSE (bukan send!)                         │
│          │                                                          │
│          ▼                                                          │
│   3. c.err di-set ke Canceled/DeadlineExceeded                     │
│          │                                                          │
│          ▼                                                          │
│   4. Semua child context juga di-cancel (propagation)              │
│          │                                                          │
│          ▼                                                          │
│   5. Semua goroutine yang listen ke ctx.Done() akan terbangun      │
│                                                                     │
└────────────────────────────────────────────────────────────────────┘
```

#### Kode Internal Cancellation (Simplified)

```go
// Struktur internal cancelCtx (disederhanakan)
type cancelCtx struct {
    Context                         // Parent context
    mu       sync.Mutex             // Protects fields below
    done     chan struct{}          // Created lazily, closed by cancel
    children map[*cancelCtx]struct{} // Set of children
    err      error                  // Set when cancelled
}

func (c *cancelCtx) cancel(removeFromParent bool, err error) {
    c.mu.Lock()
    if c.err != nil {
        c.mu.Unlock()
        return // Already cancelled
    }
    c.err = err
    
    // Close the done channel
    if c.done == nil {
        c.done = closedchan
    } else {
        close(c.done)
    }
    
    // Cancel all children
    for child := range c.children {
        child.cancel(false, err)
    }
    c.children = nil
    c.mu.Unlock()
    
    // Remove from parent if needed
    if removeFromParent {
        removeChild(c.Context, c)
    }
}
```

### Mekanisme Deadline/Timeout

```go
// timerCtx menambahkan timer ke cancelCtx
type timerCtx struct {
    cancelCtx
    timer    *time.Timer // Underlying timer
    deadline time.Time   // Waktu deadline
}

// Ketika timer fires:
func (c *timerCtx) cancel(removeFromParent bool, err error) {
    c.cancelCtx.cancel(false, err)
    if removeFromParent {
        removeChild(c.cancelCtx.Context, c)
    }
    c.mu.Lock()
    if c.timer != nil {
        c.timer.Stop() // Stop timer untuk mencegah leak
        c.timer = nil
    }
    c.mu.Unlock()
}
```

**Timeline Eksekusi:**

```
Timeline untuk context.WithTimeout(ctx, 3*time.Second):

    t=0s                                      t=3s
    │                                         │
    ▼                                         ▼
    ┌─────────────────────────────────────────┐
    │   Context Aktif                         │
    │                                         │
    │   timer := time.AfterFunc(3s, func(){  │
    │       cancel()                          │──────► ctx.Done() closed
    │   })                                    │       ctx.Err() = DeadlineExceeded
    └─────────────────────────────────────────┘
```

### Mekanisme Value Lookup

Context Value bekerja dengan **chain lookup** - mencari dari context saat ini ke parent-nya.

```go
// valueCtx menyimpan key-value pair
type valueCtx struct {
    Context    // Parent context
    key, val any
}

func (c *valueCtx) Value(key any) any {
    if c.key == key {
        return c.val // Ditemukan!
    }
    return c.Context.Value(key) // Cari ke parent
}
```

**Visualisasi Value Lookup:**

```
Mencari key="userID":

    ┌─────────────────┐
    │   Background    │  Value("userID") → nil
    └────────┬────────┘
             │
    ┌────────▼────────┐
    │ WithValue       │  
    │ key="requestID" │  Value("userID") → cari ke parent
    │ val="req-123"   │
    └────────┬────────┘
             │
    ┌────────▼────────┐
    │ WithValue       │
    │ key="userID"    │  Value("userID") → "user-456" ✓ DITEMUKAN!
    │ val="user-456"  │
    └────────┬────────┘
             │
    ┌────────▼────────┐
    │ WithValue       │
    │ key="role"      │  ctx.Value("userID") dimulai di sini
    │ val="admin"     │  ↑ Pencarian ke atas
    └─────────────────┘

Kompleksitas: O(n) di mana n = kedalaman context chain
```

### Parent-Child Relationship

```go
// Demonstrasi hubungan parent-child
func demonstrateParentChild() {
    // Root
    ctx1 := context.Background()
    
    // Child 1 of root
    ctx2, cancel2 := context.WithCancel(ctx1)
    defer cancel2()
    
    // Child 2 of ctx2
    ctx3, cancel3 := context.WithTimeout(ctx2, 5*time.Second)
    defer cancel3()
    
    // Child 3 of ctx3
    ctx4 := context.WithValue(ctx3, "key", "value")
    
    // Relationship visualization:
    //
    //    ctx1 (Background)
    //      └── ctx2 (WithCancel)
    //            └── ctx3 (WithTimeout)
    //                  └── ctx4 (WithValue)
    //
    // Jika cancel2() dipanggil:
    // - ctx2.Done() closed ✓
    // - ctx3.Done() closed ✓ (child dari ctx2)
    // - ctx4.Done() closed ✓ (child dari ctx3)
    
    _ = ctx4
}
```

### Race Condition Prevention

Context dirancang untuk **thread-safe**:

```go
// Context aman digunakan dari multiple goroutines secara bersamaan
func demonstrateThreadSafety(ctx context.Context) {
    var wg sync.WaitGroup
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            // Semua operasi ini thread-safe!
            select {
            case <-ctx.Done():
                return
            default:
                // ctx.Deadline() - thread-safe
                // ctx.Err() - thread-safe
                // ctx.Value(key) - thread-safe
            }
        }(i)
    }
    
    wg.Wait()
}
```

---

## Skenario Penggunaan Context

### Skenario 1: HTTP Server dengan Timeout

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

// Simulasi database query yang lambat
func fetchFromDatabase(ctx context.Context, userID string) (map[string]any, error) {
    // Simulasi query database yang memakan waktu 3 detik
    select {
    case <-time.After(3 * time.Second):
        return map[string]any{
            "id":    userID,
            "name":  "John Doe",
            "email": "john@example.com",
        }, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    // HTTP request sudah memiliki context bawaan!
    ctx := r.Context()
    
    // Tambahkan timeout tambahan jika diperlukan
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    userID := r.URL.Query().Get("id")
    if userID == "" {
        userID = "default-user"
    }
    
    log.Printf("Fetching user: %s", userID)
    
    user, err := fetchFromDatabase(ctx, userID)
    if err != nil {
        if err == context.DeadlineExceeded {
            http.Error(w, "Request timeout", http.StatusGatewayTimeout)
            log.Printf("Request timeout for user: %s", userID)
            return
        }
        if err == context.Canceled {
            // Client disconnect
            log.Printf("Client disconnected for user: %s", userID)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func main() {
    http.HandleFunc("/user", userHandler)
    
    server := &http.Server{
        Addr:         ":8080",
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    
    log.Println("Server starting on :8080")
    log.Fatal(server.ListenAndServe())
}
```

### Skenario 2: Database Query dengan Context

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"
    
    _ "github.com/lib/pq"
)

type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
    // Context akan meng-cancel query jika timeout atau dibatalkan
    query := `SELECT id, name, email FROM users WHERE id = $1`
    
    row := r.db.QueryRowContext(ctx, query, id)
    
    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        if err == context.DeadlineExceeded {
            return nil, fmt.Errorf("query timeout: %w", err)
        }
        return nil, err
    }
    
    return &user, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]User, error) {
    query := `SELECT id, name, email FROM users`
    
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        // Check context sebelum setiap iterasi
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }
        
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, nil
}

type User struct {
    ID    int64
    Name  string
    Email string
}

func main() {
    db, err := sql.Open("postgres", "postgres://user:pass@localhost/dbname?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    repo := &UserRepository{db: db}
    
    // Contoh penggunaan dengan timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    user, err := repo.FindByID(ctx, 1)
    if err != nil {
        log.Printf("Error finding user: %v", err)
        return
    }
    
    fmt.Printf("Found user: %+v\n", user)
}
```

### Skenario 3: Parallel API Calls dengan Graceful Shutdown

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
)

type APIResult struct {
    Source string
    Data   map[string]interface{}
    Error  error
}

func fetchAPI(ctx context.Context, url string, source string) APIResult {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return APIResult{Source: source, Error: err}
    }
    
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return APIResult{Source: source, Error: err}
    }
    defer resp.Body.Close()
    
    var data map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return APIResult{Source: source, Error: err}
    }
    
    return APIResult{Source: source, Data: data}
}

func aggregateAPIs(ctx context.Context) ([]APIResult, error) {
    apis := []struct {
        url    string
        source string
    }{
        {"https://api.github.com", "github"},
        {"https://jsonplaceholder.typicode.com/users/1", "jsonplaceholder"},
        {"https://httpbin.org/get", "httpbin"},
    }
    
    // Create a child context with timeout
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    results := make([]APIResult, len(apis))
    var wg sync.WaitGroup
    
    for i, api := range apis {
        wg.Add(1)
        go func(idx int, url, source string) {
            defer wg.Done()
            results[idx] = fetchAPI(ctx, url, source)
        }(i, api.url, api.source)
    }
    
    // Wait for all goroutines to complete
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        return results, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func main() {
    ctx := context.Background()
    
    results, err := aggregateAPIs(ctx)
    if err != nil {
        fmt.Printf("Aggregation failed: %v\n", err)
        return
    }
    
    for _, result := range results {
        if result.Error != nil {
            fmt.Printf("[%s] Error: %v\n", result.Source, result.Error)
        } else {
            fmt.Printf("[%s] Success: got %d fields\n", result.Source, len(result.Data))
        }
    }
}
```

### Skenario 4: Worker Pool dengan Cancellation

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID   int
    Output  string
    Err     error
}

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: Shutting down due to cancellation\n", id)
            return
        case job, ok := <-jobs:
            if !ok {
                fmt.Printf("Worker %d: No more jobs, exiting\n", id)
                return
            }
            
            // Proses job dengan awareness terhadap context
            result := processJob(ctx, id, job)
            
            // Kirim hasil jika context masih aktif
            select {
            case results <- result:
            case <-ctx.Done():
                return
            }
        }
    }
}

func processJob(ctx context.Context, workerID int, job Job) Result {
    // Simulasi proses yang memakan waktu
    processingTime := time.Duration(rand.Intn(500)+100) * time.Millisecond
    
    select {
    case <-time.After(processingTime):
        return Result{
            JobID:  job.ID,
            Output: fmt.Sprintf("Worker %d processed: %s", workerID, job.Data),
        }
    case <-ctx.Done():
        return Result{
            JobID: job.ID,
            Err:   ctx.Err(),
        }
    }
}

func main() {
    // Context dengan timeout 2 detik
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    jobs := make(chan Job, 100)
    results := make(chan Result, 100)
    
    // Start worker pool
    var wg sync.WaitGroup
    numWorkers := 3
    
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(ctx, i, jobs, results, &wg)
    }
    
    // Send jobs
    go func() {
        for i := 1; i <= 20; i++ {
            jobs <- Job{ID: i, Data: fmt.Sprintf("Task-%d", i)}
        }
        close(jobs)
    }()
    
    // Wait for workers and close results
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    successCount := 0
    errorCount := 0
    
    for result := range results {
        if result.Err != nil {
            fmt.Printf("Job %d failed: %v\n", result.JobID, result.Err)
            errorCount++
        } else {
            fmt.Printf("Job %d: %s\n", result.JobID, result.Output)
            successCount++
        }
    }
    
    fmt.Printf("\nSummary: %d success, %d errors\n", successCount, errorCount)
}
```

### Skenario 5: Request Tracing dengan Context Values

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"
    
    "github.com/google/uuid"
)

// Custom context key types
type contextKey string

const (
    requestIDKey contextKey = "requestID"
    userIDKey    contextKey = "userID"
    startTimeKey contextKey = "startTime"
)

// Middleware untuk menambahkan request ID
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        
        // Tambahkan request ID ke context
        ctx := context.WithValue(r.Context(), requestIDKey, requestID)
        ctx = context.WithValue(ctx, startTimeKey, time.Now())
        
        // Tambahkan ke response header untuk debugging
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Middleware untuk authentication
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Simulasi authentication
        userID := r.Header.Get("X-User-ID")
        if userID == "" {
            userID = "anonymous"
        }
        
        ctx := context.WithValue(r.Context(), userIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Helper untuk logging dengan context
func logWithContext(ctx context.Context, message string) {
    requestID, _ := ctx.Value(requestIDKey).(string)
    userID, _ := ctx.Value(userIDKey).(string)
    startTime, _ := ctx.Value(startTimeKey).(time.Time)
    
    elapsed := time.Since(startTime)
    
    log.Printf("[%s] [user:%s] [elapsed:%v] %s",
        requestID, userID, elapsed, message)
}

// Handler
func apiHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    logWithContext(ctx, "Handler started")
    
    // Simulasi processing
    result, err := processWithContext(ctx)
    if err != nil {
        logWithContext(ctx, fmt.Sprintf("Error: %v", err))
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    logWithContext(ctx, "Handler completed successfully")
    fmt.Fprintf(w, "Result: %s", result)
}

func processWithContext(ctx context.Context) (string, error) {
    logWithContext(ctx, "Processing started")
    
    // Check context sebelum operasi yang mahal
    select {
    case <-ctx.Done():
        return "", ctx.Err()
    default:
    }
    
    // Simulasi processing
    time.Sleep(100 * time.Millisecond)
    
    logWithContext(ctx, "Processing completed")
    return "Success!", nil
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api", apiHandler)
    
    // Chain middlewares
    handler := requestIDMiddleware(authMiddleware(mux))
    
    log.Println("Server starting on :8080")
    http.ListenAndServe(":8080", handler)
}
```

---

*Lanjutan di file 03_context_best_practices.md*
