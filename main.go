package main

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
	redoc "github.com/go-openapi/runtime/middleware"

	"github.com/identitatem/idp-configs-api/config"
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

func initDependencies() {
	config.Init()
	// l.InitLogger()
	// db.InitDB()
}

func main() {
	initDependencies()
	cfg := config.Get()

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
		http.ServeFile(w, r, cfg.OpenAPIFilePath)
	})

	fmt.Println("Listening on port", cfg.WebPort)

	// Listen and serve using the router
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.WebPort), r)
}
