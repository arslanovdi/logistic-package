// Package metric метрики prometheus
package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PackageCounter - счетчик не найденных запросов
var PackageCounter = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "logistic_package",
	Subsystem: "events",
	Name:      "unique_packages",
	Help:      "Total unique packages",
})

// EventsMinute - кол-во событий в минуту, полученных из кафки
var EventsMinute = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "logistic_package",
	Subsystem: "events",
	Name:      "events_minute",
	Help:      "Events per last minute",
})
