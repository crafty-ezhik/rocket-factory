package order

import (
	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/google/uuid"
)

func (r *repository) Create(order serviceModel.Order) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
