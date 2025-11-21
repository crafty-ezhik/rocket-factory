package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConv "github.com/crafty-ezhik/rocket-factory/assembly/internal/converter/kafka"
	def "github.com/crafty-ezhik/rocket-factory/assembly/internal/service"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer     kafka.Consumer
	orderPaidDecoder      kafkaConv.OrderPaidDecoder
	orderAssembledService def.OrderProducerService
}

func NewService(
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConv.OrderPaidDecoder,
	orderAssembledService def.OrderProducerService,
) *service {
	return &service{
		orderPaidConsumer:     orderPaidConsumer,
		orderPaidDecoder:      orderPaidDecoder,
		orderAssembledService: orderAssembledService,
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
