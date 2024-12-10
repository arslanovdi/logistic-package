// Package grpc предоставляет функции для работы с gRPC сервером
package grpc

import (
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"os"
)

import pb "github.com/arslanovdi/logistic-package/pkg/logistic-package-api"

// Client - GrpcClient
type Client struct {
	send pb.LogisticPackageApiServiceClient
	conn *grpc.ClientConn
}

// NewGrpcClient инициализирует соединение с gRPC сервером
func NewGrpcClient() *Client {

	log := slog.With("func", "GrpcClient.NewGrpcClient")

	cfg := config.GetConfigInstance()

	// подключение к grpc серверу без TLS
	conn, err := grpc.NewClient(
		cfg.GRPC.Host+":"+cfg.GRPC.Port,
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // Трассировка
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Warn("did not connect", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("gRPC client connected", slog.Any("address", cfg.GRPC.Host+":"+cfg.GRPC.Port))
	return &Client{
		send: pb.NewLogisticPackageApiServiceClient(conn), // инициализируем интерфейс через который будут вызываться удаленные методы
		conn: conn,
	}
}

// Close закрывает соединение с gRPC сервером
func (client *Client) Close() {
	log := slog.With("func", "GrpcClient.Close")

	err := client.conn.Close()
	if err != nil {
		log.Warn("did not close gRPC connection", slog.String("error", err.Error()))
		return
	}

	log.Info("gRPC client disconnected")
}
