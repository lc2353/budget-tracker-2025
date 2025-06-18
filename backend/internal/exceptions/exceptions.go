package exceptions

import (
	"errors"
	"fmt"
)

// TransactionNotFoundError is returned when a transaction is not found.
type TransactionNotFoundError struct {
	TransactionID string
}

func (e *TransactionNotFoundError) Error() string {
	return fmt.Sprintf("transaction '%s' not found", e.TransactionID)
}

func TransactionNotFound(transactionID string) error {
	return &TransactionNotFoundError{TransactionID: transactionID}
}

// ErrInvalidInput is returned when the input provided is invalid.
var ErrInvalidInput = errors.New("invalid input provided")

// ErrUnauthorized is returned when the user is not authorized to perform an action.
var ErrUnauthorized = errors.New("unauthorized access")

// UserForbiddenError  is returned when the user is forbidden from accessing a resource.
type UserForbiddenError struct {
	UserID string
}

func (e *UserForbiddenError) Error() string {
	return fmt.Sprintf("user '%s' is forbidden from accessing this resource", e.UserID)
}

func UserForbidden(userID string) error {
	return &UserForbiddenError{UserID: userID}
}
