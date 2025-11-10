package order_producer

import (
	"context"
	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	eventsV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{orderPaidProducer: orderPaidProducer}
}

func (p *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &eventsV1.OrderPaid{
		EventUuid:       event.EventUUID.String(),
		OrderUuid:       event.OrderUUID.String(),
		UserUuid:        event.UserUUID.String(),
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID.String(),
	}

	// Преобразуем структуру в слайс байт для передачи в Kafka
	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Failed to marshal order paid payload", zap.Error(err))
		return err
	}

	// Отправляем сообщение в топик
	err = p.orderPaidProducer.Send(ctx, []byte(event.OrderUUID.String()), payload)
	if err != nil {
		logger.Error(ctx, "Failed to publish OrderPaid", zap.Error(err))
		return err
	}

	return nil
}
