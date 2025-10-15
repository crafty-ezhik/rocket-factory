package part

import (
	"sync"

	def "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository"
	repoModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]repoModel.Part),
	}
}
