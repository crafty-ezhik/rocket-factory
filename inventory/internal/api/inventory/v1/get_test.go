package v1

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/inventory/internal/converter"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (s *ApiSuite) TestGetOrderSuccess() {
	partUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	basePart := model.Part{
		UUID:          partUUID,
		Name:          "Engine",
		Description:   "Ratata",
		Price:         999.99,
		StockQuantity: 10,
		Category:      "engine",
		Dimensions:    &model.Dimensions{},
		Manufacturer:  &model.Manufacturer{},
		Tags:          []string{"V8", "Diesel", "high torque"},
		Metadata:      nil,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tests := []struct {
		name         string
		req          *inventoryV1.GetPartRequest
		expectedResp *inventoryV1.GetPartResponse
		expectedErr  error
		setupMock    func()
	}{
		{
			name: "success",
			req: &inventoryV1.GetPartRequest{
				Uuid: partUUID.String(),
			},
			expectedResp: &inventoryV1.GetPartResponse{
				Part: converter.PartToProto(basePart),
			},
			expectedErr: nil,
			setupMock: func() {
				s.inventoryService.On("Get", s.ctx, partUUID).
					Return(basePart, nil).
					Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			result, err := s.api.GetPart(s.ctx, tt.req)

			s.Require().NoError(err)
			s.Require().Equal(tt.expectedResp, result)
		})
	}
}

func (s *ApiSuite) TestGetOrderFailure() {
	partUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	invalidPartUUID := "00000000-0000-0000-0000-0000000000033333"
	dbErr := errors.New("something went wrong")

	tests := []struct {
		name           string
		req            *inventoryV1.GetPartRequest
		expectedErrMsg string
		setupMock      func()
	}{
		{
			name: "invalid part uuid",
			req: &inventoryV1.GetPartRequest{
				Uuid: invalidPartUUID,
			},
			expectedErrMsg: "invalid UUID",
			setupMock:      func() {},
		},
		{
			name: "part not found",
			req: &inventoryV1.GetPartRequest{
				Uuid: partUUID.String(),
			},
			expectedErrMsg: "part not found",
			setupMock: func() {
				s.inventoryService.On("Get", s.ctx, partUUID).
					Return(model.Part{}, model.ErrPartNotFound).
					Once()
			},
		},
		{
			name: "service timeout",
			req: &inventoryV1.GetPartRequest{
				Uuid: partUUID.String(),
			},
			expectedErrMsg: "context deadline exceeded",
			setupMock: func() {
				s.inventoryService.On("Get", s.ctx, partUUID).
					Return(model.Part{}, context.DeadlineExceeded).
					Once()
			},
		},
		{
			name: "service canceled",
			req: &inventoryV1.GetPartRequest{
				Uuid: partUUID.String(),
			},
			expectedErrMsg: "canceled",
			setupMock: func() {
				s.inventoryService.On("Get", s.ctx, partUUID).
					Return(model.Part{}, context.Canceled).
					Once()
			},
		},
		{
			name: "service internal error",
			req: &inventoryV1.GetPartRequest{
				Uuid: partUUID.String(),
			},
			expectedErrMsg: "something went wrong",
			setupMock: func() {
				s.inventoryService.On("Get", s.ctx, partUUID).
					Return(model.Part{}, dbErr).
					Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			res, err := s.api.GetPart(s.ctx, tt.req)

			s.Require().Nil(res)
			s.Require().Error(err)
			s.Require().Contains(err.Error(), tt.expectedErrMsg)
		})
	}
}
