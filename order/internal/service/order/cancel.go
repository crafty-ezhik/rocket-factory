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

	updatedInfo := model.UpdateOrderInfo{
		UUID:        orderID,
		OrderStatus: model.OrderStatusCANCELLED,
	}

	err = s.orderRepo.Update(ctx, updatedInfo, model.OrderUpdateCANCEL)
	if err != nil {
		return err
	}

	return nil
}
