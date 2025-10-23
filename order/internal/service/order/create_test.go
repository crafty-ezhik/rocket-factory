package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrderSuccess() {
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	orderID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	partIDs := []uuid.UUID{
		uuid.MustParse("d195a37b-f2cb-48e6-b739-29db0ddcc197"),
		uuid.MustParse("a79178c5-a082-4884-b214-ee69e3972840"),
	}

	tests := []struct {
		name               string
		userID             uuid.UUID
		partIDs            []uuid.UUID
		expectedOrderID    uuid.UUID
		expectedTotalPrice float64
		expectedErr        error
		setupMocks         func()
	}{
		{
			name:               "success",
			userID:             userID,
			partIDs:            partIDs,
			expectedOrderID:    orderID,
			expectedTotalPrice: 300,
			expectedErr:        nil,
			setupMocks: func() {
				s.inventoryClient.On("ListParts", mock.Anything, model.PartsFilter{
					UUIDs: []string{partIDs[0].String(), partIDs[1].String()},
				}).
					Return([]model.Part{
						{UUID: partIDs[0], Price: 100},
						{UUID: partIDs[1], Price: 200},
					},
						nil).Once()

				s.repo.On("Create", s.ctx, mock.Anything).Return(orderID, nil)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMocks()

			orderUUID, totalPrice, err := s.service.Create(s.ctx, tt.userID, tt.partIDs)

			s.Require().NoError(err)
			s.Require().Equal(tt.expectedOrderID, orderUUID)
			s.Require().Equal(tt.expectedTotalPrice, totalPrice)
		})
	}
}

func (s *ServiceSuite) TestCreateOrder() {
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	partIDs := []uuid.UUID{
		uuid.MustParse("d195a37b-f2cb-48e6-b739-29db0ddcc197"),
		uuid.MustParse("a79178c5-a082-4884-b214-ee69e3972840"),
	}

	clientErr := errors.New("client error")
	dbErr := errors.New("something went wrong")

	tests := []struct {
		name               string
		userID             uuid.UUID
		partIDs            []uuid.UUID
		expectedOrderID    uuid.UUID
		expectedTotalPrice float64
		expectedErr        error
		setupMocks         func()
	}{
		{
			name:               "client error",
			userID:             userID,
			partIDs:            partIDs,
			expectedOrderID:    uuid.Nil,
			expectedTotalPrice: 0,
			expectedErr:        context.DeadlineExceeded,
			setupMocks: func() {
				s.inventoryClient.On("ListParts", mock.Anything, model.PartsFilter{
					UUIDs: []string{partIDs[0].String(), partIDs[1].String()},
				}).
					Return([]model.Part{}, clientErr).
					Once()
			},
		},
		{
			name:               "db error",
			userID:             userID,
			partIDs:            partIDs,
			expectedOrderID:    uuid.Nil,
			expectedTotalPrice: 0,
			expectedErr:        dbErr,
			setupMocks: func() {
				s.inventoryClient.On("ListParts", mock.Anything, model.PartsFilter{
					UUIDs: []string{partIDs[0].String(), partIDs[1].String()},
				}).
					Return([]model.Part{
						{UUID: partIDs[0], Price: 100},
						{UUID: partIDs[1], Price: 200},
					},
						nil).
					Once()

				s.repo.On("Create", s.ctx, mock.Anything).
					Return(uuid.Nil, dbErr).
					Once()
			},
		},
		{
			name:               "part not found",
			userID:             userID,
			partIDs:            partIDs,
			expectedOrderID:    uuid.Nil,
			expectedTotalPrice: 0,
			expectedErr:        fmt.Errorf("part with uuid a79178c5-a082-4884-b214-ee69e3972840 not found"),
			setupMocks: func() {
				s.inventoryClient.On("ListParts", mock.Anything, model.PartsFilter{
					UUIDs: []string{partIDs[0].String(), partIDs[1].String()},
				}).
					Return([]model.Part{
						{UUID: partIDs[0], Price: 100},
					}, nil).
					Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMocks()

			orderUUID, totalPrice, err := s.service.Create(s.ctx, tt.userID, tt.partIDs)
			s.Require().Error(err)
			s.Require().Contains(tt.expectedErr.Error(), err.Error())

			s.Equal(tt.expectedOrderID, orderUUID)
			s.Equal(tt.expectedTotalPrice, totalPrice)
		})
	}
}
