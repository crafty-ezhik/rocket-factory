package order

import (
	"sync"

	"github.com/google/uuid"

	def "github.com/crafty-ezhik/rocket-factory/order/internal/repository"
	repoModel "github.com/crafty-ezhik/rocket-factory/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[uuid.UUID]repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		data: make(map[uuid.UUID]repoModel.Order),
	}
}
