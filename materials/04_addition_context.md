# Context di Golang: Panduan Lengkap dan Mendalam

## Pendahuluan

Context adalah salah satu fitur paling penting dalam ekosistem Go yang diperkenalkan dalam Go 1.7. Package `context` menyediakan cara standar untuk membawa deadline, sinyal pembatalan, dan nilai-nilai lain yang bersifat request-scoped melintasi batas-batas API dan antar goroutine.

## Konsep Dasar Context

### Apa itu Context?

Context adalah sebuah interface yang membawa informasi tentang deadline, sinyal pembatalan, dan nilai-nilai lain yang diperlukan oleh operasi-operasi yang berjalan dalam suatu request atau transaksi. Context sangat penting dalam aplikasi concurrent karena memungkinkan kita untuk:

1. **Membatalkan operasi yang sedang berjalan** - Ketika user membatalkan request atau timeout tercapai
2. **Menyebarkan deadline** - Memberitahu semua operasi kapan mereka harus berhenti
3. **Membawa nilai request-scoped** - Seperti request ID, user authentication, dll

### Interface Context

```go
type Context interface {
    // Deadline mengembalikan waktu ketika work yang dilakukan atas nama context ini
    // harus dibatalkan. Deadline mengembalikan ok==false ketika tidak ada deadline yang diset.
    Deadline() (deadline time.Time, ok bool)
    
    // Done mengembalikan channel yang ditutup ketika work yang dilakukan atas nama context
    // ini harus dibatalkan.
    Done() <-chan struct{}
    
    // Err mengembalikan error yang menjelaskan mengapa Done channel ditutup.
    Err() error
    
    // Value mengembalikan nilai yang diasosiasikan dengan key ini.
    Value(key interface{}) interface{}
}
```

## Mekanisme Kerja Context

### 1. Context Tree (Pohon Context)

Context di Go bekerja dalam struktur pohon hierarkis. Ketika Anda membuat context child dari parent context, child tersebut mewarisi properti dari parent-nya. Namun, pembatalan hanya mengalir dari parent ke child, tidak sebaliknya.

```
                    Background/TODO (root)
                           |
                    WithCancel (parent)
                      /          \
            WithTimeout         WithValue
                 |                  |
           WithValue          WithDeadline
```

**Aturan penting:**
- Jika parent context dibatalkan, semua child context juga akan dibatalkan
- Membatalkan child context tidak mempengaruhi parent atau sibling contexts
- Setiap context hanya bisa dibatalkan sekali

### 2. Channel Komunikasi

Context menggunakan channel `Done()` untuk komunikasi pembatalan. Ketika context dibatalkan:
1. Channel `Done()` ditutup
2. Semua goroutine yang mendengarkan channel ini akan menerima signal
3. Method `Err()` akan mengembalikan alasan pembatalan

### 3. Propagasi Deadline

Ketika Anda membuat context dengan deadline:
- Context akan otomatis dibatalkan ketika deadline tercapai
- Child context dapat memiliki deadline lebih pendek dari parent, tapi tidak bisa lebih panjang
- Deadline yang paling pendek akan selalu digunakan

## Jenis-Jenis Context

### 1. context.Background()

Context kosong yang tidak pernah dibatalkan, tidak memiliki deadline, dan tidak membawa nilai. Biasanya digunakan sebagai root context.

```go
ctx := context.Background()
```

**Kapan menggunakan:**
- Di fungsi main()
- Di inisialisasi
- Di tests
- Sebagai top-level context untuk incoming requests

### 2. context.TODO()

Mirip dengan Background(), tapi digunakan ketika Anda tidak yakin context apa yang harus digunakan atau ketika fungsi belum menerima context parameter.

```go
ctx := context.TODO()
```

**Kapan menggunakan:**
- Ketika refactoring code yang belum menggunakan context
- Sebagai placeholder sementara

### 3. context.WithCancel()

Membuat context yang dapat dibatalkan secara manual.

```go
ctx, cancel := context.WithCancel(parent)
defer cancel() // Penting: selalu panggil cancel untuk mencegah resource leak
```

**Cara kerja:**
1. Membuat copy dari parent context
2. Mengembalikan context baru dan fungsi cancel
3. Ketika cancel() dipanggil, channel Done() akan ditutup
4. Semua goroutine yang listening ke Done() akan menerima signal

### 4. context.WithDeadline()

Membuat context yang otomatis dibatalkan pada waktu tertentu.

