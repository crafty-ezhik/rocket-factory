package part

import (
	"errors"

	"github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestListPartsSuccess() {
	storage := []model.Part{
		{Name: "first"},
		{Name: "second"},
		{Name: "third"},
	}

	tests := []struct {
		name           string
		filters        model.PartsFilter
		expectedResult []model.Part
		repoErr        error
		expectedError  error
	}{
		{
			name:           "success_empty_filters",
			filters:        model.PartsFilter{},
			expectedResult: storage,
			repoErr:        nil,
			expectedError:  nil,
		},
		{
			name:    "success_apply_filters",
			filters: model.PartsFilter{Names: []string{"first"}},
			expectedResult: []model.Part{
				{Name: "first"},
			},
			repoErr:       nil,
			expectedError: nil,
		},
		{
			name:           "success_empty_list",
			filters:        model.PartsFilter{Names: []string{"fourth"}},
			expectedResult: []model.Part{},
			repoErr:        nil,
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.inventoryRepo.On("List", s.ctx, test.filters).
				Return(test.expectedResult, test.repoErr).Once()

			res, err := s.inventoryRepo.List(s.ctx, test.filters)

			s.Require().NoError(err)
			s.Require().Equal(test.expectedResult, res)
		})
	}
}

func (s *ServiceSuite) TestListPartsFailure() {
	dbErr := errors.New("db error")

	tests := []struct {
		name          string
		filters       model.PartsFilter
		repoErr       error
		expectedError error
	}{
		{
			name:          "failure_repo_err",
			filters:       model.PartsFilter{Names: []string{"first"}},
			repoErr:       dbErr,
			expectedError: dbErr,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.inventoryRepo.On("List", s.ctx, test.filters).
				Return([]model.Part{}, test.repoErr).Once()

			res, err := s.inventoryRepo.List(s.ctx, test.filters)

			s.Require().Error(err)
			s.Empty(res)
		})
	}
}
