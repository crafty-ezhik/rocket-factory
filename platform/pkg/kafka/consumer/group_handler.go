package consumer

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	"go.uber.org/zap"
)

// MessageHandler — обработчик сообщений.
type MessageHandler func(ctx context.Context, msg kafka.Message) error

// Middleware — функция middleware для дополнительной обработки.
type Middleware func(next MessageHandler) MessageHandler

// groupHandler — обёртка для sarama.ConsumerGroupHandler
type groupHandler struct {
	handler MessageHandler
	logger  Logger
}

// NewGroupHandler создаёт новый groupHandler с middleware цепочкой.
func NewGroupHandler(handler MessageHandler, logger Logger, middlewares ...Middleware) *groupHandler {
	// Применяем middleware цепочку
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return &groupHandler{
		logger:  logger,
		handler: handler,
	}
}

func (g *groupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (g *groupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (g *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				g.logger.Info(session.Context(), "Kafka message channel closed")
				return nil
			}

			// Преобразуем сообщение из Kafka к нашей универсальной обертке
			msg := kafka.Message{
				Key:            message.Key,
				Value:          message.Value,
				Topic:          message.Topic,
				Partition:      message.Partition,
				Offset:         message.Offset,
				Timestamp:      message.Timestamp,
				BlockTimestamp: message.BlockTimestamp,
				Headers:        extractHeaders(message.Headers),
			}

			// Передаем полученное сообщение в обработчик
			if err := g.handler(session.Context(), msg); err != nil {
				g.logger.Info(session.Context(), "Kafka message handler failed", zap.Error(err))
				continue
			}

			// Помечаем сообщение прочитанным
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			g.logger.Info(session.Context(), "Kafka session closing")
			return nil
		}
	}
}

func extractHeaders(headers []*sarama.RecordHeader) map[string][]byte {
	result := make(map[string][]byte)
	for _, h := range headers {
		if h != nil && h.Key != nil {
			result[string(h.Key)] = h.Value
		}
	}

	return result
}
