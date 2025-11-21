package decoder

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/crafty-ezhik/rocket-factory/notification/internal/model"
	eventsV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/events/v1"
)

type decoderPaid struct{}

func NewOrderPaidDecoder() *decoderPaid {
	return &decoderPaid{}
}

func (d *decoderPaid) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	var event model.OrderPaidEvent

	eventUUID, err := uuid.Parse(pb.EventUuid)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to parse event uuid: %w", err)
	}
	event.EventUUID = eventUUID

	orderUUID, err := uuid.Parse(pb.OrderUuid)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to parse order uuid: %w", err)
	}
	event.OrderUUID = orderUUID

	userUUID, err := uuid.Parse(pb.UserUuid)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to parse user uuid: %w", err)
	}
	event.UserUUID = userUUID

	transactionUUID, err := uuid.Parse(pb.TransactionUuid)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to parse transaction uuid: %w", err)
	}
	event.TransactionUUID = transactionUUID

	event.PaymentMethod = pb.PaymentMethod

	return event, nil
}
