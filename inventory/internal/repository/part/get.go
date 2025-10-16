package part

import (
	"context"

	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, partID uuid.UUID) (serviceModel.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.data[partID.String()]
	if !ok {
		return serviceModel.Part{}, serviceModel.ErrPartNotFound
	}
	return converter.PartToServiceModel(part), nil
}
