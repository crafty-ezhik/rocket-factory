package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	grpcInventoryAddr = "localhost:50052"
	grpcPaymentAddr   = "localhost:50051"
)

var (
	ErrNotFound    = errors.New("order not found")
	ErrOrderIsPaid = errors.New("order is paid")
)

// OrderStorage - хранилище для заказов
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

// GetOrder - получает заказ из хранилища по UUID
func (s *OrderStorage) GetOrder(orderUUID string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return nil
	}

	return order
}

// CreateOrder - добавляет новый заказ в хранилище
func (s *OrderStorage) CreateOrder(user uuid.UUID, parts []uuid.UUID, totalPrice float64) (uuid.UUID, error) {
	orderUUID := uuid.New()
	newOrder := &orderV1.OrderDto{
		OrderUUID:       orderUUID,
		UserUUID:        user,
		PartUuids:       parts,
		TotalPrice:      totalPrice,
		TransactionUUID: orderV1.OptNilUUID{},
		PaymentMethod:   orderV1.OptNilPaymentMethod{},
		Status:          orderV1.OrderStatusPENDINGPAYMENT,
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[orderUUID.String()] = newOrder

	return orderUUID, nil
}

// UpdateOrder - меняет статут оплаты, метод оплаты и уникальный идентификатор транзакции
func (s *OrderStorage) UpdateOrder(orderUUID, transactionUUID uuid.UUID, paymentMethod orderV1.PaymentMethod) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	order := s.orders[orderUUID.String()]
	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = orderV1.OptNilUUID{Value: transactionUUID}
	order.PaymentMethod = orderV1.OptNilPaymentMethod{Value: paymentMethod}
	return nil
}

// CancelOrder - меняет статут заказа на CANCELLED или возвращает ошибку, если он имеет статут PAID
func (s *OrderStorage) CancelOrder(orderUUID string) error {
	// Ищем заказ по переданному UUID
	s.mu.RLock()
	order, ok := s.orders[orderUUID]
	if !ok {
		return ErrNotFound
	}
	s.mu.RUnlock()

	// Проверяем статус на PAID
	if order.Status == orderV1.OrderStatusPAID {
		return ErrOrderIsPaid
	}

	// Меняем статус заказа на CANCELLED
	order.Status = orderV1.OrderStatusCANCELLED

	return nil
}

// OrderHandler - реализует интерфейс Handler для работы с заказами
type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewOrderHandler(
	storage *OrderStorage,
	grpcInventory inventoryV1.InventoryServiceClient,
	grpcPayment paymentV1.PaymentServiceClient,
) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: grpcInventory,
		paymentClient:   grpcPayment,
	}
}

// OrderCancel - отменяет заказ
func (h *OrderHandler) OrderCancel(ctx context.Context, req orderV1.OrderCancelParams) (orderV1.OrderCancelRes, error) {
	if err := uuid.Validate(req.OrderUUID); err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "uuid validation failed. " + err.Error(),
		}, nil
	}

	err := h.storage.CancelOrder(req.OrderUUID)
	if errors.Is(err, ErrNotFound) {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("order %s not found", req.OrderUUID),
		}, nil
	}

	if errors.Is(err, ErrOrderIsPaid) {
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "a paid order cannot be canceled",
		}, nil
	}

	return &orderV1.OrderCancelNoContent{}, nil
}

// OrderCreate - Создает заказ и возвращает его уникальный идентификатор
func (h *OrderHandler) OrderCreate(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.OrderCreateRes, error) {
	if err := validateCreateRequest(req); err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	partsUUID, err := convertUUIDStoString(req)
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	// Идем в InventoryService для получения списка запрашиваемых деталей
	parts, err := h.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{Filter: &inventoryV1.PartsFilter{
		Uuids: partsUUID,
	}})
	if err != nil {
		log.Printf("ListParts failed: %v", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "internal server error",
		}, nil
	}

	totalPrice, err := checkPartsAndCountTotalPrice(partsUUID, parts)
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	orderUUID, err := h.storage.CreateOrder(req.GetUserUUID(), req.GetPartUuids(), totalPrice)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "internal server error",
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

