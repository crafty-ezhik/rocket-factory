package repository

import (
	"github.com/google/uuid"

	repoModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/model"
)

type InventoryRepository interface {
	Get(partID uuid.UUID) (repoModel.Part, error)
	List() ([]repoModel.Part, error)
	Init(n int)
}
