package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "bugtracker-backend/internal/db"
    "bugtracker-backend/internal/handlers"

    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    log.Println("Starting Bug Tracker backend server...")

    if err := db.Init(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Cleanup()

    srv := createServer()

    serverErrors := make(chan error, 1)

    go func() {
        log.Printf("Server starting on %s...\n", srv.Addr)
        serverErrors <- srv.ListenAndServe()
    }()

    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

    select {
    case err := <-serverErrors:
        log.Fatalf("server error: %v", err)

    case sig := <-shutdown:
        log.Printf("Received signal %v: initiating shutdown", sig)

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        if err := srv.Shutdown(ctx); err != nil {
            log.Printf("Graceful shutdown timed out: %v", err)
        } else {
            log.Println("Server shut down gracefully.")
        }
    }
}

func createServer() *http.Server {
    port := os.Getenv("PORT")
    if port == "" {
        port = "10000"
    }

    address := fmt.Sprintf("0.0.0.0:%s", port)
    log.Printf("Server binding to %s", address)

    r := mux.NewRouter()

    // Required for Render health check
    r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

    // Your API
    r.HandleFunc("/api/health", handlers.HealthCheck).Methods("GET")
    api := r.PathPrefix("/api").Subrouter()
    handlers.RegisterRoutes(api)

    // CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"*"},
        AllowCredentials: true,
    })

    handler := c.Handler(r)

    return &http.Server{
        Addr:    address,
        Handler: handler,
    }
}
