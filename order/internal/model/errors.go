package model

import "errors"

var (
	ErrNotFound    = errors.New("order not found")
	ErrOrderIsPaid = errors.New("order is paid")
)
