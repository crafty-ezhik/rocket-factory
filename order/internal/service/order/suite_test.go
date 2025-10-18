package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMock "github.com/crafty-ezhik/rocket-factory/order/internal/client/grpc/mocks"
	repoMock "github.com/crafty-ezhik/rocket-factory/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx             context.Context //nolint:containedctx
	repo            *repoMock.MockOrderRepository
	inventoryClient *clientMock.MockInventoryClient
	paymentClient   *clientMock.MockPaymentClient
	service         *service
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()
	s.inventoryClient = clientMock.NewMockInventoryClient(s.T())
	s.paymentClient = clientMock.NewMockPaymentClient(s.T())
	s.repo = repoMock.NewMockOrderRepository(s.T())
	s.service = &service{
		inventoryClient: s.inventoryClient,
		paymentClient:   s.paymentClient,
		orderRepo:       s.repo,
	}
}

func (s *ServiceSuite) TearDownSuite() {
	s.inventoryClient.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
