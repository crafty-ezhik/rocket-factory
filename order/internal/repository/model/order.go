package model

import (
	"github.com/google/uuid"

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
}

type UpdateOrderInfo struct {
	UUID            uuid.UUID
	TransactionUUID uuid.UUID
	PaymentMethod   model.PaymentMethod
	OrderStatus     model.OrderStatus
}
