package service

import "github.com/google/uuid"

type OrderService interface {
	Get(orderID uuid.UUID) (string, error)
	Create(userID uuid.UUID, parts []uuid.UUID) (string, error)
	Cancel(orderID uuid.UUID) error
	Pay(orderID uuid.UUID) (string, error)
}
