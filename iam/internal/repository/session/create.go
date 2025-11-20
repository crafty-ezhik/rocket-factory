package session

import (
	"context"
	"github.com/google/uuid"
)

func (r *repository) Create(ctx context.Context, userUUID uuid.UUID) (uuid.UUID, error) {
	return uuid.New(), nil
}
