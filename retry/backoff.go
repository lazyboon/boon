package retry

import (
	"math"
	"time"
)

type IBackoff interface {
	Next(count int) (stop bool, duration time.Duration)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ZeroBackoff struct {
	Count int
}

func NewZeroBackoff(count int) *ZeroBackoff {
	return &ZeroBackoff{Count: count}
}

func (z *ZeroBackoff) Next(cnt int) (bool, time.Duration) {
	return cnt > z.Count, 0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type AvgBackoff struct {
	Count int
	Time  time.Duration
}

func NewAvgBackoff(cnt int, t time.Duration) *AvgBackoff {
	return &AvgBackoff{Count: cnt, Time: t}
}

func (a *AvgBackoff) Next(cnt int) (bool, time.Duration) {
	return cnt > a.Count, a.Time
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ArithmeticBackoff Arithmetic Sequence backoff
type ArithmeticBackoff struct {
	Count int
	D     time.Duration
}

func NewArithmeticBackoff(cnt int, d time.Duration) *ArithmeticBackoff {
	return &ArithmeticBackoff{Count: cnt, D: d}
}

func (a *ArithmeticBackoff) Next(cnt int) (bool, time.Duration) {
	return cnt > a.Count, time.Duration(cnt-1) * a.D
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GeometricBackoff Geometric Sequence backoff
type GeometricBackoff struct {
	Count int
	First time.Duration
	Q     float64
}

func NewGeometricBackoff(cnt int, first time.Duration, q float64) *GeometricBackoff {
	return &GeometricBackoff{Count: cnt, First: first, Q: q}
}

func (g *GeometricBackoff) Next(cnt int) (bool, time.Duration) {
	return cnt > g.Count, g.First * time.Duration(math.Pow(g.Q, float64(cnt)))
}
