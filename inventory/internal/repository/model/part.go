package model

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Part struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UUID          uuid.UUID          `bson:"part_uuid"`
	Name          string             `bson:"name"`
	Description   string             `bson:"description"`
	Price         float64            `bson:"price"`
	StockQuantity int64              `bson:"stock_quantity"`
	Category      string             `bson:"category"`
	Dimensions    *Dimensions        `bson:"dimensions"`
	Manufacturer  *Manufacturer      `bson:"manufacturer"`
	Tags          []string           `bson:"tags"`
	Metadata      map[string]any     `bson:"metadata"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}
