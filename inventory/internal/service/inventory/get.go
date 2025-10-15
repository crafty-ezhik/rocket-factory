package inventory

import (
	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
)

func (s *service) Get(partID uuid.UUID) (serviceModel.Part, error) {
	return serviceModel.Part{}, nil
}
