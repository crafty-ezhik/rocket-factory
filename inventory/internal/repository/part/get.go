package part

import (
	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	repoModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) Get(partID uuid.UUID) (repoModel.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.data[partID.String()]
	if !ok {
		return repoModel.Part{}, serviceModel.ErrPartNotFound
	}
	return part, nil
}
