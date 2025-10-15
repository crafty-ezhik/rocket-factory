package inventory

import serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"

func (s *service) List(filters serviceModel.PartsFilter) ([]serviceModel.Part, error) {
	return []serviceModel.Part{}, nil
}
