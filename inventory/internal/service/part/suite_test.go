package part

import (
	"context"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

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
