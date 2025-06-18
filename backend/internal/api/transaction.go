package api

import (
	"backend/internal/exceptions"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
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

	userID, ok := GetUserIDFromHeader(w, r)
	if !ok {
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
	transaction.UserID = ""
	EncodeJSONResponse(w, transaction)
}

func (deps *RouterDeps) GetTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]
	if transactionID == "" {
		http.Error(w, "Missing transaction ID", http.StatusBadRequest)
		return
	}

	userID, ok := GetUserIDFromHeader(w, r)
	if !ok {
		return
	}

	transaction, err := deps.Repo.GetTransactionByID(context.Background(), userID, transactionID)
	if err != nil {
		var forbiddenErr *exceptions.UserForbiddenError
		if errors.As(err, &forbiddenErr) {
			log.Print(exceptions.UserForbidden(userID))
			http.Error(w, "", http.StatusForbidden)
			return
		}
		var notFoundErr *exceptions.TransactionNotFoundError
		if errors.As(err, &notFoundErr) {
			log.Print(exceptions.TransactionNotFound(transactionID))
			http.Error(w, "transaction not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting transaction: %v", err)
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	EncodeJSONResponse(w, transaction)

}

func (deps *RouterDeps) ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := GetUserIDFromHeader(w, r)
	if !ok {
		return
	}

	transactions, err := deps.Repo.ListTransactions(context.Background(), userID, nil)
	if err != nil {
		log.Printf("Error listing transactions: %v", err)
		http.Error(w, "Failed to list transactions", http.StatusInternalServerError)
		return
	}

	EncodeJSONResponse(w, transactions)
}

func (deps *RouterDeps) BulkAddTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	var transactions []models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transactions); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	userID, ok := GetUserIDFromHeader(w, r)
	if !ok {
		return
	}

	for i := range transactions {
		transactions[i].UserID = userID
		transactions[i].InsertedAt = time.Now()
		transactions[i].UpdatedAt = time.Now()
	}

	transactions, err := deps.Repo.BulkAddTransactions(context.Background(), transactions)
	if err != nil {
		log.Printf("Error bulk adding transactions: %v", err)
		http.Error(w, "Failed to bulk add transactions", http.StatusInternalServerError)
		return
	}

	EncodeJSONResponse(w, transactions)
}

func (deps *RouterDeps) UpdateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]
	if transactionID == "" {
		http.Error(w, "Missing transaction ID", http.StatusBadRequest)
		return
	}

	userID, ok := GetUserIDFromHeader(w, r)
	if !ok {
		return
	}

	var updateData models.TransactionUpdate
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	transaction, err := deps.Repo.UpdateTransaction(context.Background(), userID, transactionID, updateData)
	if err != nil {
		var forbiddenErr *exceptions.UserForbiddenError
		if errors.As(err, &forbiddenErr) {
			log.Print(exceptions.UserForbidden(userID))
			http.Error(w, "", http.StatusForbidden)
			return
		}
		var notFoundErr *exceptions.TransactionNotFoundError
		if errors.As(err, &notFoundErr) {
			log.Print(exceptions.TransactionNotFound(transactionID))
			http.Error(w, "transaction not found", http.StatusNotFound)
			return
		}
		log.Printf("Error updating transaction: %v", err)
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}

	EncodeJSONResponse(w, transaction)
}

func (deps *RouterDeps) DeleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]
	if transactionID == "" {
		http.Error(w, "Missing transaction ID", http.StatusBadRequest)
		return
	}

	userID, ok := GetUserIDFromHeader(w, r)
	if !ok {
		return
	}

	if err := deps.Repo.DeleteTransaction(context.Background(), userID, transactionID); err != nil {
		var forbiddenErr *exceptions.UserForbiddenError
		if errors.As(err, &forbiddenErr) {
			log.Print(exceptions.UserForbidden(userID))
			http.Error(w, "", http.StatusForbidden)
			return
		}
		var notFoundErr *exceptions.TransactionNotFoundError
		if errors.As(err, &notFoundErr) {
			log.Print(exceptions.TransactionNotFound(transactionID))
			http.Error(w, "transaction not found", http.StatusNotFound)
			return
		}
		log.Printf("Error deleting transaction: %v", err)
		http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetUserIDFromHeader(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID := r.Header.Get("user-id")
	if userID == "" {
		http.Error(w, "Missing 'user-id' header", http.StatusBadRequest)
		return "", false
	}
	return userID, true
}

func EncodeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
