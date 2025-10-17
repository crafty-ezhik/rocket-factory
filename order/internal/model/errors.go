package model

import "errors"

var (
	ErrOrderNotFound  = errors.New("order not found")
	ErrOrderIsPaid    = errors.New("order has already been paid for")
	ErrOrderIsCancel  = errors.New("order has already been cancelled")
	ErrOrderCannotPay = errors.New("order has already been paid or cancelled")
)
