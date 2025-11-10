package kafka

import (
	"context"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, handler consumer.MessageHandler) error
}

type Producer interface {
	Send(ctx context.Context, key, value []byte) error
}
