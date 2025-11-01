package config

type PaymentGRPCConfig interface {
	Address() string
}

type PaymentHTTPConfig interface {
	Address() string
}

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}
