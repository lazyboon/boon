package retry

import (
	"math"
	"time"
)

type IBackoff interface {
	Next(count int) time.Duration
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ZeroBackoff struct{}

func NewZeroBackoff() *ZeroBackoff {
	return &ZeroBackoff{}
}

func (z *ZeroBackoff) Next(_ int) time.Duration {
	return 0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type AvgBackoff struct {
	Time time.Duration
}

func NewAvgBackoff(t time.Duration) *AvgBackoff {
	return &AvgBackoff{Time: t}
}

func (z *AvgBackoff) Next(_ int) time.Duration {
	return z.Time
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ArithmeticBackoff Arithmetic Sequence backoff
type ArithmeticBackoff struct {
	D time.Duration
}

func NewArithmeticBackoff(d time.Duration) *ArithmeticBackoff {
	return &ArithmeticBackoff{D: d}
}

func (a *ArithmeticBackoff) Next(count int) time.Duration {
	return time.Duration(count-1) * a.D
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GeometricBackoff Geometric Sequence backoff
type GeometricBackoff struct {
	First time.Duration
	Q     float64
}

func NewGeometricBackoff(first time.Duration, q float64) *GeometricBackoff {
	return &GeometricBackoff{First: first, Q: q}
}

func (g *GeometricBackoff) Next(count int) time.Duration {
	if count < 0 {
		count = 0
	}
	return g.First * time.Duration(math.Pow(g.Q, float64(count)))
}
