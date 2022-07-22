package xredis

import "time"

var (
	WithConfig     = withConfig{}
	WithLock       = withLock{}
	WithDelayQueue = withDelayQueue{}
)

//----------------------------------------------------------------------------------------------------------------------

type config struct {
	alias        string
	host         string
	port         uint
	username     string
	password     string
	db           uint
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
	poolTimeout  time.Duration
	maxConnAge   time.Duration
	idleTimeout  time.Duration
	poolSize     uint
	minIdleConn  uint
}

func newConfig(options ...ConfigOption) *config {
	c := &config{
		alias: "",
		host:  "127.0.0.1",
		port:  6379,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

type ConfigOption func(c *config)

type withConfig struct{}

func (withConfig) Alias(alias string) ConfigOption {
	return func(c *config) {
		c.alias = alias
	}
}

func (withConfig) Host(host string) ConfigOption {
	return func(c *config) {
		c.host = host
	}
}

func (withConfig) Port(port uint) ConfigOption {
	return func(c *config) {
		c.port = port
	}
}

func (withConfig) Username(username string) ConfigOption {
	return func(c *config) {
		c.username = username
	}
}

func (withConfig) Password(password string) ConfigOption {
	return func(c *config) {
		c.password = password
	}
}

func (withConfig) DB(db uint) ConfigOption {
	return func(c *config) {
		c.db = db
	}
}

func (withConfig) DialTimeout(timeout time.Duration) ConfigOption {
	return func(c *config) {
		c.dialTimeout = timeout
	}
}

func (withConfig) ReadTimeout(timeout time.Duration) ConfigOption {
	return func(c *config) {
		c.readTimeout = timeout
	}
}

func (withConfig) WriteTimeout(timeout time.Duration) ConfigOption {
	return func(c *config) {
		c.writeTimeout = timeout
	}
}

func (withConfig) PoolTimeout(timeout time.Duration) ConfigOption {
	return func(c *config) {
		c.poolTimeout = timeout
	}
}

func (withConfig) MaxConnAge(age time.Duration) ConfigOption {
	return func(c *config) {
		c.maxConnAge = age
	}
}

func (withConfig) IdleTimeout(timeout time.Duration) ConfigOption {
	return func(c *config) {
		c.idleTimeout = timeout
	}
}
func (withConfig) PoolSize(size uint) ConfigOption {
	return func(c *config) {
		c.poolSize = size
	}
}

func (withConfig) MinIdleConn(count uint) ConfigOption {
	return func(c *config) {
		c.minIdleConn = count
	}
}

//----------------------------------------------------------------------------------------------------------------------

type lockOptions struct {
	value           string
	blockingTimeout *time.Duration
	backoff         IBackoff
}

func newLockOptions(options ...LockOption) *lockOptions {
	c := &lockOptions{
		backoff: &ZeroBackoff{},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

type LockOption func(c *lockOptions)

type withLock struct{}

func (withLock) Value(value string) LockOption {
	return func(c *lockOptions) {
		c.value = value
	}
}

func (withLock) Backoff(backoff IBackoff) LockOption {
	return func(c *lockOptions) {
		c.backoff = backoff
	}
}

func (withLock) BlockingTimeout(timeout time.Duration) LockOption {
	return func(c *lockOptions) {
		c.blockingTimeout = &timeout
	}
}

//----------------------------------------------------------------------------------------------------------------------

type delayQueueOptions struct {
	bucketCount      uint64
	bucketPrefix     string
	scheduleInterval time.Duration
}

func newDelayQueueOptions(options ...DelayQueueOption) *delayQueueOptions {
	c := &delayQueueOptions{
		bucketCount:  1,
		bucketPrefix: "dq-bucket",
	}
	for _, option := range options {
		option(c)
	}
	return c
}

type DelayQueueOption func(c *delayQueueOptions)

type withDelayQueue struct{}

func (withDelayQueue) BucketCount(count uint64) DelayQueueOption {
	return func(c *delayQueueOptions) {
		c.bucketCount = count
	}
}

func (withDelayQueue) BucketPrefix(prefix string) DelayQueueOption {
	return func(c *delayQueueOptions) {
		c.bucketPrefix = prefix
	}
}

func (withConfig) ScheduleInterval(interval time.Duration) DelayQueueOption {
	return func(c *delayQueueOptions) {
		c.scheduleInterval = interval
	}
}
