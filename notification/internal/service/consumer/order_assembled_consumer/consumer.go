package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConv "github.com/crafty-ezhik/rocket-factory/notification/internal/converter/kafka"
	def "github.com/crafty-ezhik/rocket-factory/notification/internal/service"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConv.OrderAssembledDecoder
	tgService              def.TelegramService
}

func NewService(orderAssembledConsumer kafka.Consumer, orderAssembledDecoder kafkaConv.OrderAssembledDecoder, tgService def.TelegramService) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		tgService:              tgService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderAssembledConsumer service")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}
	return nil
}
