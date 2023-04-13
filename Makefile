.PHONY: build requirements api buildswagger

VERSION := $(shell git describe --always |sed -e "s/^v//")
GRPC_GW_PATH := $(shell go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway)
GOOGLEAPIS_PATH := "$(GRPC_GW_PATH)/../third_party/googleapis"

build: swagger
	mkdir -p build
	go build $(GO_EXTRA_BUILD_ARGS) -ldflags "-s -w -X main.version=$(VERSION)" -o build/xm cmd/app/main.go

requirements:
	@go mod download
	@go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	@go install github.com/golang/protobuf/protoc-gen-go
	@go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

api:
	protoc -I=$(GOOGLEAPIS_PATH) -I=protobuf -I=protobuf/internal/api/ --go_out=plugins=grpc,paths=source_relative:. internal/api/company.proto
	#grpc-gw
	protoc -I=$(GOOGLEAPIS_PATH) -I=protobuf -I=protobuf/internal/api --grpc-gateway_out=paths=source_relative,logtostderr=true:. internal/api/company.proto

swagger: requirements buildswagger static/swagger/api.swagger.json

buildswagger:
	@rm -rf static/swagger
	@mkdir -p static/swagger
	@cp -r ui/swagger/* static
	protoc -I=$(GOOGLEAPIS_PATH) -I=protobuf -I=protobuf/internal/api --swagger_out=json_names_for_fields=true:static/swagger internal/api/company.proto

static/swagger/api.swagger.json:
	@echo "Fetching Swagger definitions and generate combined Swagger JSON"

	@cp static/swagger/internal/api/*.json static/swagger
	@GOOS="" GOARCH="" go run internal/tools/swagger/main.go static/swagger/internal/api > static/swagger/api.swagger.json

test:
	@echo "Running tests"
	@rm -f coverage.out
	@golint ./...
	@go vet ./...
	@go test -p 1 -v -cover ./... -coverprofile coverage.out	