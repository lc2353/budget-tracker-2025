package db

import (
	"backend/internal/models"
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"google.golang.org/api/option"
)

type Repository interface {
	AddTransaction(ctx context.Context, transaction models.Transaction) (string, error)
	GetTransactionByID(ctx context.Context, userID, transactionID string) (*models.Transaction, error)
	ListTransactions(ctx context.Context, userID string, filters map[string]string) ([]models.Transaction, error)
	UpdateTransaction(ctx context.Context, userID, transactionID string, updates map[string]interface{}) error
	BulkAddTransactions(ctx context.Context, transactions []models.Transaction) error
}

type FirestoreRepository struct {
	client *firestore.Client
}

func (r *FirestoreRepository) GetTransactionByID(ctx context.Context, userID, transactionID string) (*models.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (r *FirestoreRepository) ListTransactions(ctx context.Context, userID string, filters map[string]string) ([]models.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (r *FirestoreRepository) UpdateTransaction(ctx context.Context, userID, transactionID string, updates map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (r *FirestoreRepository) BulkAddTransactions(ctx context.Context, transactions []models.Transaction) error {
	//TODO implement me
	panic("implement me")
}

func NewFirestoreClient(ctx context.Context, credentialsPath string) (*firestore.Client, error) {
	var opts []option.ClientOption
	if credentialsPath != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsPath))
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, fmt.Errorf("firebase.NewApp: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("app.Firestore: %w", err)
	}
	return client, nil
}

func NewFirestoreRepository(client *firestore.Client) *FirestoreRepository {
	return &FirestoreRepository{client: client}
}

func (r *FirestoreRepository) AddTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	// For nested collection: users/{userID}/transactions
	ref, _, err := r.client.Collection("transactions").Add(ctx, transaction)
	if err != nil {
		return "", fmt.Errorf("failed to add transaction: %w", err)
	}
	return ref.ID, nil
}
