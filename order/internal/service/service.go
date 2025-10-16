package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

type OrderService interface {
	Get(ctx context.Context, orderID uuid.UUID) (model.Order, error)
	Create(ctx context.Context, userID uuid.UUID, parts []uuid.UUID) (uuid.UUID, error)
	Cancel(ctx context.Context, orderID uuid.UUID) error
	Pay(ctx context.Context, orderID, userID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}
