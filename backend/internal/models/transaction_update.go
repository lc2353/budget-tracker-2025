package models

import "time"

type TransactionUpdate struct {
	Description *string   `json:"description,omitempty" firestore:"description,omitempty"`
	Amount      *int32    `json:"amount,omitempty" firestore:"amount,omitempty"`
	Category    *string   `json:"category,omitempty" firestore:"category,omitempty"`
	Type        *string   `json:"type,omitempty" firestore:"type,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" firestore:"updatedAt,omitempty"`
}
