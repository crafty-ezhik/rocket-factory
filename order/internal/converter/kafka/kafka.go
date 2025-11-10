package kafka

import "github.com/crafty-ezhik/rocket-factory/order/internal/model"

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.OrderAssembledEvent, error)
}
