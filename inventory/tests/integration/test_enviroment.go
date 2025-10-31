package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/model"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

// InsertTestPart - вставляет тестовую деталь в коллекцию Mongo и возвращает ее UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := uuid.New()
	now := time.Now()

	partDoc := model.Part{
		ID:            primitive.NewObjectID(),
		UUID:          partUUID,
		Name:          gofakeit.Name(),
		Description:   gofakeit.Phrase(),
		Price:         gofakeit.Float64Range(0, 1000),
		StockQuantity: int64(gofakeit.IntN(100)),
		Category:      "FUEL",
		Dimensions: &model.Dimensions{
			Length: gofakeit.Float64Range(0, 1000),
			Width:  gofakeit.Float64Range(0, 1000),
			Height: gofakeit.Float64Range(0, 1000),
			Weight: gofakeit.Float64Range(0, 1000),
		},
		Manufacturer: &model.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags:      []string{gofakeit.Word(), gofakeit.Word(), gofakeit.Word(), gofakeit.Word()},
		Metadata:  nil,
		CreatedAt: now,
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "parts" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID.String(), nil
}

func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, part *model.Part) (string, error) {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "parts" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, part)
	if err != nil {
		return "", err
	}

	return part.UUID.String(), nil
}

// GetTestPart - возвращает информацию о детали
func (env *TestEnvironment) GetTestPart() *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          uuid.NewString(),
		Name:          "Двигатель М123 Y7",
		Description:   "Стартовый двигатель нового поколения",
		Price:         12024,
		StockQuantity: 2,
		Category:      2,
		Dimensions: &inventoryV1.Dimensions{
			Length: gofakeit.Float64Range(0, 1000),
			Width:  gofakeit.Float64Range(0, 1000),
			Height: gofakeit.Float64Range(0, 1000),
			Weight: gofakeit.Float64Range(0, 1000),
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags:      []string{gofakeit.Word(), gofakeit.Word()},
		Metadata:  nil,
		CreatedAt: timestamppb.New(time.Now().Add(-2 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now()),
	}
}

// GetTestParts - возвращает информацию о нескольких деталях
func (env *TestEnvironment) GetTestParts() []*inventoryV1.Part {
	return []*inventoryV1.Part{
		env.GetTestPart(),
		env.GetTestPart(),
	}
}

// ClearPartsCollection - удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "parts" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
