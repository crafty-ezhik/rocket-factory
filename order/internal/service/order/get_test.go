package order

import (
	"errors"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestGetOrder() {
	dbErr := errors.New("DB error")
	orderUUID := uuid.New()

	tests := []struct {
		name        string
		orderID     uuid.UUID
		expectedRes model.Order
		expectedErr error
	}{
		{
			name:        "success",
			orderID:     orderUUID,
			expectedRes: model.Order{UUID: orderUUID},
			expectedErr: nil,
		},
		{
			name:        "order not found",
			orderID:     uuid.Nil,
			expectedRes: model.Order{},
			expectedErr: model.ErrOrderNotFound,
		},
		{
			name:        "db error",
			orderID:     orderUUID,
			expectedRes: model.Order{},
			expectedErr: dbErr,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.repo.On("Get", s.ctx, test.orderID).
				Return(test.expectedRes, test.expectedErr).Once()

			res, err := s.service.Get(s.ctx, test.orderID)

			s.Require().Equal(test.expectedRes, res)
			s.Require().Equal(test.expectedErr, err)
		})
	}
}
