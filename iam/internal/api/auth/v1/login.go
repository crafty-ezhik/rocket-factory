package v1

import (
	"context"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	authV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	if req == nil {
		return &authV1.LoginResponse{}, model.ErrInvalidCredentials
	}

	if req.Login == "" || req.Password == "" {
		return &authV1.LoginResponse{}, model.ErrInvalidCredentials
	}

	sessionUUID, err := a.service.Login(ctx, req.Login, req.Password)
	if err != nil {
		return &authV1.LoginResponse{}, err
	}

	return &authV1.LoginResponse{
		SessionUuid: sessionUUID.String(),
	}, nil
}
