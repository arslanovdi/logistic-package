// Package process содержит логику обработки событий
package process

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/arslanovdi/logistic-package/events/internal/metric"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"go.opentelemetry.io/otel/trace"
)

var (
	packages map[string]struct{}
	avg      TimeAvg
)

// Инициализация глобальных переменных
func init() {
	packages = make(map[string]struct{})

	avg = TimeAvg{
		Duration: 1 * time.Minute,
	}
}

// PrintPackageEvent выводит событие на консоль
func PrintPackageEvent(ctx context.Context, packageID string, msg model.PackageEvent, offset int64) {
	log := slog.With("func", "process.PrintPackageEvent")

	if span := trace.SpanContextFromContext(ctx); span.IsSampled() { // вытягиваем span из контекста и пробрасываем в лог
		log = log.With("trace_id", span.TraceID().String())
	}

	_, ok := packages[packageID]
	if !ok {
		packages[packageID] = struct{}{}
		metric.PackageCounter.Inc() // счетчик уникальных пакетов
	}

	_, err := fmt.Printf("offset: %d. event: %s\n", offset, msg.String())
	if err != nil {
		log.Error("Error formating", slog.String("error", err.Error()))
	}

	log.Debug("PrintPackageEvent", slog.String("event", msg.String()))

	metric.EventsMinute.Set(float64(avg.Add())) // кол-во событий в минуту
}
