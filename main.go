package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	chi "github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
	redoc "github.com/go-openapi/runtime/middleware"
)

func setupDocsMiddleware(handler http.Handler) http.Handler {
	opt := redoc.RedocOpts{
		SpecURL: "/api/hello-world-service/v0/openapi.json",
	}
	return redoc.Redoc(opt, handler)
}

// Handler function that responds with Hello World
func helloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello world")
}

func main() {
    // Initialize router
	r := chi.NewRouter()
	
	r.Use(
		middleware.Logger,
		setupDocsMiddleware,
	)

    // Register handler functions on server routes
	r.Get("/api/hello-world-service/v0/ping", helloWorld)

	// OpenAPI Spec
	r.Get("/api/hello-world-service/v0/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./cmd/spec/openapi.json")
	})

	// Initialize server
	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Handle graceful shutdown of server
    idleConnsClosed := make(chan struct{})
    go func() {
        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, os.Interrupt)
        <-sigint

        // We received an interrupt signal, shut down.
        if err := srv.Shutdown(context.Background()); err != nil {
            // Error from closing listeners, or context timeout:
            log.Printf("HTTP server Shutdown: %v", err)
        } else {
			log.Printf("HTTP server successfully shutdown")
		}
        close(idleConnsClosed)
    }()

    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        // Error starting or closing listener:
        log.Printf("HTTP server ListenAndServe: %v", err)
    }

    <-idleConnsClosed	
}
