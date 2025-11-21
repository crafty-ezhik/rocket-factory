package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/crafty-ezhik/rocket-factory/assembly/internal/model"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	eventsV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/events/v1"
)

type service struct {
	orderAssembledProducer kafka.Producer
}

func NewService(orderAssembledProducer kafka.Producer) *service {
	return &service{orderAssembledProducer: orderAssembledProducer}
}

func (p *service) ProduceOrderAssembled(ctx context.Context, event model.OrderAssembledEvent) error {
	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUUID.String(),
		OrderUuid:    event.OrderUUID.String(),
		UserUuid:     event.UserUUID.String(),
		BuildTimeSec: int64(event.BuildTimeSec),
	}

	// Преобразуем структуру в слайс байт для передачи в kafka
	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Failed to marshal order assembled payload", zap.Error(err))
		return err
	}

	// Отправляем сообщение в топик
	err = p.orderAssembledProducer.Send(ctx, []byte(event.OrderUUID.String()), payload)
	if err != nil {
		logger.Error(ctx, "Failed to publish OrderAssembled", zap.Error(err))
		return err
	}

	return nil
}
