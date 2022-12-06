package xamqp

type QOS struct {
	PrefetchCount int
	PrefetchSize  int
	Global        bool
}
