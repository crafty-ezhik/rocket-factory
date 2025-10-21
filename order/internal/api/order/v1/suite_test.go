package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/crafty-ezhik/rocket-factory/order/internal/service/mocks"
)

type ApiSuite struct {
	suite.Suite
	ctx          context.Context
	orderService *mocks.MockOrderService
	api          *api
}

func (s *ApiSuite) SetupSuite() {
	s.ctx = context.Background()
	s.orderService = mocks.NewMockOrderService(s.T())
	s.api = NewAPI(s.orderService)
}
func (s *ApiSuite) TearDownSuite() {}

func TestApiIntegration(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
