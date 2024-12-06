# syntax=docker/dockerfile:1-labs
#builder
FROM golang:1.23.3-bullseye AS builder

COPY ../Makefile /home/Makefile

COPY --parents ../go.mod /home/go.mod

COPY --parents ../go.sum /home/go.sum

WORKDIR /home

RUN make deps-go

COPY ../.. /home
# собираем
RUN make build-go

# gRPC Server

FROM golang:1.23.3-bullseye AS server

ARG GITHUB_PATH=github.com/arslanovdi/logistic-package/logistic-package-api/

LABEL org.opencontainers.image.source=https://${GITHUB_PATH}

WORKDIR /root/
# копируем все что нужно для работы приложения
COPY --from=builder /home/bin/grpc-server .
COPY --from=builder /home/config.yml .
COPY --from=builder /home/migrations/ ./migrations
COPY --from=builder /home/swagger ./swagger
COPY --from=builder /home/swagger-ui ./swagger-ui

RUN chown root:root grpc-server

EXPOSE 8080
EXPOSE 8082
EXPOSE 9100
EXPOSE 8000

CMD ["./grpc-server"]
