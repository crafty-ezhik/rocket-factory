package repository

import (
	"context"

	"github.com/google/uuid"

	repoModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/model"
)

type InventoryRepository interface {
	Get(ctx context.Context, partID uuid.UUID) (repoModel.Part, error)
	List(ctx context.Context, filters repoModel.PartsFilter) ([]repoModel.Part, error)
	Init(n int)
}
