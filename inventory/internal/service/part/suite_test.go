package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	inventoryRepo *mocks.MockInventoryRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.inventoryRepo = mocks.NewMockInventoryRepository(s.T())
	s.service = NewService(s.inventoryRepo)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
