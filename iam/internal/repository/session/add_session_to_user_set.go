package session

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

func (r *repository) AddToUserSet(ctx context.Context, userUUID uuid.UUID, sessionUUID uuid.UUID) error {

	err := r.redis.SAdd(ctx, userUUID.String(), sessionUUID.String())
	if err != nil {
		return fmt.Errorf("failed to add to user set: %v", err)
	}
	return errors.New("not implemented")
}
