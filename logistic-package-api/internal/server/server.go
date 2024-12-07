// Package server - все веб серверы приложения
package server

import (
	"context"
	"fmt"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/service"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/api"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	pb "github.com/arslanovdi/logistic-package/pkg/logistic-package-api"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

// GrpcServer is gRPC server
type GrpcServer struct {
	server *grpc.Server
	lis    net.Listener
}

// grpcMiddleware Перехватчик унарных методов, считаем метрики
func grpcMiddleware(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	d := time.Now()

	m, err := handler(ctx, req)

	GRPC2.Observe(time.Since(d).Seconds())

	if status.Code(err) == codes.NotFound {
		GRPCNotFoundCounter.Inc()
	}

	CRUDCounter.Inc()

	return m, err
}

// NewGrpcServer returns gRPC server
func NewGrpcServer(packageService *service.PackageService) *GrpcServer {

	cfg := config.GetConfigInstance()

	s := &GrpcServer{}

	// дефолтные grpc метрики в прометеус
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120})))
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}

	s.server = grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(cfg.Grpc.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(cfg.Grpc.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(cfg.Grpc.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(cfg.Grpc.Timeout) * time.Minute,
		}),
		grpc.StatsHandler(otelgrpc.NewServerHandler()), // openTelemetry трассировка
		grpc.ChainUnaryInterceptor( // последовательное исполнение middleware, с общим контекстом
			grpcMiddleware,
			srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			grpcrecovery.UnaryServerInterceptor(), // дефолтный перехватчик паник
		),
	)

	pb.RegisterLogisticPackageApiServiceServer(s.server, api.NewPackageAPI(packageService)) // регистрируем имплементацию интерфейса в gRPC-сервере

	grpcPrometheus.EnableHandlingTimeHistogram()
	grpcPrometheus.Register(s.server)

	if cfg.Project.Debug {
		reflection.Register(s.server) // в дебаге регестрируем отражение методов gRPC-сервера: предоставляет сведения о публично доступных методах
	}

	return s
}

// Start method runs server
func (s *GrpcServer) Start() {

	log := slog.With("func", "GrpcServer.Start")

	cfg := config.GetConfigInstance()

	grpcAddr := fmt.Sprintf("%s:%v", cfg.Grpc.Host, cfg.Grpc.Port)

	var err1 error
	s.lis, err1 = net.Listen("tcp", grpcAddr)
	if err1 != nil {
		log.Error("failed to listen", slog.String("error", err1.Error()))
		os.Exit(1)
	}

	go func() {
		log.Info("GRPC Server is listening", slog.String("address", grpcAddr))
		if err2 := s.server.Serve(s.lis); err2 != nil {
			log.Error("Failed running gRPC server", slog.String("error", err2.Error()))
			os.Exit(1)
		}
	}()
}

// Stop - stop gRPC server
func (s *GrpcServer) Stop() error {

	s.server.GracefulStop()

	err := s.lis.Close()
	if err != nil {
		return err
	}
	return nil
}
