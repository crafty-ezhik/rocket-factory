package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *repository) AddToUserSet(ctx context.Context, userUUID, sessionUUID uuid.UUID) error {
	err := r.redis.SAdd(ctx, userUUID.String(), sessionUUID.String())
	if err != nil {
		return fmt.Errorf("failed to add to user set: %w", err)
	}
	return nil
}
