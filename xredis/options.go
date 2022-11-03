package xredis

import (
	"github.com/lazyboon/boon/retry"
	"time"
)

type Options struct {
	Config withConfig
}

var (
	WithConfig = withConfig{}
	WithLock   = withLock{}
	WithDelay  = withDelay{}
)

//----------------------------------------------------------------------------------------------------------------------

type config struct {
	alias           string
	host            string
	port            uint
	username        string
	password        string
	db              uint
	dialTimeout     time.Duration
	readTimeout     time.Duration
	writeTimeout    time.Duration
	poolSize        uint
	poolTimeout     time.Duration
	minIdleConn     uint
	maxIdleConn     uint
	connMaxIdleTime time.Duration
	connMaxLifetime time.Duration
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

func (withConfig) MaxIdleConn(count uint) ConfigOption {
	return func(c *config) {
		c.maxIdleConn = count
	}
}

func (withConfig) ConnMaxIdleTime(duration time.Duration) ConfigOption {
	return func(c *config) {
		c.connMaxIdleTime = duration
	}
}

func (withConfig) ConnMaxLifetime(duration time.Duration) ConfigOption {
	return func(c *config) {
		c.connMaxLifetime = duration
	}
}

//----------------------------------------------------------------------------------------------------------------------

type lockOptions struct {
	value           string
	blockingTimeout *time.Duration
	backoff         retry.IBackoff
}

func newLockOptions(options ...LockOption) *lockOptions {
	c := &lockOptions{
		backoff: &retry.ZeroBackoff{},
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

func (withLock) Backoff(backoff retry.IBackoff) LockOption {
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

type delayOptions struct {
	namespace     string
	listenTopics  []string
	errorCallback func(err error)
}

func newDelayOptions(options ...DelayOption) *delayOptions {
	c := &delayOptions{
		namespace:     "",
		listenTopics:  []string{},
		errorCallback: nil,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

type DelayOption func(c *delayOptions)

type withDelay struct{}

func (withDelay) Namespace(namespace string) DelayOption {
	return func(c *delayOptions) {
		c.namespace = namespace
	}
}

func (withDelay) ListenTopics(topics ...string) DelayOption {
	return func(c *delayOptions) {
		c.listenTopics = topics
	}
}

func (withDelay) ErrorCallback(callback func(err error)) DelayOption {
	return func(c *delayOptions) {
		c.errorCallback = callback
	}
}
