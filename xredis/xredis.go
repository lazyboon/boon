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

func (c *Client) AcquireLock(ctx context.Context, key string, expiration time.Duration, options ...LockOption) (*Lock, error) {
	return NewLock(ctx, c.Client, key, expiration, options...)
}

func InitWithConfigs(configs []*Config) {
	for _, cfg := range configs {
		AddConnectPool(cfg.ToOptions()...)
	}
}

func AddConnectPool(options ...ConfigOption) {
	lock.Lock()
	defer lock.Unlock()

	initInstancesContainer()
	conf := newConfig(options...)
	key := connectPoolKey(conf)

	// check connection pool for conflicts, the same host、port、db can only appear once
	if _, ok := connectPoolSet[key]; ok {
		panic("conflict when building connection pool container")
	}
	connectPoolSet[key] = struct{}{}

	client := redis.NewClient(&redis.Options{
		Network:         "tcp",
		Addr:            fmt.Sprintf("%s:%d", conf.host, conf.port),
		Username:        conf.username,
		Password:        conf.password,
		DB:              int(conf.db),
		DialTimeout:     conf.dialTimeout,
		ReadTimeout:     conf.readTimeout,
		WriteTimeout:    conf.writeTimeout,
		PoolSize:        int(conf.poolSize),
		PoolTimeout:     conf.poolTimeout,
		MinIdleConns:    int(conf.minIdleConn),
		MaxIdleConns:    int(conf.maxIdleConn),
		ConnMaxIdleTime: conf.connMaxIdleTime,
		ConnMaxLifetime: conf.connMaxLifetime,
	})
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
	instanceMap[conf.alias] = &Client{Client: client}
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

func connectPoolKey(c *config) string {
	return fmt.Sprintf("%s+%d+%d", c.host, c.port, c.db)
}
