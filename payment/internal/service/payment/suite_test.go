package payment_test

import (
	"context"
	"github.com/crafty-ezhik/rocket-factory/payment/internal/service/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

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
