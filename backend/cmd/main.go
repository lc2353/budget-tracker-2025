package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"backend/internal"
)

var firestoreClient *firestore.Client

func main() {
	cfg := config.LoadConfig()

	ctx := context.Background()
	var opts []option.ClientOption

	if cfg.FirestoreCredentialPath != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.FirestoreCredentialPath))
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	firestoreClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v", err)
	}

	log.Println("Firestore client initialized successfully")

	testDoc := map[string]interface{}{
		"transactionName": "coffee",
		"description":     "starbucks latte",
		"amount":          4.76,
		"category":        "Coffee",
		"insertedAt":      time.Now(),
	}

	_, _, err = firestoreClient.Collection("test-txns").Add(ctx, testDoc)
	if err != nil {
		log.Printf("Failed to add setup test document: %v", err)
	} else {
		log.Println("Setup test document added to 'test-txns' collection.")
	}
}
