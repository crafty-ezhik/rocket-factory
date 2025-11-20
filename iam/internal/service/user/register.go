package user

import (
	"context"
	"errors"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/hasher"
	"github.com/google/uuid"
	"net/mail"
)

func (s *service) Register(ctx context.Context, userInfo model.UserRegistrationInfo) (uuid.UUID, error) {
	if _, err := mail.ParseAddress(userInfo.Info.Email); err != nil {
		return uuid.Nil, model.ErrInvalidEmail
	}

	hashedPassword, err := s.hasher.Hash(userInfo.Password)
	if err != nil {
		if errors.Is(err, hasher.ErrWeakPassword) {
			return uuid.Nil, model.ErrWeakPassword
		}
		return uuid.Nil, err
	}

	userUUID, err := s.userRepo.Create(ctx, userInfo, hashedPassword)
	if err != nil {
		return uuid.Nil, err
	}
	return userUUID, nil
}
