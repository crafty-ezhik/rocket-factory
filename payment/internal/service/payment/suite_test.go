package payment_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/crafty-ezhik/rocket-factory/payment/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	service *mocks.MockPaymentService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	s.service = mocks.NewMockPaymentService(s.T())
}

func (s *ServiceSuite) TearDownSuite() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
