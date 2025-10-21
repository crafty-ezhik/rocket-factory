package model

import (
	"github.com/google/uuid"
	"time"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

type Order struct {
	UUID            uuid.UUID
	UserUUID        uuid.UUID
	PartUUIDs       []uuid.UUID
	TotalPrice      float64
	TransactionUUID uuid.UUID
	PaymentMethod   model.PaymentMethod
	Status          model.OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

type UpdateOrderInfo struct {
	UUID            uuid.UUID
	TransactionUUID uuid.UUID
	PaymentMethod   model.PaymentMethod
	OrderStatus     model.OrderStatus
}
