package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/crafty-ezhik/rocket-factory/payment/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	paymentService *mocks.MockPaymentService

	api *API
}

func (s *APISuite) SetupSuite() {
	s.ctx = context.Background()

	s.paymentService = mocks.NewMockPaymentService(s.T())

	s.api = NewAPI(s.paymentService)
}

func (s *APISuite) TestPaymentFail() {}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
