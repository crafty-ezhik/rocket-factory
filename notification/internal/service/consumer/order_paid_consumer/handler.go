package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembled event", zap.Error(err))
		return err
	}

	// Отправка в ТГ
	err = s.tgService.SendOrderPaidNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send order paid notification", zap.Error(err))
		return err
	}
	return nil
}
