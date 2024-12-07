package process

import (
	"fmt"
	"github.com/arslanovdi/logistic-package/pkg/model"
)

var packages map[string]struct{}

func init() {
	packages = make(map[string]struct{})
}

func PrintPackageEvent(packageID string, msg model.PackageEvent, offset int64) {
	_, ok := packages[packageID]
	if !ok {
		packages[packageID] = struct{}{}
		//metrics.PackageCounter.Inc()
	}
	fmt.Println(msg, offset)
}
