package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderUpdateKind string

const (
	OrderUpdateCANCEL     OrderUpdateKind = "cancel"
	OrderUpdateUPDATEINFO OrderUpdateKind = "update_info"
)

type Order struct {
	UUID            uuid.UUID
	UserUUID        uuid.UUID
	PartUUIDs       []uuid.UUID
	TotalPrice      float64
	TransactionUUID uuid.UUID
	PaymentMethod   PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

type UpdateOrderInfo struct {
	UUID            uuid.UUID
	TransactionUUID uuid.UUID
	PaymentMethod   PaymentMethod
	OrderStatus     OrderStatus
}
