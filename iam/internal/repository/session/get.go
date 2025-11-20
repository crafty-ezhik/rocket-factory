package session

import (
	"context"
	"errors"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	"github.com/google/uuid"
)

func (r *repository) Get(ctx context.Context, sessionUUID uuid.UUID) (model.Session, error) {
	return model.Session{}, errors.New("not implemented")
}
