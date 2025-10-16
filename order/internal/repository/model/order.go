package model

import "github.com/google/uuid"

type Order struct {
	UUID            uuid.UUID
	UserUUID        uuid.UUID
	PartUUIDs       []uuid.UUID
	TotalPrice      float64
	TransactionUUID uuid.UUID
	PaymentMethod   PaymentMethod
	Status          OrderStatus
}

type UpdateOrderInfo struct {
	UUID            uuid.UUID
	TransactionUUID uuid.UUID
	PaymentMethod   PaymentMethod
	OrderStatus     OrderStatus
}
