package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.orderRepo.Get(ctx, orderID)
	if err != nil {
		return err
	}

	switch order.Status {
	case model.OrderStatusCANCELLED:
		return model.ErrOrderIsCancel
	case model.OrderStatusPAID:
		return model.ErrOrderIsPaid
	}

	order.Status = model.OrderStatusCANCELLED

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return err
	}

	return nil
}
