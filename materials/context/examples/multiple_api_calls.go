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