package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderID uuid.UUID) (model.Order, error) {
	return s.orderRepo.Get(ctx, orderID)
}
