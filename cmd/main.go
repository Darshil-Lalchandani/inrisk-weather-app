package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"example.com/app"
	"github.com/go-chi/chi"
)

func main() {
	// create new GCS client and pass it to manager as dependency
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create GCS client: %w", err)
	}

	manager := app.NewManager(client)

	controller := app.NewController(manager)

	r := chi.NewRouter()

	controller.MountRoutes(r)

	port := "8080"

	// Start the HTTP server
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
