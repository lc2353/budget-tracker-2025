package models

import (
	"time"
)

type Transaction struct {
	ID                  string    `json:"id" firestore:"-"`
	UserID              string    `json:"userId" firestore:"userId"`
	TransactionDateTime time.Time `json:"transactionDateTime" firestore:"transactionDateTime"`
	Description         string    `json:"description,omitempty" firestore:"description"`
	Amount              int32     `json:"amount" firestore:"amount"`
	Category            string    `json:"category,omitempty" firestore:"category"`
	Type                string    `json:"type" firestore:"type"`
	BankReference       string    `json:"bankReference,omitempty" firestore:"bankReference,omitempty"`
	InsertedAt          time.Time `json:"insertedAt" firestore:"insertedAt"`
	UpdatedAt           time.Time `json:"updatedAt" firestore:"updatedAt"`
}
