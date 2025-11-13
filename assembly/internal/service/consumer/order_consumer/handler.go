package order_consumer

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/crafty-ezhik/rocket-factory/assembly/internal/model"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Получен запрос на сборку заказа",
		zap.String("order_uuid", event.OrderUUID.String()),
		zap.String("user_uuid", event.UserUUID.String()),
	)
	// Задержка перед отправкой обратного сообщения
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Second * 10):
		logger.Info(ctx, "✅ Заказ успешно собран!")
	}

	// Отправляем сообщение о завершении сборки
	err = s.orderAssembledService.ProduceOrderAssembled(ctx, model.OrderAssembledEvent{
		EventUUID:    event.EventUUID,
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: 10,
	})
	if err != nil {
		logger.Error(ctx, "Failed to produce order assembled event", zap.Error(err))
		return err
	}

	return nil
}
