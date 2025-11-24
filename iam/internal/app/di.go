package app

import (
	"context"
	"fmt"
	authAPIV1 "github.com/crafty-ezhik/rocket-factory/iam/internal/api/auth/v1"
	userAPIV1 "github.com/crafty-ezhik/rocket-factory/iam/internal/api/user/v1"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/config"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/repository"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/repository/session"
	userRepo "github.com/crafty-ezhik/rocket-factory/iam/internal/repository/user"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/service"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/service/auth"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/service/user"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/cache"
	redisWrap "github.com/crafty-ezhik/rocket-factory/platform/pkg/cache/redis"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/closer"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/hasher"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/hasher/bcrypt"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	authV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/auth/v1"
	userV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/user/v1"
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type diContainer struct {
	authV1API         authV1.AuthServiceServer
	userV1API         userV1.UserServiceServer
	authService       service.AuthService
	userService       service.UserService
	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository
	redisClient       cache.RedisClient
	hasher            hasher.PasswordHasher
	pgConnPool        *pgxpool.Pool
	redisConn         *redis.Pool
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = authAPIV1.NewAuthAPI(d.AuthService(ctx))
	}
	return d.authV1API
}

func (d *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userAPIV1.NewUserAPI(d.UserService(ctx))
	}
	return d.userV1API
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = auth.NewService(d.UserRepository(ctx), d.SessionRepository(ctx), d.Hasher(ctx))
	}
	return d.authService
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = user.NewService(d.UserRepository(ctx), d.Hasher(ctx))
	}
	return d.userService
}

func (d *diContainer) SessionRepository(ctx context.Context) repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = session.NewRepository(d.RedisClient(ctx))
	}
	return d.sessionRepository
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepo.NewRepository(d.PgConn(ctx))
	}
	return d.userRepository
}

func (d *diContainer) RedisClient(ctx context.Context) cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redisWrap.NewClient(
			d.RedisConn(ctx),
			logger.Logger(),
			config.AppConfig().Redis.ConnectionTimeout(),
		)
	}
	return d.redisClient
}

func (d *diContainer) Hasher(_ context.Context) hasher.PasswordHasher {
	if d.hasher == nil {
		d.hasher = bcrypt.NewBcryptPasswordHasher(10)
	}
	return d.hasher
}

func (d *diContainer) PgConn(ctx context.Context) *pgxpool.Pool {
	if d.pgConnPool == nil {
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

		// Добавляем закрытие пула в closer
		closer.AddNamed("Postgres connection pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pgConnPool = pool
	}
	return d.pgConnPool
}

func (d *diContainer) RedisConn(ctx context.Context) *redis.Pool {
	if d.redisConn == nil {
		redisPool := &redis.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redis.Conn, error) {
				return redis.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}

		closer.AddNamed("Redis connection pool", func(ctx context.Context) error {
			if err := redisPool.Close(); err != nil {
				logger.Error(ctx, "❌ Ошибка закрытия подключения с Redis", zap.Error(err))
				return err
			}
			return nil
		})

		d.redisConn = redisPool

	}
	return d.redisConn
}
