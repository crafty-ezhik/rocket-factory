package decoder

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
	eventsV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderAssembledDecoder() *decoder { return &decoder{} }

func (d *decoder) Decode(data []byte) (model.OrderAssembledEvent, error) {
	var pb eventsV1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	var event model.OrderAssembledEvent

	eventUUID, err := uuid.Parse(pb.EventUuid)
	if err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to parse event uuid: %w", err)
	}
	event.EventUUID = eventUUID

	orderUUID, err := uuid.Parse(pb.OrderUuid)
	if err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to parse order uuid: %w", err)
	}
	event.OrderUUID = orderUUID

	userUUID, err := uuid.Parse(pb.UserUuid)
	if err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to parse user uuid: %w", err)
	}
	event.UserUUID = userUUID

	event.BuildTimeSec = int(pb.BuildTimeSec)

	return event, nil
}
