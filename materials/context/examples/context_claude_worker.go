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