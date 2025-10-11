package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/crafty-ezhik/rocket-factory/inventory/internal/interceptor"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

const (
	HOST          = "localhost"
	PathToSwagger = "./shared/pkg/swagger/inventory/v1"
	grpcPort      = 50052
	httpPort      = 8082
)

// mapParts -псевдоним для map[string]*inventoryV1.Part
type mapParts map[string]*inventoryV1.Part

// FilterType - ограничение для типов фильтра
type FilterType interface {
	string | inventoryV1.Category
}

// inventoryService - реализует gRPC сервис для работы с оплатами
type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	store mapParts
}

// applyFilter - применяет фильтр к мапе. Передается список значений фильтра по которому фильтровать
// и геттер для получения значения у n-го элемента
func applyFilter[T FilterType](mp mapParts, filter []T, getter func(p *inventoryV1.Part) T) {
	for uuidPart, part := range mp {
		exist := false
		for _, item := range filter {
			if getter(part) == item {
				exist = true
				break
			}
		}
		if !exist {
			delete(mp, uuidPart)
		}
	}
}

// map2slice - преобразует мапу в слайс для ответа клиенту
func (is *inventoryService) map2slice(data map[string]*inventoryV1.Part) []*inventoryV1.Part {
	sliceParts := make([]*inventoryV1.Part, 0, len(is.store))

	is.mu.RLock()
	defer is.mu.RUnlock()

	for _, v := range data {
		sliceParts = append(sliceParts, v)
	}
	return sliceParts
}

// GetPart - возвращает деталь по переданном уникальному идентификатору
func (is *inventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	is.mu.RLock()
	defer is.mu.RUnlock()

	item, ok := is.store[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "item not found")
	}
	return &inventoryV1.GetPartResponse{Part: item}, nil
}

// ListParts - возвращает список деталей по указанным фильтрам
func (is *inventoryService) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	tempMap := make(mapParts)
	is.mu.RLock()
	for k, v := range is.store {
		tempMap[k] = v
	}
	is.mu.RUnlock()

	filter := req.GetFilter()

	if filter == nil {
		return &inventoryV1.ListPartsResponse{Parts: is.map2slice(is.store)}, nil
	}

	if filter.GetUuids() != nil {
		applyFilter[string](tempMap, filter.GetUuids(), func(p *inventoryV1.Part) string { return p.GetUuid() })
	}

	if filter.GetNames() != nil {
		applyFilter[string](tempMap, filter.GetNames(), func(p *inventoryV1.Part) string { return p.GetName() })
	}

	if filter.GetCategories() != nil {
		applyFilter[inventoryV1.Category](
			tempMap,
			filter.GetCategories(),
			func(p *inventoryV1.Part) inventoryV1.Category { return p.GetCategory() },
		)
	}

	if filter.GetManufacturerCountries() != nil {
		applyFilter[string](
			tempMap,
			filter.GetManufacturerCountries(),
			func(p *inventoryV1.Part) string { return p.GetManufacturer().GetCountry() },
		)
	}

	if filter.GetTags() != nil {
		for uuidPart, part := range tempMap {
			exist := false
			for _, tag := range filter.GetTags() {
				for _, partTag := range part.GetTags() {
					if partTag == tag {
						exist = true
						break
					}
				}
			}
			if !exist {
				delete(tempMap, uuidPart)
			}
		}
	}

	return &inventoryV1.ListPartsResponse{
		Parts: is.map2slice(tempMap),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if err := lis.Close(); err != nil {
			log.Fatalf("failed to close listener: %v\n", err)
		}
	}()

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerInterceptor(),
			interceptor.ValidatorInterceptor(),
		),
	)

	// Регистрируем сервис inventoryService
	service := &inventoryService{store: generateFakeData(10)}
	inventoryV1.RegisterInventoryServiceServer(grpcServer, service)

	// Включаем рефлексию для отладки
	reflection.Register(grpcServer)

	// Запускаем сервер
	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Поднимаем HTTP сервер для gRPC-gateway + Swagger UI
	var gwServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Создаем мультиплексор для запросов
		mux := runtime.NewServeMux()

		// Настраиваем опции для соединения с gRPC сервером
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// Регистрируем gRPC-gateway хендлеры
		err = inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("%s:%d", HOST, grpcPort),
			opts)
		if err != nil {
			log.Printf("failed to register gateway: %v\n", err)
			return
		}

		// Создаем файловый сервер для swagger-ui
		fileServer := http.FileServer(http.Dir(PathToSwagger))

		// Создаем HTTP маршрутизатор
		httpMux := http.NewServeMux()

		// Регистрируем API ручку
		httpMux.Handle("/api/v1/inventory/", mux)

		// Swagger UI ручки
		httpMux.Handle("/swagger-ui.html", fileServer)
		httpMux.Handle("/inventory.swagger.json", fileServer)

		// Настраиваем редирект с корня на Swagger UI
		httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
				return
			}
			fileServer.ServeHTTP(w, r)
		}))

		// Создаем HTTP сервер
		gwServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           httpMux,
			ReadHeaderTimeout: 10 * time.Second,
		}

		// Запускаем HTTP сервер
		log.Printf("🌐 HTTP server with gRPC-Gateway and Swagger UI listening on %d\n", httpPort)
		err = gwServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")

	// Останавливаем HTTP сервер
	if gwServer != nil {
		shutdownctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := gwServer.Shutdown(shutdownctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		log.Println("✅ HTTP server stopped")
	}

	// Останавливаем gRPC сервер
	grpcServer.GracefulStop()
	log.Println("✅ gRPC Server stopped")
}

func generateFakeData(n int) mapParts {
	fakeData := make(mapParts)
	catSlice := []inventoryV1.Category{
		inventoryV1.Category_ENGINE,
		inventoryV1.Category_FUEL,
		inventoryV1.Category_PORTHOLE,
		inventoryV1.Category_WING,
		inventoryV1.Category_UNKNOW_UNSPECIFIED,
	}

	for range n {
		data := &inventoryV1.Part{
			Uuid:          uuid.NewString(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.HackerPhrase(),
			Price:         math.Floor(gofakeit.Float64Range(1, 1000)*100) / 100,
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      catSlice[gofakeit.IntRange(0, len(catSlice)-1)],
			Dimensions: &inventoryV1.Dimensions{
				Length: gofakeit.Float64Range(1, 10000),
				Width:  gofakeit.Float64Range(1, 10000),
				Height: gofakeit.Float64Range(1, 10000),
				Weight: gofakeit.Float64Range(1, 10000),
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    gofakeit.Name(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags:      []string{gofakeit.Word(), gofakeit.Word(), gofakeit.Word()},
			Metadata:  nil,
			CreatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix(), Nanos: 0},
			UpdatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix(), Nanos: 0},
		}

		fakeData[data.GetUuid()] = data
	}
	return fakeData
}