```go
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(parent, deadline)
defer cancel()
```

**Cara kerja:**
1. Membuat context dengan deadline absolut
2. Otomatis membatalkan context ketika deadline tercapai
3. `Err()` akan mengembalikan `context.DeadlineExceeded`

### 5. context.WithTimeout()

Shorthand untuk WithDeadline yang menggunakan duration relatif.

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()
```

**Cara kerja:**
1. Internally memanggil `WithDeadline(parent, time.Now().Add(timeout))`
2. Lebih mudah digunakan karena menggunakan durasi relatif

### 6. context.WithValue()

Membawa nilai key-value dalam context.

```go
ctx := context.WithValue(parent, "userID", 12345)
```

**Cara kerja:**
1. Membuat context baru yang menyimpan pasangan key-value
2. Nilai dapat diakses dengan `ctx.Value(key)`
3. Pencarian nilai mengikuti rantai parent sampai ke root

## Mekanisme Pembatalan Context

### Cara Kerja Pembatalan

```go
func demonstrasiPembatalan() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go func() {
        select {
        case <-ctx.Done():
            fmt.Println("Goroutine dibatalkan:", ctx.Err())
            return
        case <-time.After(10 * time.Second):
            fmt.Println("Goroutine selesai normal")
        }
    }()
    
    // Simulasi pembatalan setelah 2 detik
    time.Sleep(2 * time.Second)
    cancel() // Channel Done() ditutup di sini
    
    time.Sleep(1 * time.Second)
}
```

**Alur pembatalan:**
1. `cancel()` dipanggil
2. Channel internal ditutup
3. Semua goroutine yang menunggu di `<-ctx.Done()` terbangun
4. `ctx.Err()` mengembalikan error pembatalan

### Pembatalan Cascade (Beruntun)

```go
func demonstrasiCascade() {
    // Parent context
    parentCtx, parentCancel := context.WithCancel(context.Background())
    defer parentCancel()
    
    // Child context 1
    childCtx1, cancel1 := context.WithCancel(parentCtx)
    defer cancel1()
    
    // Child context 2 dari child 1
    childCtx2, cancel2 := context.WithCancel(childCtx1)
    defer cancel2()
    
    go monitarContext(childCtx2, "Child 2")
    go monitarContext(childCtx1, "Child 1")
    go monitarContext(parentCtx, "Parent")
    
    time.Sleep(1 * time.Second)
    
    // Membatalkan parent akan membatalkan semua child
    parentCancel()
    
    time.Sleep(1 * time.Second)
}

