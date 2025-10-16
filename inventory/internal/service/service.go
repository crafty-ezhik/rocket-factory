package service

import (
	"context"
	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
)

type InventoryService interface {
	Get(ctx context.Context, partID uuid.UUID) (serviceModel.Part, error)
	List(ctx context.Context, filters serviceModel.PartsFilter) ([]serviceModel.Part, error)
}
