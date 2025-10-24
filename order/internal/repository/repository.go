package repository

import (
	"context"

	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order serviceModel.Order) (uuid.UUID, error)
	Get(ctx context.Context, orderID uuid.UUID) (serviceModel.Order, error)
	Update(ctx context.Context, order serviceModel.Order) error
}
