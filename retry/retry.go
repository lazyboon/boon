package retry

import "time"

type options struct {
	backoff IBackoff
}

type Option func(o *options)

func WithBackoff(backoff IBackoff) Option {
	return func(o *options) {
		o.backoff = backoff
	}
}

type TaskHandler func(count int) error

func Do(task TaskHandler, opts ...Option) error {
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
		if o.backoff != nil {
			stop, sleep := o.backoff.Next(i)
			if stop {
				break
			}
			time.Sleep(sleep)
		}
	}
	return err
}
