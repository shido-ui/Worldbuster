package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","service":"worldbuster"}`))
	})
	log.Println("Worldbuster backend listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
