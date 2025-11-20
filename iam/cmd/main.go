package main

import (
	"context"
	"fmt"
	v1 "github.com/crafty-ezhik/rocket-factory/iam/internal/api/user/v1"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/interceptor"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/repository/user"
	user2 "github.com/crafty-ezhik/rocket-factory/iam/internal/service/user"
	sharedIns "github.com/crafty-ezhik/rocket-factory/platform/pkg/grpc/interceptors"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/hasher/bcrypt"
	userV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/user/v1"
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

	userRepo := user.NewRepository(pool)
	hasher := bcrypt.NewBcryptPasswordHasher(10)
	userService := user2.NewService(userRepo, hasher)
	api := v1.NewUserAPI(userService)

	server := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerInterceptor(),

			sharedIns.UnaryErrorInterceptor(),
		))

	reflection.Register(server)

	userV1.RegisterUserServiceServer(server, api)

	listener, err := net.Listen("tcp", config.AppConfig().IamGRPC.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("🚀 grpc server listening at %v\n", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
