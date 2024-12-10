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
	count    int
	Duration time.Duration
}

// Add увеличивает счетчик событий с ограничением по времени
func (c *TimeAvg) Add() int {
	if c.root == nil {
		c.root = &record{
			time: time.Now(),
		}
	} else {
		// go to end
		end := c.root
		for end.next != nil {
			end = end.next
		}
		end.next = &record{
			time: time.Now(),
		}
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
