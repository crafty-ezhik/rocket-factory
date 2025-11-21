package auth

import (
	"context"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	"github.com/google/uuid"
)

func (s *service) Login(ctx context.Context, login, password string) (uuid.UUID, error) {
	user, err := s.userRepo.Exist(ctx, login)
	if err != nil {
		return uuid.Nil, err
	}

	err = s.hasher.Verify(user.Info.PasswordHash, password)
	if err != nil {
		return uuid.Nil, model.ErrInvalidCredentials
	}

	sessionUUID, err := s.sessionRepo.Create(ctx, user.UUID)
	if err != nil {
		return uuid.Nil, err
	}

	return sessionUUID, nil
}