func monitarContext(ctx context.Context, name string) {
    <-ctx.Done()
    fmt.Printf("%s dibatalkan: %v\n", name, ctx.Err())
}
```

## Skenario Penggunaan Real-World

### Skenario 1: HTTP Server dengan Timeout

**Problem:** Server perlu membatasi waktu pemrosesan request untuk mencegah resource exhaustion.

**Solusi dengan Context:**

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
func queryDatabase(ctx context.Context, userID int) (map[string]interface{}, error) {
    // Channel untuk mengembalikan hasil
    resultChan := make(chan map[string]interface{})
    errChan := make(chan error)
    
    go func() {
        // Simulasi query yang memakan waktu 3 detik
        time.Sleep(3 * time.Second)
        
        result := map[string]interface{}{
            "user_id": userID,
            "name":    "John Doe",
            "email":   "john@example.com",
        }
        
        resultChan <- result
    }()
    
    // Tunggu hasil atau context dibatalkan
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case result := <-resultChan:
        return result, nil
    case err := <-errChan:
        return nil, err
    }
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    // Buat context dengan timeout 2 detik
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()
    
    // Query database dengan context
    userData, err := queryDatabase(ctx, 123)
    
    if err != nil {
        if err == context.DeadlineExceeded {
            http.Error(w, "Request timeout", http.StatusRequestTimeout)
            log.Println("Database query timeout")
            return
        }
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userData)
}

func main() {
    http.HandleFunc("/user", userHandler)
    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**Penjelasan mekanisme:**
1. Request masuk ke handler
2. Context dengan timeout 2 detik dibuat
3. Database query dimulai (simulasi 3 detik)
4. Setelah 2 detik, context timeout
5. Select case menerima signal dari `ctx.Done()`
6. Error `DeadlineExceeded` dikembalikan
7. Client menerima response 408 Request Timeout

### Skenario 2: Worker Pool dengan Graceful Shutdown

**Problem:** Aplikasi perlu menghentikan semua worker dengan aman saat shutdown.

**Solusi dengan Context:**

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

type Job struct {
    ID   int
    Data string
}

type Worker struct {
    ID       int
    JobQueue chan Job
    wg       *sync.WaitGroup
}

func (w *Worker) Start(ctx context.Context) {
    defer w.wg.Done()
    
    fmt.Printf("Worker %d started\n", w.ID)
    
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d shutting down: %v\n", w.ID, ctx.Err())
            return
            
        case job, ok := <-w.JobQueue:
            if !ok {
                fmt.Printf("Worker %d: job queue closed\n", w.ID)
                return
            }
            
            fmt.Printf("Worker %d processing job %d: %s\n", w.ID, job.ID, job.Data)
            
            // Simulasi pemrosesan
            processingTime := time.Duration(rand.Intn(3)+1) * time.Second
            
            select {
            case <-ctx.Done():
                fmt.Printf("Worker %d: job %d interrupted\n", w.ID, job.ID)
                return
            case <-time.After(processingTime):
                fmt.Printf("Worker %d completed job %d\n", w.ID, job.ID)
            }
        }
    }
}

type WorkerPool struct {
    Workers   []*Worker
    JobQueue  chan Job
    ctx       context.Context
    cancel    context.CancelFunc
    wg        sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    jobQueue := make(chan Job, 100)
    
    pool := &WorkerPool{
        Workers:  make([]*Worker, numWorkers),
        JobQueue: jobQueue,
        ctx:      ctx,
        cancel:   cancel,
    }
    
    for i := 0; i < numWorkers; i++ {
        pool.Workers[i] = &Worker{
            ID:       i + 1,
            JobQueue: jobQueue,
            wg:       &pool.wg,
        }
    }
    
    return pool
}

func (p *WorkerPool) Start() {
    for _, worker := range p.Workers {
        p.wg.Add(1)
        go worker.Start(p.ctx)
    }
}

func (p *WorkerPool) Shutdown() {
    fmt.Println("\n=== Initiating graceful shutdown ===")
    
    // 1. Stop accepting new jobs
    close(p.JobQueue)
    
    // 2. Cancel context to signal all workers
    p.cancel()
    
    // 3. Wait for all workers to finish
    p.wg.Wait()
    
    fmt.Println("=== All workers stopped ===")
}

func (p *WorkerPool) AddJob(job Job) {
    select {
    case <-p.ctx.Done():
        fmt.Println("Cannot add job: pool is shutting down")
    case p.JobQueue <- job:
        fmt.Printf("Job %d added to queue\n", job.ID)
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    // Buat worker pool dengan 3 workers
    pool := NewWorkerPool(3)
    pool.Start()
    
    // Channel untuk menangkap signal OS
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    // Goroutine untuk menambahkan jobs
    go func() {
        for i := 1; i <= 10; i++ {
            job := Job{
                ID:   i,
                Data: fmt.Sprintf("Task %d", i),
            }
            pool.AddJob(job)
            time.Sleep(500 * time.Millisecond)
        }
    }()
    
    // Tunggu signal shutdown
    <-sigChan
    
    // Graceful shutdown
    pool.Shutdown()
}
```

**Penjelasan mekanisme:**
1. Worker pool dibuat dengan context yang dapat dibatalkan
2. Setiap worker mendengarkan `ctx.Done()` dan job queue
3. Saat SIGTERM/SIGINT diterima, `cancel()` dipanggil
4. Semua worker menerima signal pembatalan
5. Worker menyelesaikan job yang sedang diproses atau langsung berhenti
6. `wg.Wait()` memastikan semua worker selesai sebelum program exit

### Skenario 3: Multiple API Calls dengan Timeout

**Problem:** Aplikasi perlu memanggil beberapa API eksternal, tapi tidak mau menunggu terlalu lama.

**Solusi dengan Context:**

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type ServiceResponse struct {
    ServiceName string
    Data        interface{}
    Error       error
    Duration    time.Duration
}

