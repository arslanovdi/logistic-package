GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.22","$(shell printf "$(GO_VERSION_SHORT)\n1.22" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.22. Found: $(GO_VERSION_SHORT))
endif

.PHONY: run
run:
	go run cmd/bot/main.go

.PHONY: build
build:
	go build -o bot cmd/bot/main.go

.PHONY: build-events
build-events: .build-events

.build-events:
	go mod download && CGO_ENABLED=1  go build \
    		-o ./bin/events$(shell go env GOEXE) ./events/cmd/events/main.go
