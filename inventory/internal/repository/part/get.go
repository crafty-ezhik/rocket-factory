package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/model"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

func (r *repository) Get(ctx context.Context, partID uuid.UUID) (serviceModel.Part, error) {
	var part repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{partFieldPartUUID: partID}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return serviceModel.Part{}, serviceModel.ErrPartNotFound
		}
		logger.Error(ctx, "Ошибка получения детали", zap.Error(err))
		return serviceModel.Part{}, fmt.Errorf("error receiving data: %w", err)
	}
	return converter.PartToServiceModel(part), nil
}