// Simulasi service eksternal
func callExternalService(ctx context.Context, serviceName string, url string, delay time.Duration) ServiceResponse {
    start := time.Now()
    response := ServiceResponse{ServiceName: serviceName}
    
    // Simulasi delay network
    select {
    case <-ctx.Done():
        response.Error = ctx.Err()
        response.Duration = time.Since(start)
        return response
    case <-time.After(delay):
        // Continue with request
    }
    
    // Buat HTTP request dengan context
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        response.Error = err
        response.Duration = time.Since(start)
        return response
    }
    
    // Execute request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        response.Error = err
        response.Duration = time.Since(start)
        return response
    }
    defer resp.Body.Close()
    
    // Read response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        response.Error = err
        response.Duration = time.Since(start)
        return response
    }
    
    var data interface{}
    json.Unmarshal(body, &data)
    
    response.Data = data
    response.Duration = time.Since(start)
    return response
}

// Aggregate data dari multiple services
func aggregateServiceData(ctx context.Context) map[string]ServiceResponse {
    results := make(map[string]ServiceResponse)
    resultChan := make(chan ServiceResponse, 3)
    
    // Define services dengan delay yang berbeda
    services := []struct {
        name  string
        url   string
        delay time.Duration
    }{
        {"UserService", "https://jsonplaceholder.typicode.com/users/1", 500 * time.Millisecond},
        {"PostService", "https://jsonplaceholder.typicode.com/posts/1", 1 * time.Second},
        {"CommentService", "https://jsonplaceholder.typicode.com/comments/1", 1500 * time.Millisecond},
    }
    
    // Launch goroutine untuk setiap service
    for _, service := range services {
        go func(s struct {
            name  string
            url   string
            delay time.Duration
        }) {
            result := callExternalService(ctx, s.name, s.url, s.delay)
            resultChan <- result
        }(service)
    }
    
    // Collect results
    for i := 0; i < len(services); i++ {
        select {
        case <-ctx.Done():
            fmt.Println("Main context cancelled, stopping collection")
            return results
        case result := <-resultChan:
            results[result.ServiceName] = result
        }
    }
    
    return results
}

func main() {
    // Scenario 1: Timeout terlalu pendek
    fmt.Println("=== Scenario 1: Short Timeout (800ms) ===")
    ctx1, cancel1 := context.WithTimeout(context.Background(), 800*time.Millisecond)
    defer cancel1()
    
    results1 := aggregateServiceData(ctx1)
    for name, result := range results1 {
        if result.Error != nil {
            fmt.Printf("%s: ERROR - %v (took %v)\n", name, result.Error, result.Duration)
        } else {
            fmt.Printf("%s: SUCCESS (took %v)\n", name, result.Duration)
        }
    }
    
    time.Sleep(2 * time.Second)
    
    // Scenario 2: Timeout cukup
    fmt.Println("\n=== Scenario 2: Adequate Timeout (2s) ===")
    ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel2()
    
    results2 := aggregateServiceData(ctx2)
    for name, result := range results2 {
        if result.Error != nil {
            fmt.Printf("%s: ERROR - %v (took %v)\n", name, result.Error, result.Duration)
        } else {
            fmt.Printf("%s: SUCCESS (took %v)\n", name, result.Duration)
        }
    }
    
    // Scenario 3: Manual cancellation
    fmt.Println("\n=== Scenario 3: Manual Cancellation ===")
    ctx3, cancel3 := context.WithCancel(context.Background())
    
    go func() {
        time.Sleep(600 * time.Millisecond)
        fmt.Println("Cancelling manually...")
        cancel3()
    }()
    
    results3 := aggregateServiceData(ctx3)
    for name, result := range results3 {
        if result.Error != nil {
            fmt.Printf("%s: ERROR - %v (took %v)\n", name, result.Error, result.Duration)
        } else {
            fmt.Printf("%s: SUCCESS (took %v)\n", name, result.Duration)
        }
    }
}
```

**Penjelasan mekanisme:**
1. Context dengan timeout dibuat
2. Multiple goroutine diluncurkan untuk memanggil API
3. Setiap goroutine menerima context yang sama
4. Jika timeout tercapai, semua request yang belum selesai dibatalkan
5. HTTP client otomatis membatalkan request yang belum selesai
6. Results dikumpulkan dengan informasi sukses/gagal

### Skenario 4: Context Values untuk Request Tracing

**Problem:** Perlu melacak request ID melalui berbagai layer aplikasi untuk logging dan debugging.

**Solusi dengan Context:**

```go
package main

import (
    "context"
    "fmt"
    "log"
    "math/rand"
    "time"
)

// Key types untuk context values
type contextKey string

const (
    requestIDKey contextKey = "requestID"
    userIDKey    contextKey = "userID"
    sessionIDKey contextKey = "sessionID"
)

