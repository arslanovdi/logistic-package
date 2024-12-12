package test

import (
	"testing"
	"time"

	"github.com/arslanovdi/logistic-package/events/internal/process"
	"github.com/magiconair/properties/assert"
)

func TestTimeAvg(t *testing.T) {
	t.Parallel()

	const n = 20 // кол-во событий

	avg := process.TimeAvg{
		Duration: 1 * time.Second,
	}

	res := 0

	for range n {
		res = avg.Add()
	}
	assert.Equal(t, res, n) // 20 в секунду

	time.Sleep(1 * time.Second)

	res = avg.Add()
	assert.Equal(t, res, 1) // 1 в секунду
}
