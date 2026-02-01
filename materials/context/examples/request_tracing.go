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