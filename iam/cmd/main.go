package main

import (
	"context"
	"fmt"
	authAPIV1 "github.com/crafty-ezhik/rocket-factory/iam/internal/api/auth/v1"
	userAPIV1 "github.com/crafty-ezhik/rocket-factory/iam/internal/api/user/v1"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/interceptor"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/repository/session"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/repository/user"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/service/auth"
	serviceUser "github.com/crafty-ezhik/rocket-factory/iam/internal/service/user"
	redisWrap "github.com/crafty-ezhik/rocket-factory/platform/pkg/cache/redis"
	sharedIns "github.com/crafty-ezhik/rocket-factory/platform/pkg/grpc/interceptors"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/hasher/bcrypt"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	authV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/auth/v1"
	userV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/user/v1"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"

	"github.com/crafty-ezhik/rocket-factory/iam/internal/config"
)

const configPath = "../deploy/compose/iam/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		log.Println(err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// === POSTGRES ===
	pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
	if err != nil {
		panic(fmt.Sprintf("❌ Ошибка подключения к базе данных: %v\n", err))
	}

	// Проверка соединения
	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	err = pool.Ping(pingCtx)
	if err != nil {
		panic(fmt.Sprintf("failed to ping Postgres: %v\n", err))
	}

	// === REDIS ===
	redisPool := &redigo.Pool{
		MaxIdle:     config.AppConfig().Redis.MaxIdle(),
		IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
		DialContext: func(ctx context.Context) (redigo.Conn, error) {
			return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
		},
	}

	// === ЗАВИСИМОСТИ ===
	redisClient := redisWrap.NewClient(redisPool, logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	hasher := bcrypt.NewBcryptPasswordHasher(10)

	userRepo := user.NewRepository(pool)
	sessionRepo := session.NewRepository(redisClient)
	userService := serviceUser.NewService(userRepo, hasher)
	authService := auth.NewService(userRepo, sessionRepo, hasher)
	apiUser := userAPIV1.NewUserAPI(userService)
	apiAuth := authAPIV1.NewAuthAPI(authService)

	// === СЕРВЕР ===
	server := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerInterceptor(),

			sharedIns.UnaryErrorInterceptor(),
		))

	reflection.Register(server)

	userV1.RegisterUserServiceServer(server, apiUser)
	authV1.RegisterAuthServiceServer(server, apiAuth)

	listener, err := net.Listen("tcp", config.AppConfig().IamGRPC.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("🚀 grpc server listening at %v\n", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
