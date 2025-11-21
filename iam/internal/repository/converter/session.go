package converter

import (
	"time"

	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	repoModel "github.com/crafty-ezhik/rocket-factory/iam/internal/repository/model"
)

func SessionToServiceModel(data repoModel.Session) serviceModel.Session {
	var updatedAt *time.Time
	if data.UpdatedAt != nil {
		tmp := time.Unix(*data.UpdatedAt, 0)
		updatedAt = &tmp
	}

	return serviceModel.Session{
		UserUUID:  uuid.MustParse(data.UserUUID),
		CreatedAt: time.Unix(data.CreatedAt, 0),
		UpdatedAt: updatedAt,
		ExpiresAt: time.Unix(data.ExpiresAt, 0),
	}
}
