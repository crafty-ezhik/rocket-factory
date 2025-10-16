package part

import (
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/repository"
	def "github.com/crafty-ezhik/rocket-factory/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	inventoryRepo repository.InventoryRepository
}

func NewService(inventoryRepo repository.InventoryRepository) *service {
	return &service{inventoryRepo: inventoryRepo}
}
