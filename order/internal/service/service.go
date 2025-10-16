package service

import (
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository/model"
	"github.com/google/uuid"
)

type OrderService interface {
	Get(orderID uuid.UUID) (model.Order, error)
	Create(userID uuid.UUID, parts []uuid.UUID) (uuid.UUID, error)
	Cancel(orderID uuid.UUID) error
	Pay(orderID uuid.UUID, paymentMethod string) (uuid.UUID, error)
}
