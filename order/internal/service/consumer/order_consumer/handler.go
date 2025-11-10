package order_consumer

import (
	"context"
	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) OrderAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembled event", zap.Error(err))
		return err
	}

	order, err := s.orderRepo.Get(ctx, event.OrderUUID)
	if err != nil {
		logger.Error(ctx, "Failed to get order", zap.Error(err))
		return err
	}

	order.Status = model.OrderStatusASSEMBLED

	if err = s.orderRepo.Update(ctx, order); err != nil {
		logger.Error(ctx, "Failed to update order", zap.Error(err))
		return err
	}
	return nil
}
