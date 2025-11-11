package service

import (
	"context"
	"github.com/crafty-ezhik/rocket-factory/notification/internal/model"
)

type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, msg model.OrderPaidEvent) error
	SendOrderAssembledNotification(ctx context.Context, msg model.OrderAssembledEvent) error
}

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}
