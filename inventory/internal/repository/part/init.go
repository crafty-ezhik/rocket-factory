package part

import (
	"math"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	repoModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/model"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

// Init - наполняет хранилище n случайными запчастями
func (r *repository) Init(n int) {
	catSlice := []string{
		inventoryV1.Category_ENGINE.String(),
		inventoryV1.Category_FUEL.String(),
		inventoryV1.Category_PORTHOLE.String(),
		inventoryV1.Category_WING.String(),
		inventoryV1.Category_UNKNOWN_UNSPECIFIED.String(),
	}

	for range n {
		data := repoModel.Part{
			UUID:          uuid.New(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.HackerPhrase(),
			Price:         math.Floor(gofakeit.Float64Range(1, 1000)*100) / 100,
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      catSlice[gofakeit.IntRange(0, len(catSlice)-1)],
			Dimensions: &repoModel.Dimensions{
				Length: gofakeit.Float64Range(1, 10000),
				Width:  gofakeit.Float64Range(1, 10000),
				Height: gofakeit.Float64Range(1, 10000),
				Weight: gofakeit.Float64Range(1, 10000),
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    gofakeit.Name(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags:      []string{gofakeit.Word(), gofakeit.Word(), gofakeit.Word()},
			Metadata:  nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		r.data[data.UUID.String()] = data
	}
}
