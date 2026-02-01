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