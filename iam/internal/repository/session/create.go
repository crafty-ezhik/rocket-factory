package session

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const ttl = time.Hour * 24

const (
	fieldUserUUID  = "user_uuid"
	fieldCreatedAt = "created_at"
	fieldUpdatedAt = "updated_at"
	fieldExpiresAt = "expires_at"
)

func (r *repository) Create(ctx context.Context, userUUID uuid.UUID) (uuid.UUID, error) {
	sessionUUID := uuid.New()
	createdAt := time.Now().Unix()
	expiresAt := time.Now().Add(ttl).Unix()
	updatedAt := time.Now().Unix()

	fields := map[string]any{
		fieldUserUUID:  userUUID.String(),
		fieldCreatedAt: createdAt,
		fieldUpdatedAt: updatedAt,
		fieldExpiresAt: expiresAt,
	}

	err := r.redis.HashSet(ctx, sessionUUID.String(), fields)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to set hash set: %v", err)
	}

	err = r.AddToUserSet(ctx, userUUID, sessionUUID)
	if err != nil {
		return uuid.Nil, err
	}

	return sessionUUID, nil
}
