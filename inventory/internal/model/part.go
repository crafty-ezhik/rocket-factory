package model

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      string
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]any
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type PartsFilter struct {
	UUIDs               []string
	Names               []string
	Categories          []string
	ManufacturerCountry []string
	Tags                []string
}

func (pf *PartsFilter) IsEmpty() bool {
	return pf.UUIDs == nil &&
		pf.Names == nil &&
		pf.Categories == nil &&
		pf.ManufacturerCountry == nil &&
		pf.Tags == nil
}
