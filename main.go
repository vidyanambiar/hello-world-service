package main

import (
	"fmt"
	"net/http"
)

// Handler function that responds with Hello World
func helloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello world")
}

func main() {
	// Register handler function on server route
    http.HandleFunc("/", helloWorld)
	
    fmt.Println("Listening on localhost:8080")
    http.ListenAndServe(":8080", nil)
}
