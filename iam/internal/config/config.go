package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/crafty-ezhik/rocket-factory/iam/internal/config/env"
)

var appConfig *config

type config struct {
	IamGRPC  GRPCConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Logger   LoggerConfig
	Session  SessionConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	iamGRPCConfig, err := env.NewIAMGRPCConfig()
	if err != nil {
		return err
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}
	redisConfig, err := env.NewRedisConfig()
	if err != nil {
		return err
	}
	sessionConfig, err := env.NewSessionConfig()
	if err != nil {
		return err
	}
	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		IamGRPC:  iamGRPCConfig,
		Postgres: postgresConfig,
		Redis:    redisConfig,
		Session:  sessionConfig,
		Logger:   loggerConfig,
	}
	return nil
}

func AppConfig() *config {
	return appConfig
}
