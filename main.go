package main

import (
	"fmt"
	"net/http"

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

    fmt.Println("Listening on localhost:8080")

    // Listen and serve using the router
    http.ListenAndServe(":8080", r)
}
