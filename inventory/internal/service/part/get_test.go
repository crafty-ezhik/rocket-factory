package part

import (
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	"github.com/google/uuid"
	"time"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	partUUID := uuid.New()
	part := model.Part{
		UUID:          partUUID,
		Name:          gofakeit.Name(),
		Description:   gofakeit.HackerPhrase(),
		Price:         gofakeit.Float64Range(0, 1000),
		StockQuantity: gofakeit.Int64(),
		Category:      gofakeit.Word(),
		Dimensions: &model.Dimensions{
			Length: gofakeit.Float64Range(0, 100),
			Height: gofakeit.Float64Range(0, 100),
			Width:  gofakeit.Float64Range(0, 100),
			Weight: gofakeit.Float64Range(0, 100),
		},
		Manufacturer: &model.Manufacturer{
			Name:    gofakeit.Name(),
			Website: gofakeit.URL(),
			Country: gofakeit.Country(),
		},
		Tags:      []string{gofakeit.Word(), gofakeit.Word(), gofakeit.Word()},
		Metadata:  map[string]any{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.inventoryRepo.On("Get", s.ctx, partUUID).
		Return(part, nil).Once()

	resp, err := s.service.Get(s.ctx, partUUID)

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Equal(partUUID, resp.UUID)
}

func (s *ServiceSuite) TestGetPartFailure() {
	partUUID := uuid.New()
	dbErr := errors.New("db error")

	tests := []struct {
		name        string
		repoErr     error
		expectedErr error
		expectWrap  bool
	}{
		{
			name:        "not_found",
			repoErr:     model.ErrPartNotFound,
			expectedErr: model.ErrPartNotFound,
			expectWrap:  false,
		},
		{
			name:        "db_error",
			repoErr:     fmt.Errorf("wrapped error: %w", dbErr),
			expectedErr: dbErr,
			expectWrap:  true,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.inventoryRepo.On("Get", s.ctx, partUUID).
				Return(model.Part{}, test.repoErr).Once()

			resp, err := s.service.Get(s.ctx, partUUID)

			s.Require().Error(err)
			if test.expectWrap {
				s.Require().ErrorIs(err, test.expectedErr)
				s.Contains(err.Error(), test.expectedErr.Error())
			} else {
				s.Require().EqualError(err, test.expectedErr.Error())
			}

			s.Zero(resp)
		})

	}

}
