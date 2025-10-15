package service

import (
	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
)

type InventoryService interface {
	Get(partID uuid.UUID) (serviceModel.Part, error)
	List(filters serviceModel.PartsFilter) ([]serviceModel.Part, error)
}
