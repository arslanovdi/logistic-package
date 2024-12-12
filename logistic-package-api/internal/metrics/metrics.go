package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// GRPCNotFoundCounter - счетчик не найденных запросов
var GRPCNotFoundCounter = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "logistic",
	Subsystem: "package_api",
	Name:      "grpc_not_found",
	Help:      "Total gRPC not found calls",
})

// CRUDCounter - счетчик CRUD запросов
var CRUDCounter = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "logistic",
	Subsystem: "package_api",
	Name:      "crud",
	Help:      "Total CRUD calls",
})

// GRPC2 - гистограмма времени выполнения gRPC запросов
var GRPC2 = promauto.NewHistogram(prometheus.HistogramOpts{
	Namespace: "logistic",
	Subsystem: "package_api",
	Name:      "grpc2",
	Help:      "grpc2 calls",
}, // []string{"method"}, // метка для метрики, для каждой метки будет свой график
)

// RetranslatorEvents - счетчик событий которые сейчас отправляются в кафку
var RetranslatorEvents = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "logistic",
	Subsystem: "package_api",
	Name:      "retranslator",
	Help:      "Retranslator events in work",
})

// ProcessingTime - метрика времени обработки батча событий ретранслятором
var ProcessingTime = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "logistic_package",
	Subsystem: "events",
	Name:      "processing_time",
	Help:      "Processing time (seconds)",
})
