package order_consumer

import (
	"context"
	kafkaConv "github.com/crafty-ezhik/rocket-factory/order/internal/converter/kafka"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository"
	def "github.com/crafty-ezhik/rocket-factory/order/internal/service"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConv.OrderAssembledDecoder
	orderRepo              repository.OrderRepository
}

func NewService(orderAssembledConsumer kafka.Consumer, orderRepo repository.OrderRepository, decoder kafkaConv.OrderAssembledDecoder) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  decoder,
		orderRepo:              orderRepo,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderAssembledConsumer service")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from ufo.recorded topic error", zap.Error(err))
		return err
	}

	return nil
}
