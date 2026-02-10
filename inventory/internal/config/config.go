package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/crafty-ezhik/rocket-factory/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	InventoryGRPC InventoryGRPCConfig
	InventoryHTTP InventoryHTTPConfig
	Mongo         MongoConfig
	IamGRPC       IAMConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	inventoryHTTPCfg, err := env.NewInventoryHTTPConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	iamGRPCCfg, err := env.NewIAMGRPCConfig()
	if err != nil {
		return err
	}
	appConfig = &config{
		Logger:        loggerCfg,
		InventoryGRPC: inventoryGRPCCfg,
		InventoryHTTP: inventoryHTTPCfg,
		Mongo:         mongoCfg,
		IamGRPC:       iamGRPCCfg,
	}
	return nil
}

func AppConfig() *config { return appConfig }
