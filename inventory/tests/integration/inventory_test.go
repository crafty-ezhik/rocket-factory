//go:build integration

package integration

import (
	"context"

	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/grpc/errors"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
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

		ginkgo.It("Должно возвращать 404 Not Found", func() {
			expectedErr := errors.BusinessErrorToGRPCStatus(model.ErrPartNotFound).Err()

			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: uuid.NewString(),
			})

			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err).To(gomega.Equal(expectedErr))
			gomega.Expect(resp).To(gomega.BeNil())
		})

		ginkgo.It("Должно вернуть 400 Bad Request", func() {
			expectedErr := errors.BusinessErrorToGRPCStatus(model.ErrInvalidUUID).Err()

			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: "invalid-uuid",
			})

			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err).To(gomega.Equal(expectedErr))
			gomega.Expect(resp).To(gomega.BeNil())
		})
	})

	// Test 2: Получение списка деталей
	ginkgo.Describe("ListParts", func() {
		partUUIDs := []string{}

		// Добавляем детали
		ginkgo.BeforeEach(func() {
			partUUIDs = nil

			for range 3 {
				partUUID, err := env.InsertTestPart(ctx)
				gomega.Expect(err).ToNot(gomega.HaveOccurred(), "Ожидали успешно добавление детали в MongoDB")
				partUUIDs = append(partUUIDs, partUUID)
			}
		})

		ginkgo.It("Фильтры пустые. Должны вернуться все детали", func() {
			logger.Info(ctx, "🧪 Тест: Фильтры пустые. Должны вернуться все детали")
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: nil,
			})

			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(resp.Parts).ToNot(gomega.BeEmpty())
			gomega.Expect(len(resp.Parts)).ToNot(gomega.BeZero())
			gomega.Expect(len(resp.Parts)).To(gomega.BeNumerically("==", 3))
		})
		ginkgo.It("Передан имеющийся UUID. Должна вернуться 1 детали", func() {
			logger.Info(ctx, "🧪 Тест: Передан имеющийся UUID. Должна вернуться 1 детали")
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: []string{partUUIDs[0]},
				},
			})

			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(resp.Parts[0].GetUuid()).To(gomega.Equal(partUUIDs[0]))
			gomega.Expect(len(resp.Parts)).To(gomega.BeNumerically("==", 1))
		})
		ginkgo.It("Передан отсутствующий UUID. Должен вернуться пустой список", func() {
			logger.Info(ctx, "🧪 Тест: Передан отсутствующий UUID. Должен вернуться пустой список")
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: []string{uuid.NewString()},
				},
			})

			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(resp.Parts).To(gomega.BeEmpty())
		})
	})
})
