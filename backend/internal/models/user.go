package models

type User struct {
	ID        string `json:"id,omitempty" firestore:"-"`
	Email     string `json:"email" firestore:"email"`
	FirstName string `json:"firstName" firestore:"firstName"`
	LastName  string `json:"lastName" firestore:"lastName"`
	CreatedAt string `json:"createdAt" firestore:"createdAt"`
	UpdatedAt string `json:"updatedAt" firestore:"updatedAt"`
}
