package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"backend/internal/models"
)

func (deps *RouterDeps) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("user-id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	transaction.UserID = userID
	transaction.InsertedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	transactionID, err := deps.Repo.AddTransaction(context.Background(), transaction)
	if err != nil {
		log.Printf("Error adding transaction to DB: %v", err)
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	transaction.ID = transactionID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(transaction)
	if err != nil {
		return
	}
}
