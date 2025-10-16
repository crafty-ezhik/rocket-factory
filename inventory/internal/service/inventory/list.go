package inventory

import (
	"context"

	serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filters serviceModel.PartsFilter) ([]serviceModel.Part, error) {
	parts, err := s.inventoryRepo.List(ctx, filters)
	if err != nil {
		return []serviceModel.Part{}, err
	}
	return parts, nil
}
