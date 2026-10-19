package loan

import "errors"

// Sentinel errors. RecordPayment, RemainingOf, and MarkDefaulted must return these.
var (
	ErrUnauthorized            = errors.New("unauthorized")
	ErrInvalidAddress          = errors.New("invalid address")
	ErrInvalidAmount           = errors.New("invalid amount")
	ErrLoanExists              = errors.New("loan exists")
	ErrLoanNotFound            = errors.New("loan not found")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)