// Helper functions untuk set/get values
func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, requestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
    if requestID, ok := ctx.Value(requestIDKey).(string); ok {
        return requestID
    }
    return "unknown"
}

func WithUserID(ctx context.Context, userID int) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) int {
    if userID, ok := ctx.Value(userIDKey).(int); ok {
        return userID
    }
    return 0
}

// Logger yang aware terhadap context
type ContextLogger struct{}

func (l *ContextLogger) Info(ctx context.Context, message string) {
    requestID := GetRequestID(ctx)
    userID := GetUserID(ctx)
    log.Printf("[INFO] [RequestID: %s] [UserID: %d] %s", requestID, userID, message)
}

func (l *ContextLogger) Error(ctx context.Context, message string, err error) {
    requestID := GetRequestID(ctx)
    userID := GetUserID(ctx)
    log.Printf("[ERROR] [RequestID: %s] [UserID: %d] %s: %v", requestID, userID, message, err)
}

// Simulasi layers aplikasi
func controllerLayer(ctx context.Context) error {
    logger := &ContextLogger{}
    logger.Info(ctx, "Controller: Processing request")
    
    // Pass context ke service layer
    return serviceLayer(ctx)
}

func serviceLayer(ctx context.Context) error {
    logger := &ContextLogger{}
    logger.Info(ctx, "Service: Executing business logic")
    
    // Simulasi proses
    time.Sleep(100 * time.Millisecond)
    
    // Pass context ke repository layer
    return repositoryLayer(ctx)
}

func repositoryLayer(ctx context.Context) error {
    logger := &ContextLogger{}
    logger.Info(ctx, "Repository: Querying database")
    
    // Simulasi database query
    time.Sleep(200 * time.Millisecond)
    
    // Simulasi error random
    if rand.Intn(3) == 0 {
        err := fmt.Errorf("database connection failed")
        logger.Error(ctx, "Repository: Database error", err)
        return err
    }
    
    logger.Info(ctx, "Repository: Query successful")
    return nil
}

// Middleware untuk generate request ID
func requestIDMiddleware(next func(context.Context) error) func(context.Context) error {
    return func(ctx context.Context) error {
        requestID := fmt.Sprintf("REQ-%d", rand.Intn(10000))
        ctx = WithRequestID(ctx, requestID)
        
        return next(ctx)
    }
}

// Middleware untuk extract user ID
func authMiddleware(userID int) func(func(context.Context) error) func(context.Context) error {
    return func(next func(context.Context) error) func(context.Context) error {
        return func(ctx context.Context) error {
            ctx = WithUserID(ctx, userID)
            return next(ctx)
        }
    }
}

