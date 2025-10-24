package v1

import (
	"context"
	"errors"

	"github.com/crafty-ezhik/rocket-factory/inventory/internal/converter"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (s *ApiSuite) TestListPartsSuccess() {
	parts := []model.Part{
		{Name: "B57D30"},
		{Name: "OM137"},
	}

	tests := []struct {
		name         string
		req          *inventoryV1.ListPartsRequest
		filters      *inventoryV1.PartsFilter
		expectedResp *inventoryV1.ListPartsResponse
		expectedErr  error
		setupMock    func(filters *inventoryV1.PartsFilter)
	}{
		{
			name: "success empty filters",
			req: &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{},
			},
			filters: &inventoryV1.PartsFilter{},
			expectedResp: &inventoryV1.ListPartsResponse{
				Parts: converter.SlicePartToProto(parts),
			},
			expectedErr: nil,
			setupMock: func(filters *inventoryV1.PartsFilter) {
				s.inventoryService.On("List", s.ctx, converter.PartsFilterToServiceModel(filters)).
					Return(parts, nil).
					Once()
			},
		},
		{
			name: "success empty list",
			req: &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Names: []string{"B58"},
				},
			},
			filters: &inventoryV1.PartsFilter{
				Names: []string{"B58"},
			},
			expectedResp: &inventoryV1.ListPartsResponse{
				Parts: converter.SlicePartToProto([]model.Part{}),
			},
			expectedErr: nil,
			setupMock: func(filters *inventoryV1.PartsFilter) {
				s.inventoryService.On("List", s.ctx, converter.PartsFilterToServiceModel(filters)).
					Return([]model.Part{}, nil).
					Once()
			},
		},
		{
			name: "success filtered parts",
			req: &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Names: []string{"B57D30"},
				},
			},
			filters: &inventoryV1.PartsFilter{
				Names: []string{"B57D30"},
			},
			expectedResp: &inventoryV1.ListPartsResponse{
				Parts: converter.SlicePartToProto(parts[:1]),
			},
			expectedErr: nil,
			setupMock: func(filters *inventoryV1.PartsFilter) {
				s.inventoryService.On("List", s.ctx, converter.PartsFilterToServiceModel(filters)).
					Return(parts[:1], nil).
					Once()
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock(tt.filters)

			result, err := s.api.ListParts(s.ctx, tt.req)

			s.Require().NoError(err)
			s.Require().Equal(tt.expectedResp, result)
		})
	}
}

func (s *ApiSuite) TestListPartsFailure() {
	dbErr := errors.New("something went wrong")

	tests := []struct {
		name           string
		req            *inventoryV1.ListPartsRequest
		filters        *inventoryV1.PartsFilter
		expectedResp   *inventoryV1.ListPartsResponse
		expectedErrMsg string
		setupMock      func(filters *inventoryV1.PartsFilter)
	}{
		{
			name: "failure service timeout",
			req: &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Names: []string{"B57D30"},
				},
			},
			filters: &inventoryV1.PartsFilter{
				Names: []string{"B57D30"},
			},
			expectedResp:   nil,
			expectedErrMsg: "context deadline exceeded",
			setupMock: func(filters *inventoryV1.PartsFilter) {
				s.inventoryService.On("List", s.ctx, converter.PartsFilterToServiceModel(filters)).
					Return([]model.Part{}, context.DeadlineExceeded).
					Once()
			},
		},
		{
			name: "failure service cancelled",
			req: &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Names: []string{"B57D30"},
				},
			},
			filters: &inventoryV1.PartsFilter{
				Names: []string{"B57D30"},
			},
			expectedResp:   nil,
			expectedErrMsg: "context canceled",
			setupMock: func(filters *inventoryV1.PartsFilter) {
				s.inventoryService.On("List", s.ctx, converter.PartsFilterToServiceModel(filters)).
					Return([]model.Part{}, context.Canceled).
					Once()
			},
		},
		{
			name: "failure internal server error",
			req: &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Names: []string{"B57D30"},
				},
			},
			filters: &inventoryV1.PartsFilter{
				Names: []string{"B57D30"},
			},
			expectedResp:   nil,
			expectedErrMsg: "something went wrong",
			setupMock: func(filters *inventoryV1.PartsFilter) {
				s.inventoryService.On("List", s.ctx, converter.PartsFilterToServiceModel(filters)).
					Return([]model.Part{}, dbErr).
					Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock(tt.filters)
			res, err := s.api.ListParts(s.ctx, tt.req)

			s.Require().Nil(res)
			s.Require().Error(err)
			s.Require().Contains(err.Error(), tt.expectedErrMsg)
		})
	}
}
