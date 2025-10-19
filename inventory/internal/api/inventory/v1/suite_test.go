package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/crafty-ezhik/rocket-factory/inventory/internal/service/mocks"
)

type ApiSuite struct {
	suite.Suite
	ctx              context.Context //nolint:containedctx
	inventoryService *mocks.MockInventoryService

	api *api
}

func (s *ApiSuite) SetupTest() {
	s.ctx = context.Background()
	s.inventoryService = mocks.NewMockInventoryService(s.T())

	s.api = NewAPI(s.inventoryService)
}

func (s *ApiSuite) TearDownTest() {}

func TestIntegrationAPI(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
