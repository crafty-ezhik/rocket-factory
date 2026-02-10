package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, login, password string) (uuid.UUID, error)
	Whoami(ctx context.Context, sessionUUID uuid.UUID) (model.WhoamiResponse, error)
}

type UserService interface {
	Register(ctx context.Context, userInfo model.UserRegistrationInfo) (uuid.UUID, error)
	Get(ctx context.Context, userUUID uuid.UUID) (model.User, error)
}
