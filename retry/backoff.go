package retry

import "time"

type IBackoff interface {
	Next(count uint) (stop bool, duration time.Duration)
}

type ZeroBackoff struct{}

type UnlimitedBackoff struct{}
