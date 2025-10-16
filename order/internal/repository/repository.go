package repository

import (
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create() error
	Get(id uuid.UUID) error
	Update(id uuid.UUID) error
}
