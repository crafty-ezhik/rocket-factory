package order_assembled_consumer

import (
	"context"
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

	// Отправка в ТГ
	err = s.tgService.SendOrderAssembledNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send order assemble notification", zap.Error(err))
		return err
	}
	return nil
}
