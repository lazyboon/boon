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

	// check connection pool for conflicts, the same host、port can only appear once
	if _, ok := connectPoolSet[key]; ok {
		panic("conflict when building connection pool container")
	}
	connectPoolSet[key] = struct{}{}

	// build mongo options
	opts := mongoOptions.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%d", conf.Host, conf.Port))

	if conf.Username != nil && conf.Password != nil {
		credential := mongoOptions.Credential{
			Username: *conf.Username,
			Password: *conf.Password,
		}
		if conf.AuthSource != nil {
			credential.AuthSource = *conf.AuthSource
		}
		opts.SetAuth(credential)
	}
	if conf.MaxPoolSize != nil {
		opts.SetMaxPoolSize(*conf.MaxPoolSize)
	}
	if conf.MinPoolSize != nil {
		opts.SetMinPoolSize(*conf.MinPoolSize)
	}
	if conf.MaxConnIdleTime != nil {
		opts.SetMaxConnIdleTime(time.Duration(*conf.MaxConnIdleTime) * time.Second)
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
	instanceMap[conf.Alias] = client
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

func connectPoolKey(c *Config) string {
	return fmt.Sprintf("%s+%d", c.Host, c.Port)
}
