package part

import repoModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/model"

func (r *repository) List() ([]repoModel.Part, error) {
	partList := make([]repoModel.Part, 0)

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.data {
		partList = append(partList, item)
	}
	return partList, nil
}
