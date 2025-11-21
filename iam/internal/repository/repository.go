package repository

import (
	"context"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	Get(ctx context.Context, userUUID uuid.UUID) (model.User, error)
	Create(ctx context.Context, info model.UserRegistrationInfo, hashedPassword string) (uuid.UUID, error)
	Exist(ctx context.Context, login string) (model.User, error)
}

type SessionRepository interface {
	Get(ctx context.Context, sessionUUID uuid.UUID) (model.Session, error)
	Create(ctx context.Context, userUUID uuid.UUID) (uuid.UUID, error)
	AddToUserSet(ctx context.Context, userUUID uuid.UUID, sessionUUID uuid.UUID) error
}
