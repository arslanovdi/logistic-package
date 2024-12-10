package process

import (
	"fmt"
	"github.com/arslanovdi/logistic-package/events/internal/metric"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"log/slog"
	"time"
)

var (
	packages map[string]struct{}
	avg      TimeAvg
)

func init() {

	packages = make(map[string]struct{})

	avg = TimeAvg{
		Duration: 1 * time.Minute,
	}
}

func PrintPackageEvent(packageID string, msg model.PackageEvent, offset int64) {

	log := slog.With("process.PrintPackageEvent")

	_, ok := packages[packageID]
	if !ok {
		packages[packageID] = struct{}{}
		metric.PackageCounter.Inc()
	}

	_, err := fmt.Printf("offset: %d. event: %s\n", offset, msg.String())

	if err != nil {
		log.Error("Error formating", slog.String("error", err.Error()))
	}

	metric.EventsMinute.Set(float64(avg.Add()))
}
