package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConv "github.com/crafty-ezhik/rocket-factory/notification/internal/converter/kafka"
	def "github.com/crafty-ezhik/rocket-factory/notification/internal/service"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConv.OrderPaidDecoder
	tgService         def.TelegramService
}

func NewService(orderPaidDecoder kafkaConv.OrderPaidDecoder, orderPaidConsumer kafka.Consumer, tgService def.TelegramService) *service {
	return &service{
		orderPaidDecoder:  orderPaidDecoder,
		orderPaidConsumer: orderPaidConsumer,
		tgService:         tgService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderPaidConsumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}
	return nil
}
