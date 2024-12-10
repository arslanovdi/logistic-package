export GO111MODULE=on

.PHONY: build-tgbot
build-tgbot: .build-tgbot

.build-tgbot:
	go mod download && CGO_ENABLED=0  go build \
    		-o ./bin/tg_bot$(shell go env GOEXE) ./telegram_bot/cmd/bot/main.go
