module github.com/crafty-ezhik/rocket-factory/inventory

go 1.24.5

replace github.com/crafty-ezhik/rocket-factory/shared => ../shared

replace github.com/crafty-ezhik/rocket-factory/platform => ../platform

require (
	github.com/brianvoe/gofakeit/v7 v7.7.3
	github.com/caarlos0/env/v11 v11.3.1
	github.com/crafty-ezhik/rocket-factory/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.11.1
	go.mongodb.org/mongo-driver v1.17.4
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251007200510-49b9836ed3ff // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
