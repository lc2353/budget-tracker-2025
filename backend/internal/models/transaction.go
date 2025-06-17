package models

import (
	"time"
)

type Transaction struct {
	ID                  string    `json:"id,omitempty" firestore:"-"`
	UserID              string    `json:"userId" firestore:"userId"`
	TransactionDateTime time.Time `json:"transactionDateTime" firestore:"transactionDateTime"`
	Description         string    `json:"description" firestore:"description"`
	Amount              int32     `json:"amount" firestore:"amount"` // will be stored as pence
	Category            string    `json:"category" firestore:"category"`
	Type                string    `json:"type" firestore:"type"` // "income", "expense", "transfer"
	BankReference       string    `json:"bankReference,omitempty" firestore:"bankReference,omitempty"`
	InsertedAt          time.Time `json:"insertedAt" firestore:"insertedAt"`
	UpdatedAt           time.Time `json:"updatedAt" firestore:"updatedAt"`
}
