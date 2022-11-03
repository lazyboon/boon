package retry

import "time"

type options struct {
	count   int
	backoff IBackoff
}

type Option func(o *options)

func WithCount(c int) Option {
	return func(o *options) {
		o.count = c
	}
}

func WithBackoff(backoff IBackoff) Option {
	return func(o *options) {
		o.backoff = backoff
	}
}

func Do(task func(c int) error, opts ...Option) error {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	var err error
	for i := 1; ; i++ {
		err = task(i)
		if err == nil {
			return nil
		}
		if i > o.count {
			break
		}
		var sleep time.Duration
		if o.backoff != nil {
			sleep = o.backoff.Next(i)
		}
		time.Sleep(sleep)
	}
	return err
}