func handleRequest(ctx context.Context) error {
    logger := &ContextLogger{}
    logger.Info(ctx, "Starting request processing")
    
    // Process through layers
    err := controllerLayer(ctx)
    
    if err != nil {
        logger.Error(ctx, "Request failed", err)
        return err
    }
    
    logger.Info(ctx, "Request completed successfully")
    return nil
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    // Simulasi beberapa requests
    for i := 1; i <= 5; i++ {
        fmt.Printf("\n=== Processing Request %d ===\n", i)
        
        ctx := context.Background()
        
        // Apply middlewares
        handler := authMiddleware(100 + i)(requestIDMiddleware(handleRequest))
        
        // Execute request
        handler(ctx)
        
        time.Sleep(500 * time.Millisecond)
    }
}
```

**Penjelasan mekanisme:**
1. Request ID dan User ID disimpan dalam context
2. Context dipass melalui semua layer aplikasi
3. Setiap layer dapat mengakses metadata tanpa parameter tambahan
4. Logger menggunakan context untuk menambahkan informasi tracing
5. Memudahkan debugging karena semua log terkait request memiliki ID yang sama

### Skenario 5: Database Transaction dengan Context

**Problem:** Database transaction perlu dibatalkan jika request timeout atau user membatalkan.

**Solusi dengan Context:**

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"
    
    _ "github.com/mattn/go-sqlite3"
)

type User struct {
    ID    int
    Name  string
    Email string
}

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

// Create user dengan context
func (r *UserRepository) CreateUser(ctx context.Context, user User) error {
    query := "INSERT INTO users (name, email) VALUES (?, ?)"
    
    // Execute query dengan context
    result, err := r.db.ExecContext(ctx, query, user.Name, user.Email)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get last insert id: %w", err)
    }
    
    log.Printf("User created with ID: %d", id)
    return nil
}

// Complex transaction dengan multiple operations
func (r *UserRepository) CreateUserWithProfile(ctx context.Context, user User, profileData map[string]string) error {
    // Start transaction dengan context
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    
    // Defer rollback jika terjadi error
    defer func() {
        if err != nil {
            tx.Rollback()
            log.Println("Transaction rolled back")
        }
    }()
    
    // Step 1: Insert user
    result, err := tx.ExecContext(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
    if err != nil {
        return fmt.Errorf("failed to insert user: %w", err)
    }
    
    userID, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get user id: %w", err)
    }
    
    log.Printf("User inserted with ID: %d", userID)
    
    // Simulasi operasi yang memakan waktu
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(2 * time.Second):
        // Continue
    }
    
    // Step 2: Insert profile
    _, err = tx.ExecContext(ctx, 
        "INSERT INTO profiles (user_id, bio, location) VALUES (?, ?, ?)",
        userID, profileData["bio"], profileData["location"])
    if err != nil {
        return fmt.Errorf("failed to insert profile: %w", err)
    }
    
    log.Printf("Profile inserted for user ID: %d", userID)
    
    // Step 3: Insert preferences
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(1 * time.Second):
        // Continue
    }
    
    _, err = tx.ExecContext(ctx,
        "INSERT INTO preferences (user_id, theme, language) VALUES (?, ?, ?)",
        userID, "dark", "en")
    if err != nil {
        return fmt.Errorf("failed to insert preferences: %w", err)
    }
    
    log.Printf("Preferences inserted for user ID: %d", userID)
    
    // Commit transaction
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    log.Println("Transaction committed successfully")
    return nil
}

// Query dengan timeout
func (r *UserRepository) FindUserByID(ctx context.Context, id int) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE id = ?"
    
    row := r.db.QueryRowContext(ctx, query, id)
    
    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to scan user: %w", err)
    }
    
    return &user, nil
}

func setupDatabase() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        return nil, err
    }
    
    // Create tables
    _, err = db.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL
        );
        
        CREATE TABLE profiles (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            bio TEXT,
            location TEXT,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
        
        CREATE TABLE preferences (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            theme TEXT,
            language TEXT,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
    `)
    
    return db, err
}

func main() {
    db, err := setupDatabase()
    if err != nil {
        log.Fatal("Failed to setup database:", err)
    }
    defer db.Close()
    
    repo := NewUserRepository(db)
    
    // Scenario 1: Transaction sukses
    fmt.Println("\n=== Scenario 1: Successful Transaction ===")
    ctx1 := context.Background()
    user1 := User{Name: "John Doe", Email: "john@example.com"}
    profileData1 := map[string]string{
        "bio":      "Software Engineer",
        "location": "San Francisco",
    }
    
    err = repo.CreateUserWithProfile(ctx1, user1, profileData1)
    if err != nil {
        log.Printf("Error: %v", err)
    }
    
    // Scenario 2: Transaction timeout
    fmt.Println("\n=== Scenario 2: Transaction Timeout ===")
    ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel2()
    
    user2 := User{Name: "Jane Smith", Email: "jane@example.com"}
    profileData2 := map[string]string{
        "bio":      "Product Manager",
        "location": "New York",
    }
    
    err = repo.CreateUserWithProfile(ctx2, user2, profileData2)
    if err != nil {
        log.Printf("Error: %v", err)
    }
    
    // Scenario 3: Manual cancellation
    fmt.Println("\n=== Scenario 3: Manual Cancellation ===")
    ctx3, cancel3 := context.WithCancel(context.Background())
    
    go func() {
        time.Sleep(1500 * time.Millisecond)
        log.Println("Cancelling transaction manually...")
        cancel3()
    }()
    
    user3 := User{Name: "Bob Wilson", Email: "bob@example.com"}
    profileData3 := map[string]string{
        "bio":      "Designer",
        "location": "Los Angeles",
    }
    
    err = repo.CreateUserWithProfile(ctx3, user3, profileData3)
    if err != nil {
        log.Printf("Error: %v", err)
    }
    
    // Verify data
    fmt.Println("\n=== Verifying Data ===")
    queryCtx := context.Background()
    
    for id := 1; id <= 3; id++ {
        user, err := repo.FindUserByID(queryCtx, id)
        if err != nil {
            log.Printf("User ID %d: %v", id, err)
        } else {
            log.Printf("User ID %d: %s (%s)", user.ID, user.Name, user.Email)
        }
    }
}
```

**Penjelasan mekanisme:**
1. Database transaction dimulai dengan `BeginTx(ctx, nil)`
2. Setiap query menggunakan `*Context` variant (ExecContext, QueryRowContext)
3. Jika context timeout atau dibatalkan, query otomatis dihentikan
4. Transaction di-rollback jika ada error
5. Mencegah resource leak dengan defer rollback

## Best Practices Context

### 1. Selalu Pass Context sebagai Parameter Pertama

```go
// BENAR
func ProcessData(ctx context.Context, data string) error {
    // ...
}

