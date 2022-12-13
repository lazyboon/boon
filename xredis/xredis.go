package xredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"sync"
	"time"
)

var (
	lock           sync.RWMutex
	instanceMap    map[string]*Client
	connectPoolSet map[string]struct{}
)

type Client struct {
	*redis.Client
}

func (c *Client) AcquireLock(ctx context.Context, key string, expiration time.Duration, options ...*LockOption) (*Lock, error) {
	return NewLock(ctx, c.Client, key, expiration, options...)
}

func InitWithConfigs(configs []*Config) {
	for _, conf := range configs {
		AddConnectPool(conf)
	}
}

func AddConnectPool(conf *Config) {
	lock.Lock()
	defer lock.Unlock()

	initInstancesContainer()
	conf.init()

	key := connectPoolKey(conf)

	// check connection pool for conflicts, the same host、port、db can only appear once
	if _, ok := connectPoolSet[key]; ok {
		panic("conflict when building connection pool container")
	}
	connectPoolSet[key] = struct{}{}

	options := &redis.Options{
		Network:  "tcp",
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Username: conf.Username,
		Password: conf.Password,
		DB:       int(conf.DB),
	}
	if conf.DialTimeout != nil {
		options.DialTimeout = time.Duration(*conf.DialTimeout) * time.Second
	}
	if conf.ReadTimeout != nil {
		options.ReadTimeout = time.Duration(*conf.ReadTimeout) * time.Second
	}
	if conf.WriteTimeout != nil {
		options.WriteTimeout = time.Duration(*conf.WriteTimeout) * time.Second
	}
	if conf.PoolSize != nil {
		options.PoolSize = int(*conf.PoolSize)
	}
	if conf.PoolTimeout != nil {
		options.PoolTimeout = time.Duration(*conf.PoolTimeout) * time.Second
	}
	if conf.MinIdleConn != nil {
		options.MinIdleConns = int(*conf.MinIdleConn)
	}
	if conf.MaxIdleConn != nil {
		options.MaxIdleConns = int(*conf.MaxIdleConn)
	}
	if conf.ConnMaxIdleTime != nil {
		options.ConnMaxIdleTime = time.Duration(*conf.ConnMaxIdleTime) * time.Second
	}
	if conf.ConnMaxLifetime != nil {
		options.ConnMaxLifetime = time.Duration(*conf.ConnMaxLifetime) * time.Second
	}

	client := redis.NewClient(options)
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
	instanceMap[conf.Alias] = &Client{Client: client}
}

func Connect(alias ...string) *Client {
	k := ""
	if len(alias) > 0 {
		k = alias[len(alias)-1]
	}
	return instanceMap[k]
}

func initInstancesContainer() {
	if instanceMap == nil {
		instanceMap = make(map[string]*Client)
	}
	if connectPoolSet == nil {
		connectPoolSet = make(map[string]struct{})
	}
}

func connectPoolKey(c *Config) string {
	return fmt.Sprintf("%s+%d+%d", c.Host, c.Port, c.DB)
}
