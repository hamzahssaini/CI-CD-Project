package main

import (
    "fmt"
    "net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "✅ Go service healthy")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "👋 Hello from Go Service!")
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/health", healthHandler)
    fmt.Println("Go service running on port 3002")
    http.ListenAndServe(":3002", nil)
}