// SALAH - context harus parameter pertama
func ProcessData(data string, ctx context.Context) error {
    // ...
}
```

### 2. Jangan Simpan Context dalam Struct

```go
// SALAH - jangan simpan context dalam struct
type Server struct {
    ctx context.Context
}

// BENAR - pass context sebagai parameter method
type Server struct {
    // fields lain
}

func (s *Server) ProcessRequest(ctx context.Context) error {
    // ...
}
```

**Alasan:** Context adalah request-scoped dan memiliki lifetime terbatas. Menyimpannya dalam struct dapat menyebabkan kebingungan tentang lifecycle-nya.

### 3. Selalu Panggil Cancel Function

```go
// BENAR
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel() // Penting untuk mencegah resource leak

// SALAH - tidak memanggil cancel
ctx, _ := context.WithTimeout(parent, 5*time.Second)
```

**Alasan:** Tidak memanggil cancel() dapat menyebabkan goroutine leak dan memory leak.

### 4. Gunakan Type-Safe Keys untuk Context Values

```go
// BENAR - gunakan custom type
type contextKey string
const userIDKey contextKey = "userID"

ctx = context.WithValue(ctx, userIDKey, 12345)

// SALAH - gunakan string literal
ctx = context.WithValue(ctx, "userID", 12345) // Bisa collision dengan package lain
```

### 5. Context Values untuk Request-Scoped Data Saja

```go
// BENAR - data yang spesifik untuk request
ctx = context.WithValue(ctx, requestIDKey, "REQ-12345")
ctx = context.WithValue(ctx, userIDKey, user.ID)

// SALAH - gunakan untuk parameter function biasa
ctx = context.WithValue(ctx, "config", appConfig) // Seharusnya parameter
ctx = context.WithValue(ctx, "logger", logger)    // Seharusnya dependency injection
```

### 6. Handle Context Cancellation di Loop

```go
// BENAR
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Process item
    }
}

// ATAU lebih efisien dengan case lain
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case item := <-itemChan:
        // Process item
    }
}
```

### 7. Propagate Context ke Semua Blocking Operations

```go
// BENAR
func FetchData(ctx context.Context, url string) error {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := client.Do(req)
    // ...
}

// SALAH - tidak menggunakan context
func FetchData(ctx context.Context, url string) error {
    resp, err := http.Get(url) // Context tidak di-propagate
    // ...
}
```

### 8. Jangan Gunakan nil Context

```go
// BENAR
ctx := context.Background()
ProcessData(ctx, data)

// SALAH
ProcessData(nil, data) // Bisa panic
```

## Common Pitfalls dan Solusinya

### Pitfall 1: Goroutine Leak

**Problem:**

```go
// SALAH - goroutine akan leak
func leakyFunction() {
    ctx := context.Background()
    ch := make(chan int)
    
    go func() {
        // Goroutine ini tidak pernah berhenti
        for {
            ch <- rand.Intn(100)
            time.Sleep(1 * time.Second)
        }
    }()
    
    // Function return, tapi goroutine masih jalan
}
```

**Solusi:**

```go
// BENAR
func fixedFunction() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // Pastikan goroutine dibersihkan
    
    ch := make(chan int)
    
    go func() {
        for {
            select {
            case <-ctx.Done():
                return // Goroutine berhenti saat context dibatalkan
            case ch <- rand.Intn(100):
                time.Sleep(1 * time.Second)
            }
        }
    }()
    
    // Do work...
}
```

### Pitfall 2: Tidak Memeriksa Context Done

**Problem:**

```go
// SALAH - tidak memeriksa context
func processItems(ctx context.Context, items []string) {
    for _, item := range items {
        // Proses bisa berlanjut meskipun context sudah cancelled
        heavyProcess(item)
    }
}
```

**Solusi:**

```go
// BENAR
func processItems(ctx context.Context, items []string) error {
    for _, item := range items {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := heavyProcess(item); err != nil {
                return err
            }
        }
    }
    return nil
}
```

### Pitfall 3: Context Value Type Assertion Tanpa Check

**Problem:**

```go
// SALAH - bisa panic jika type assertion gagal
func getUserID(ctx context.Context) int {
    return ctx.Value(userIDKey).(int) // Panic jika bukan int atau nil
}
```

**Solusi:**

```go
// BENAR
func getUserID(ctx context.Context) (int, bool) {
    userID, ok := ctx.Value(userIDKey).(int)
    return userID, ok
}