// OrderGet - Возвращает данные по заказу
func (h *OrderHandler) OrderGet(ctx context.Context, req orderV1.OrderGetParams) (orderV1.OrderGetRes, error) {
	// Валидируем uuid
	if err := uuid.Validate(req.OrderUUID); err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "uuid validation failed. " + err.Error(),
		}, nil
	}

	order := h.storage.GetOrder(req.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}
	return order, nil
}

// OrderPay - Выполняет оплату заказа
func (h *OrderHandler) OrderPay(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.OrderPayParams) (orderV1.OrderPayRes, error) {
	// Валидируем uuid
	if err := uuid.Validate(params.OrderUUID); err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "uuid validation failed. " + err.Error(),
		}, nil
	}

	// Получаем заказ из хранилища
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	// Преобразуем orderV1.PaymentMethod к string для дальнейшей работы
	paymentMethod, ok := req.GetPaymentMethod().Get()
	if !ok {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "Payment method not found",
		}, nil
	}

	paymentMethodText, err := paymentMethod.MarshalText()
	if err != nil {
		log.Printf("MarshalText failed: %v", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "internal server error",
		}, nil
	}

	// Отправляем данные для оплаты в PaymentService
	transactionUUIDstr, err := h.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     params.OrderUUID,
		UserUuid:      order.GetUserUUID().String(),
		PaymentMethod: paymentV1.PaymentMethod(paymentV1.PaymentMethod_value[string(paymentMethodText)]),
	})
	if err != nil {
		log.Printf("PayOrder failed: %v", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "internal server error",
		}, nil
	}

	// Преобразуем string к UUID для добавления к заказу
	transactionUUID, err := uuid.Parse(transactionUUIDstr.GetTransactionUuid())
	if err != nil {
		log.Printf("ParseTransactionUUID failed: %v", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "internal server error",
		}, nil
	}

	// Обновляем данные по заказу
	err = h.storage.UpdateOrder(order.GetOrderUUID(), transactionUUID, paymentMethod)
	if err != nil {
		log.Printf("PayOrder failed: %v", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "internal server error",
		}, nil
	}

	return &orderV1.PayOrderResponse{TransactionUUID: transactionUUID}, nil
}

// NewError создает новую ошибку в формате GenericError
func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

func convertUUIDStoString(req *orderV1.CreateOrderRequest) ([]string, error) {
	partsUUID := make([]string, 0, len(req.GetPartUuids()))
	for _, partUUID := range req.GetPartUuids() {
		if err := uuid.Validate(partUUID.String()); err != nil {
			return nil, errors.New("uuid validation failed")
		}
		partsUUID = append(partsUUID, partUUID.String())
	}
	return partsUUID, nil
}

func validateCreateRequest(req *orderV1.CreateOrderRequest) error {
	if err := req.Validate(); err != nil {
		return errors.New("validate CreateOrder Request failed")
	}

	if err := uuid.Validate(req.GetUserUUID().String()); err != nil {
		return errors.New("uuid validation failed")
	}
	return nil
}

func checkPartsAndCountTotalPrice(partsUUID []string, parts *inventoryV1.ListPartsResponse) (float64, error) {
	totalPrice := 0.0
	for _, UUID := range partsUUID {
		exist := false
		for _, part := range parts.GetParts() {
			if part.Uuid == UUID {
				exist = true
				totalPrice += part.Price
				break
			}
		}
		if !exist {
			return 0, errors.New(fmt.Sprintf("part with uuid %s not found", UUID))
		}
	}
	return totalPrice, nil
}
func main() {
	// Создаем хранилище для данных о заказах
	storage := NewOrderStorage()

	// Создаем gRPC клиента для InventoryService
	inventoryConn, err := grpc.NewClient(
		grpcInventoryAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}

	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	// Создаем gRPC клиента для PaymentService
	paymentConn, err := grpc.NewClient(
		grpcPaymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}

	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	// Создаем обработчик для API
	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
		return
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчик OpenAPI к нашему серверу
	r.Mount("/", orderServer)

	// Создаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Запускаем сервер
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Реализуем Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
