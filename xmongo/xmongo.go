package xmongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var (
	lock           sync.RWMutex
	instanceMap    map[string]*mongo.Client
	connectPoolSet map[string]struct{}
)

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

	// check connection pool for conflicts, the same host、port can only appear once
	if _, ok := connectPoolSet[key]; ok {
		panic("conflict when building connection pool container")
	}
	connectPoolSet[key] = struct{}{}

	// build mongo options
	opts := mongoOptions.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%d", conf.host, conf.port))

	if conf.username != nil && conf.password != nil {
		opts.SetAuth(mongoOptions.Credential{
			Username: *conf.username,
			Password: *conf.password,
		})
	}
	if conf.maxPoolSize != nil {
		opts.SetMaxPoolSize(*conf.maxPoolSize)
	}
	if conf.minPoolSize != nil {
		opts.SetMinPoolSize(*conf.minPoolSize)
	}
	if conf.maxConnIdleTime != nil {
		opts.SetMaxConnIdleTime(*conf.maxConnIdleTime)
	}

	// create client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	instanceMap[conf.alias] = client
}

func Connect(alias ...string) *mongo.Client {
	k := ""
	if len(alias) > 0 {
		k = alias[len(alias)-1]
	}
	return instanceMap[k]
}

func initInstancesContainer() {
	if instanceMap == nil {
		instanceMap = make(map[string]*mongo.Client)
	}
	if connectPoolSet == nil {
		connectPoolSet = make(map[string]struct{})
	}
}

func connectPoolKey(c *config) string {
	return fmt.Sprintf("%s+%d", c.host, c.port)
}
