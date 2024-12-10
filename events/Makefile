export GO111MODULE=on

.PHONY: build-events
build-events: .build-events

.build-events:
	go mod download && CGO_ENABLED=1  go build \
    		-o ./bin/events$(shell go env GOEXE) ./events/cmd/events/main.go
