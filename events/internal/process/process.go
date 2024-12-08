package process

import (
	"fmt"
	"github.com/arslanovdi/logistic-package/events/internal/metric"
	"github.com/arslanovdi/logistic-package/pkg/model"
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
	_, ok := packages[packageID]
	if !ok {
		packages[packageID] = struct{}{}
		metric.PackageCounter.Inc()
	}

	fmt.Printf("offset: %d. event: %s\n", offset, msg.String())

	metric.EventsMinute.Set(float64(avg.Add()))
}
