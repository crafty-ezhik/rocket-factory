package session

import (
	"context"
	"fmt"
	repoModel "github.com/crafty-ezhik/rocket-factory/iam/internal/repository/model"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

func (r *repository) Get(ctx context.Context, sessionUUID uuid.UUID) (repoModel.Session, error) {
	sessionData, err := r.redis.HGetAll(ctx, sessionUUID.String())
	if err != nil {
		return repoModel.Session{}, fmt.Errorf("failed to get session data: %w", err)
	}

	var session repoModel.Session
	err = redis.ScanStruct(sessionData, &session)
	if err != nil {
		return repoModel.Session{}, fmt.Errorf("failed to scan hash fields to struct: %v\n", err)
	}

	return session, nil
}
