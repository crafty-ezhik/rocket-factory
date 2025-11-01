package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/crafty-ezhik/rocket-factory/payment/internal/config/env"
)

var appConfig *config

type config struct {
	PaymentGRPC PaymentGRPCConfig
	PaymentHTTP PaymentHTTPConfig
	Logger      LoggerConfig
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

	paymentGRPCConfig, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	paymentHTTPConfig, err := env.NewPaymentHTTPConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		PaymentGRPC: paymentGRPCConfig,
		PaymentHTTP: paymentHTTPConfig,
		Logger:      loggerCfg,
	}
	return nil
}

func AppConfig() *config { return appConfig }
