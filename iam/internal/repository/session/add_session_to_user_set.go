package session

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

func (r *repository) AddToUserSet(ctx context.Context, userUUID uuid.UUID, sessionUUID uuid.UUID) error {
	return errors.New("not implemented")
}
