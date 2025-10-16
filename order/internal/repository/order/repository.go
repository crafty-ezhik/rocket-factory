package order

import (
	repoModel "github.com/crafty-ezhik/rocket-factory/order/internal/repository/model"
	"github.com/google/uuid"
	"sync"
)

type repository struct {
	mu   sync.RWMutex
	data map[uuid.UUID]repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		data: make(map[uuid.UUID]repoModel.Order),
	}
}
