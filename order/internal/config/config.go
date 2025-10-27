package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/crafty-ezhik/rocket-factory/order/internal/config/env"
)

var appConfig *config

type config struct {
	OrderHTTP     OrderHTTPConfig
	Postgres      PostgresConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC   PaymentGRPCConfig
	Logger        LoggerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	orderHTTPConfig, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	inventoryGRPCConfig, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	paymentGRPCConfig, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		OrderHTTP:     orderHTTPConfig,
		Postgres:      postgresConfig,
		InventoryGRPC: inventoryGRPCConfig,
		PaymentGRPC:   paymentGRPCConfig,
		Logger:        loggerConfig,
	}
	return nil
}

func AppConfig() *config { return appConfig }