// ATAU dengan default value
func getUserIDOrDefault(ctx context.Context) int {
    if userID, ok := ctx.Value(userIDKey).(int); ok {
        return userID
    }
    return 0 // default value
}
```

### Pitfall 4: Deadline Lebih Panjang dari Parent

**Problem:**

```go
// Ini tidak akan bekerja seperti yang diharapkan
parentCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
childCtx, _ := context.WithTimeout(parentCtx, 5*time.Second) // Child akan cancelled setelah 2 detik, bukan 5
```

**Penjelasan:** Child context akan dibatalkan ketika parent-nya dibatalkan, jadi deadline efektif adalah yang paling pendek.

### Pitfall 5: Menggunakan Context.Background() di Dalam Request Handler

**Problem:**

```go
// SALAH - membuat context baru, kehilangan cancellation dari HTTP request
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background() // Don't do this!
    processRequest(ctx)
}
```

**Solusi:**

```go
// BENAR - gunakan context dari request
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    processRequest(ctx)
}
```

## Performance Considerations

### 1. Context Overhead

Context memiliki overhead minimal, tetapi tetap ada:
- Allocation untuk setiap WithValue, WithCancel, WithTimeout
- Lookup time untuk Value() yang traverse parent chain

**Optimasi:**

```go
// Jika perlu banyak values, lebih baik gunakan struct
type RequestMetadata struct {
    RequestID string
    UserID    int
    SessionID string
}

ctx = context.WithValue(ctx, metadataKey, RequestMetadata{...})

// Daripada
ctx = context.WithValue(ctx, requestIDKey, "...")
ctx = context.WithValue(ctx, userIDKey, 123)
ctx = context.WithValue(ctx, sessionIDKey, "...")
```

### 2. Channel Overhead

Done() channel memiliki cost untuk creation dan monitoring:

```go
// Untuk hot path, check context sesekali saja
func processMany(ctx context.Context, items []int) error {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for i, item := range items {
        // Check context setiap 100ms atau setiap 1000 items
        if i%1000 == 0 {
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-ticker.C:
                // Time to check
            default:
            }
        }
        
        process(item)
    }
    return nil
}
```

## Context dalam Testing

### Testing dengan Context

```go
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()
    
    err := longRunningOperation(ctx)
    
    if err != context.DeadlineExceeded {
        t.Errorf("Expected DeadlineExceeded, got %v", err)
    }
}

func TestWithCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    
    resultChan := make(chan error)
    go func() {
        resultChan <- operation(ctx)
    }()
    
    // Cancel after 50ms
    time.Sleep(50 * time.Millisecond)
    cancel()
    
    err := <-resultChan
    if err != context.Canceled {
        t.Errorf("Expected Canceled, got %v", err)
    }
}
```

## Kesimpulan

Context adalah tool yang sangat powerful untuk:
1. **Cancellation Propagation** - Membatalkan operasi secara hierarkis
2. **Deadline Management** - Membatasi waktu eksekusi
3. **Request-Scoped Values** - Membawa metadata request
4. **Graceful Shutdown** - Menghentikan services dengan aman
5. **Resource Management** - Mencegah goroutine dan resource leaks

**Key Takeaways:**
- Selalu pass context sebagai parameter pertama
- Selalu panggil cancel function dengan defer
- Gunakan context untuk cancellation dan deadlines, bukan untuk passing dependencies
- Check ctx.Done() di long-running operations
- Jangan simpan context dalam struct
- Gunakan type-safe keys untuk context values

Context adalah fondasi dari concurrent programming yang baik di Go. Memahami dan menggunakan context dengan benar akan membuat aplikasi Anda lebih robust, maintainable, dan production-ready.