package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net/http"
	"os"
	"time"

	pb "github.com/arslanovdi/logistic-package/pkg/logistic-package-api"
)

var (
	httpTotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_microservice_requests_total",
		Help: "The total number of incoming HTTP requests",
	})
)

// GatewayServer is HTTP gRPC-gateway server
type GatewayServer struct {
	server *http.Server
}

// NewGatewayServer returns HTTP gRPC-gateway server
func NewGatewayServer() *GatewayServer {

	cfg := config.GetConfigInstance()
	grpcAddr := fmt.Sprintf("%s:%v", cfg.Grpc.Host, cfg.Grpc.Port)
	gatewayAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)

	log := slog.With("func", "server.NewGatewayServer")

	conn, err1 := grpc.NewClient(grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // openTelemtry трассировка grpc клиента
	)

	if err1 != nil {
		log.Warn("Failed to dial gRPC server",
			slog.String("error", err1.Error()))
	}

	rmux := runtime.NewServeMux()
	if err2 := pb.RegisterLogisticPackageApiServiceHandler(context.Background(), rmux, conn); err2 != nil {
		log.Warn("Failed registration handler",
			slog.String("error", err2.Error()))
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/",
		httpMetricWrapper( // оборачиваем HTTP методы, подсчитываем метрики, установка уровня логирования для запроса
			otelhttp.NewHandler(rmux, "grpc-gateway"), // Оборачиваем HTTP методы gRPC в openTelemtry трейсы
		),
	)

	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) { // Подменяем swagger.json, указанный в файле swagger-initializer.js сгенерированным logistic_package_api.swagger.json
		http.ServeFile(w, r, "./swagger/logistic_package_api.swagger.json")
	})

	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swagger-ui/"))))

	server := &http.Server{
		Addr:              gatewayAddr,
		Handler:           mux,
		ReadHeaderTimeout: time.Second * 5,
	}

	return &GatewayServer{
		server: server,
	}
}

/*
Start - starts the gateway server and Swagger server
cancelFunc - функция отмены контекста, вызывается в случае ошибки запуска
*/
func (s *GatewayServer) Start() {
	log := slog.With("func", "GatewayServer.Start")

	cfg := config.GetConfigInstance()

	gatewayAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)

	go func() {
		log.Info("Gateway server is running", slog.String("address", gatewayAddr))
		log.Info("Swagger server is running", slog.String("address", gatewayAddr+"/swagger-ui/"))
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed running gateway server", slog.String("error", err.Error()))
			os.Exit(1) // приложение завершается с ошибкой, при ошибке запуска сервера
		}
	}()
}

/*
Stop - stops the gateway server correctly
*/
func (s *GatewayServer) Stop(ctx context.Context) {
	log := slog.With("func", "GatewayServer.Stop")

	if err := s.server.Shutdown(ctx); err != nil {
		log.Error("GatewayServer.Shutdown", slog.String("error", err.Error()))
	} else {
		log.Info("GatewayServer shut down correctly")
	}
}

/*
httpMetricWrapper - обертка для http запросов
подсчет метрик
*/
func httpMetricWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		httpTotalRequests.Inc() // метрика

		h.ServeHTTP(w, r)
	})
}
