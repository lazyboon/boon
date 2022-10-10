package xredis

import "time"

type IBackoff interface {
	Next(count uint) (stop bool, duration time.Duration)
}

type ZeroBackoff struct{}

func (ZeroBackoff) Next(_ uint) (stop bool, duration time.Duration) {
	return true, 0
}
