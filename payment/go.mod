module github.com/crafty-ezhik/rocket-factory/payment

go 1.24.8

replace github.com/crafty-ezhik/rocket-factory/shared => ../shared

replace github.com/crafty-ezhik/rocket-factory/platform => ../platform

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/crafty-ezhik/rocket-factory/platform v0.0.0-00010101000000-000000000000
	github.com/crafty-ezhik/rocket-factory/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.11.1
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.76.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	go.opentelemetry.io/otel/sdk v1.38.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
