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

	// Create server
	srv := createServer()

	// Channel to listen errors
	serverErrors := make(chan error, 1)

	// Start server
	go func() {
		log.Printf("Server starting on %s...\n", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Wait
	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)

	case sig := <-shutdown:
		log.Printf("Received signal %v: initiating shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
		} else {
			log.Println("Server shut down gracefully.")
		}
	}
}

func createServer() *http.Server {
	r := mux.NewRouter()

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://bugtracker-staging-jameswillett.fly.dev",
			"https://bugtracker-jameswillett.fly.dev",
			"https://your-netlify-frontend-domain.netlify.app",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// Routes
	r.HandleFunc("/api/health", handlers.HealthCheck).Methods("GET")
	apiRouter := r.PathPrefix("/api").Subrouter()
	handlers.RegisterRoutes(apiRouter)

	// REQUIRED FOR RENDER: read port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local development
	}

	addr := "0.0.0.0:" + port
	log.Printf("Server binding to %s", addr)

	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}
