export GO111MODULE=on

PGV_VERSION:="v1.2.1"
BUF_VERSION:="v1.53.0"

OS_NAME=$(shell uname -s)
OS_ARCH=$(shell uname -m)
GO_BIN=$(shell go env GOPATH)/bin
BUF_EXE=$(GO_BIN)/buf$(shell go env GOEXE)

ifeq ("NT", "$(findstring NT,$(OS_NAME))")
OS_NAME=Linux
endif

.PHONY: generate
generate-go: .generate-install-buf .generate-go

.generate-install-buf:
	@ command -v buf 2>&1 > /dev/null || (echo "Install buf" && \
    		curl -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(OS_NAME)-$(OS_ARCH)$(shell go env GOEXE) --create-dirs -o "$(BUF_EXE)" && \
    		chmod +x "$(BUF_EXE)")

.generate-go:
	cd protobuf && \
	$(BUF_EXE) generate

# ----------------------------------------------------------------

.PHONY: deps-go
deps-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.26.3
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.26.3
	go install github.com/envoyproxy/protoc-gen-validate@$(PGV_VERSION)

.PHONY: build-go
build-go: generate-go .build

.build:
	go mod download && CGO_ENABLED=0  go build \
		-o ./bin/grpc-server$(shell go env GOEXE) ./logistic-package-api/cmd/grpc-server/main.go

.PHONY: build-outbox
build-outbox: .build-outbox

.build-outbox:
	go mod download && CGO_ENABLED=1  go build -tags musl \
    		-o ./bin/retranslator$(shell go env GOEXE) ./logistic-package-api/cmd/retranslator/main.go
