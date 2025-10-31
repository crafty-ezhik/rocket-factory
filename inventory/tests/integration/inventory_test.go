//go:build integration

package integration

import (
	"context"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

// Описываем наш сервис
var _ = ginkgo.Describe("Inventory Service", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	// Что необходимо делать перед каждым тестом
	ginkgo.BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаем gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	// Что необходимо сделать после каждого теста
	ginkgo.AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Ожидали успешную очистку коллекции")

		cancel()
	})

	// Test 1: Получение детали
	ginkgo.Describe("Get", func() {
		var partUUID string

		// Добавляем деталь
		ginkgo.BeforeEach(func() {
			var err error
			partUUID, err = env.InsertTestPart(ctx)
			gomega.Expect(err).ToNot(gomega.HaveOccurred(), "Ожидали успешно добавление детали в MongoDB")
		})

		ginkgo.It("Должно успешно возвращать деталь по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUUID,
			})
			part := resp.GetPart()

			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(part.GetUuid()).To(gomega.Equal(partUUID))
			gomega.Expect(part.GetName()).ToNot(gomega.BeEmpty())
			gomega.Expect(part.GetCategory()).To(gomega.BeNumerically(">", 0))
			gomega.Expect(part.GetDescription()).ToNot(gomega.BeEmpty())
			gomega.Expect(part.GetManufacturer()).ToNot(gomega.BeNil())
			gomega.Expect(part.GetDimensions()).ToNot(gomega.BeNil())
			gomega.Expect(part.GetPrice()).To(gomega.BeNumerically(">", 0))
			gomega.Expect(part.GetTags()).ToNot(gomega.BeEmpty())
		})
	})
})
