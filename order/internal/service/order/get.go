package order

import (
	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) Get(orderID uuid.UUID) (model.Order, error) {
	return model.Order{}, nil
}
