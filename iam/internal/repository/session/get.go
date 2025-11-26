package session

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/crafty-ezhik/rocket-factory/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, sessionUUID uuid.UUID) (serviceModel.Session, error) {
	sessionData, err := r.redis.HGetAll(ctx, sessionUUID.String())
	if err != nil {
		return serviceModel.Session{}, fmt.Errorf("failed to get session data: %w", err)
	}

	var session repoModel.Session
	err = redis.ScanStruct(sessionData, &session)
	if err != nil {
		return serviceModel.Session{}, fmt.Errorf("failed to scan hash fields to struct: %w", err)
	}

	return converter.SessionToServiceModel(session), nil
}
