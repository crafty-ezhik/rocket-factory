module github.com/crafty-ezhik/rocket-factory/inventory

go 1.24.5

replace github.com/crafty-ezhik/rocket-factory/shared => ../shared

require (
	github.com/brianvoe/gofakeit/v7 v7.7.3
	github.com/crafty-ezhik/rocket-factory/shared v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	google.golang.org/grpc v1.76.0
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
