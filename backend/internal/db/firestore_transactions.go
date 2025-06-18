package db

import (
	"backend/internal/exceptions"
	"backend/internal/models"
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Repository interface {
	AddTransaction(ctx context.Context, transaction models.Transaction) (string, error)
	GetTransactionByID(ctx context.Context, userID, transactionID string) (*models.Transaction, error)
	ListTransactions(ctx context.Context, userID string, filters map[string]string) ([]models.Transaction, error)
	BulkAddTransactions(ctx context.Context, transactions []models.Transaction) ([]models.Transaction, error)
	UpdateTransaction(ctx context.Context, userID, transactionID string, updateData models.TransactionUpdate) (*models.Transaction, error)
	DeleteTransaction(ctx context.Context, userID, transactionID string) error
}

type FirestoreRepository struct {
	client *firestore.Client
}

func (r *FirestoreRepository) AddTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	ref, _, err := r.client.Collection("transactions").Add(ctx, transaction)
	if err != nil {
		return "", fmt.Errorf("failed to add transaction: %w", err)
	}
	return ref.ID, nil
}

func (r *FirestoreRepository) GetTransactionByID(ctx context.Context, userID, transactionID string) (*models.Transaction, error) {
	docRef := r.client.Collection("transactions").Doc(transactionID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, exceptions.TransactionNotFound(transactionID)
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	var transaction models.Transaction
	if err := doc.DataTo(&transaction); err != nil {
		return nil, fmt.Errorf("failed to convert document data to transaction: %w", err)
	}
	if transaction.UserID != userID {
		return nil, exceptions.UserForbidden(userID)
	}

	transaction.ID = doc.Ref.ID
	transaction.UserID = ""
	return &transaction, nil
}

func (r *FirestoreRepository) ListTransactions(ctx context.Context, userID string, filters map[string]string) ([]models.Transaction, error) {
	query := r.client.Collection("transactions").Where("userId", "==", userID)

	for field, value := range filters {
		query = query.Where(field, "==", value)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	var transactions []models.Transaction
	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				return transactions, nil
			}
			if status.Code(err) == codes.NotFound {
				return transactions, nil
			}
			return nil, fmt.Errorf("failed to list transactions: %w", err)
		}

		var transaction models.Transaction
		if err := doc.DataTo(&transaction); err != nil {
			return nil, fmt.Errorf("failed to convert document data to transaction: %w", err)
		}
		transaction.ID = doc.Ref.ID
		transaction.UserID = ""
		transactions = append(transactions, transaction)
	}
}

func (r *FirestoreRepository) BulkAddTransactions(ctx context.Context, transactions []models.Transaction) ([]models.Transaction, error) {
	var docRefs []*firestore.DocumentRef
	for _, transaction := range transactions {
		docRef := r.client.Collection("transactions").NewDoc()
		_, err := docRef.Set(ctx, transaction)
		if err != nil {
			return nil, fmt.Errorf("failed to add transaction: %w", err)
		}
		docRefs = append(docRefs, docRef)
	}
	for i, docRef := range docRefs {
		transactions[i].ID = docRef.ID
		transactions[i].UserID = ""
	}
	return transactions, nil
}

func (r *FirestoreRepository) UpdateTransaction(ctx context.Context, userID, transactionID string, updateData models.TransactionUpdate) (*models.Transaction, error) {
	docRef := r.client.Collection("transactions").Doc(transactionID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, exceptions.TransactionNotFound(transactionID)
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	var transaction models.Transaction
	if err := doc.DataTo(&transaction); err != nil {
		return nil, fmt.Errorf("failed to convert document data to transaction: %w", err)
	}

	if transaction.UserID != userID {
		return nil, exceptions.UserForbidden(userID)
	}

	updateMap := toUpdateMap(updateData)
	updateMap["updatedAt"] = time.Now()

	_, err = docRef.Set(ctx, updateMap, firestore.MergeAll)
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	doc, err = docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated transaction: %w", err)
	}
	if err := doc.DataTo(&transaction); err != nil {
		return nil, fmt.Errorf("failed to convert updated document data to transaction: %w", err)
	}

	transaction.ID = doc.Ref.ID
	transaction.UserID = ""
	return &transaction, nil
}

func (r *FirestoreRepository) DeleteTransaction(ctx context.Context, userID, transactionID string) error {
	docRef := r.client.Collection("transactions").Doc(transactionID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return exceptions.TransactionNotFound(transactionID)
		}
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	var transaction models.Transaction
	if err := doc.DataTo(&transaction); err != nil {
		return fmt.Errorf("failed to convert document data to transaction: %w", err)
	}

	if transaction.UserID != userID {
		return exceptions.UserForbidden(userID)
	}

	_, err = docRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	return nil
}

func toUpdateMap(update models.TransactionUpdate) map[string]interface{} {
	result := make(map[string]interface{})
	if update.Description != nil {
		result["description"] = *update.Description
	}
	if update.Amount != nil {
		result["amount"] = *update.Amount
	}
	if update.Category != nil {
		result["category"] = *update.Category
	}
	if update.Type != nil {
		result["type"] = *update.Type
	}
	return result
}
