package main

import (
	"log"
	"net/http"
	"time"

	"github.com/shido-ui/Worldbuster/backend/internal/api"
	"github.com/shido-ui/Worldbuster/backend/internal/world"
)

func main() {
	worldService := world.NewService()
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			worldService.Tick()
		}
	}()

	server := &http.Server{
		Addr:              ":8080",
		Handler:           api.NewRouter(worldService).Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("Worldbuster server listening on :8080")
	log.Fatal(server.ListenAndServe())
}
