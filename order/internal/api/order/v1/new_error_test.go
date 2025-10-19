package v1

import (
	"errors"
	"net/http"

	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
)

func (s *ApiSuite) TestNewError() {
	somethingErr := errors.New("something bad happened")

	tests := []struct {
		name        string
		err         error
		expectedRes *orderV1.GenericErrorStatusCode
	}{
		{
			name: "Undefined error",
			err:  somethingErr,
			expectedRes: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusInternalServerError,
				Response: orderV1.GenericError{
					Code:    http.StatusInternalServerError,
					Message: "something bad happened",
				},
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.api.NewError(s.ctx, tt.err)
			s.NotNil(err)
			s.Contains(err.Response.Message, tt.err.Error())
		})
	}
}
