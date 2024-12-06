GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.22","$(shell printf "$(GO_VERSION_SHORT)\n1.22" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.22. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE=on

SERVICE_NAME=logistic-package-api
SERVICE_PATH=arslanovdi/logistic-package-api

PGV_VERSION:="v1.1.0"
BUF_VERSION:="v1.47.2"
GCC_VERSION="14.2.0"

OS_NAME=$(shell uname -s)
OS_ARCH=$(shell uname -m)
GO_BIN=$(shell go env GOPATH)/bin
BUF_EXE=$(GO_BIN)/buf$(shell go env GOEXE)

ifeq ("NT", "$(findstring NT,$(OS_NAME))")
OS_NAME=Linux
endif

.PHONY: run
run:
	go run cmd/grpc-server/main.go

.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: test
test:
	go test -v -race -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep total | awk '{print $$3}'


# ----------------------------------------------------------------

.PHONY: generate
generate-go: .generate-install-buf .generate-go .generate-finalize-go

.generate-install-buf:
	@ command -v buf 2>&1 > /dev/null || (echo "Install buf" && \
    		curl -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(OS_NAME)-$(OS_ARCH)$(shell go env GOEXE) --create-dirs -o "$(BUF_EXE)" && \
    		chmod +x "$(BUF_EXE)")
#.generate-install-buf:
#	scoop install buf	for windows

.generate-go:
	$(BUF_EXE) generate

.generate-finalize-go:
	#mv pkg/$(SERVICE_NAME)/github.com/$(SERVICE_PATH)/pkg/$(SERVICE_NAME)/* pkg/$(SERVICE_NAME)
	#rm -rf pkg/$(SERVICE_NAME)/github.com/
	cd pkg/$(SERVICE_NAME) && ls go.mod || (go mod init github.com/$(SERVICE_PATH)/pkg/$(SERVICE_NAME) && go mod tidy)

# ----------------------------------------------------------------

.PHONY: deps-go
deps-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.24.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.24.0
	go install github.com/envoyproxy/protoc-gen-validate@$(PGV_VERSION)


.PHONY: build-go
build-go: generate-go .build

.build:
	go mod download && CGO_ENABLED=1  go build \
		-tags='no_mysql no_sqlite3' \
		-ldflags=" \
			-X 'github.com/$(SERVICE_PATH)/internal/config.version=$(VERSION)' \
			-X 'github.com/$(SERVICE_PATH)/internal/config.commitHash=$(COMMIT_HASH)' \
		" \
		-o ./bin/grpc-server$(shell go env GOEXE) ./cmd/grpc-server/main.go

.PHONY: build-outbox
build-outbox: .build-outbox

.build-outbox:
	go mod download && CGO_ENABLED=1  go build \
    		-tags='no_mysql no_sqlite3' \
    		-ldflags=" \
    			-X 'github.com/$(SERVICE_PATH)/internal/config.version=$(VERSION)' \
    			-X 'github.com/$(SERVICE_PATH)/internal/config.commitHash=$(COMMIT_HASH)' \
    		" \
    		-o ./bin/retranslator$(shell go env GOEXE) ./cmd/retranslator/main.go
