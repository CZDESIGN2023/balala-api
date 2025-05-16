GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

API_SWAGGER_SCAN_DIR = .
API_SWAGGER_SCAN_ENTRY = cmd/go-cs/main.go
API_SWAGGER_OUT_DIR = docs

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
	API_PROTO_BEAN_FILES=$(shell find internal/bean -name *.proto)
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
	API_PROTO_BEAN_FILES=$(shell find internal/bean -name *.proto)
endif

.PHONY: init
# init env
init:
	go env -w GOPROXY=https://goproxy.cn,direct
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: internal
# generate internal proto
internal:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	make bean
	protoc --proto_path=./api \
	       --proto_path=./third_party \
	       --proto_path=./internal \
 	       --go_out=paths=source_relative:./api \
 	       --go-http_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
 	       --openapi_out=$(API_SWAGGER_OUT_DIR) \
	       $(API_PROTO_FILES)

.PHONY: bean
# generate bean proto
bean:
	go mod tidy
	go run ./cmd/generate_apply_update_block/main.go
	make internal

.PHONY: proto
# generate all proto
proto:
	make api
	make bean

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate wire
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

#生成doc文档
doc:
	#swag init -g cmd/go-cs/main.go -o docs
	swag fmt -d ${API_SWAGGER_SCAN_DIR} -g ${API_SWAGGER_SCAN_ENTRY}
	swag init -d ${API_SWAGGER_SCAN_DIR} -g ${API_SWAGGER_SCAN_ENTRY} -o ${API_SWAGGER_OUT_DIR} --parseInternal

.PHONY: all
# generate all
all:
	make proto;
	make generate;
	

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

start-doc:
	docker run --rm -d --name balala_doc -p 80:8080 -e SWAGGER_JSON=/path/to/openapi.yaml -v ./docs/openapi.yaml:/path/to/openapi.yaml swaggerapi/swagger-ui:v5.18.3

stop-doc:
	docker stop balala_doc

restart-doc:
	docker restart balala_doc