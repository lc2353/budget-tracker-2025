package main

import (
	"backend/internal/api"
	"backend/internal/db"
	config "backend/internal/setup"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	ctx := context.Background()
	firestoreClient, err := db.NewFirestoreClient(ctx, cfg.FirestoreCredentialsPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore client: %v", err)
	}
	defer func() {
		if err := firestoreClient.Close(); err != nil {
			log.Printf("Error closing Firestore client: %v", err)
		}
		log.Println("Firestore client closed.")
	}()
	log.Println("Firestore client initialized successfully!")

	repo := db.NewFirestoreRepository(firestoreClient)

	router := api.NewRouter(repo, cfg)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
