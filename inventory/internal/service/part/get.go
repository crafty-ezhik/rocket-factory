package part

import (
	"context"

	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, partID uuid.UUID) (serviceModel.Part, error) {
	return s.inventoryRepo.Get(ctx, partID)
}
