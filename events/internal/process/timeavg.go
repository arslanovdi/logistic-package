package process

import (
	"time"
)

type record struct {
	time time.Time
	next *record
}

// TimeAvg - счетчик событий с ограничением по времени
type TimeAvg struct {
	root     *record
	end      *record
	count    int
	Duration time.Duration
}

// Add увеличивает счетчик событий с ограничением по времени
func (c *TimeAvg) Add() int {
	if c.root == nil {
		c.root = &record{
			time: time.Now(),
		}
		c.end = c.root
	} else {
		c.end.next = &record{
			time: time.Now(),
		}
		c.end = c.end.next
	}
	c.count++

	c.actualize()

	return c.count
}

// Удаляет события посчитанные раньше, чем Duration назад
func (c *TimeAvg) actualize() {
	pointer := c.root
	for pointer != nil && pointer.time.Before(time.Now().Add(-c.Duration)) {
		pointer = pointer.next
		c.count--
	}
	c.root = pointer
}
