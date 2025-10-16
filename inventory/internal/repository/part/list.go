package part

import (
	"context"
	repoModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) List(ctx context.Context, filters repoModel.PartsFilter) ([]repoModel.Part, error) {
	// TODO: Реализовать фильтрацию
	partList := make([]repoModel.Part, 0)

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.data {
		partList = append(partList, item)
	}
	return partList, nil
}
