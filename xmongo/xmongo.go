package xmongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	lock        sync.RWMutex
	instanceMap = make(map[string]*mongo.Client)
)

func InitWithConfigs(configs []*Config) {
	for _, conf := range configs {
		AddConnectPool(conf)
	}
}

func AddConnectPool(conf *Config) {
	lock.Lock()
	defer lock.Unlock()

	conf.init()

	// check alias
	if conf.Alias == "" {
		panic("mongo alias required")
	}
	if _, ok := instanceMap[conf.Alias]; ok {
		panic("duplicate mongo alias: " + conf.Alias)
	}

	// build URI
	uri := fmt.Sprintf("mongodb://%s:%d", conf.Host, conf.Port)

	// build options
	opts := mongoOptions.Client().ApplyURI(uri)

	// auth
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

	// pool config
	if conf.MaxPoolSize != nil {
		opts.SetMaxPoolSize(*conf.MaxPoolSize)
	}
	if conf.MinPoolSize != nil {
		opts.SetMinPoolSize(*conf.MinPoolSize)
	}
	if conf.MaxConnIdleTime != nil {
		opts.SetMaxConnIdleTime(time.Duration(*conf.MaxConnIdleTime) * time.Second)
	}
	if conf.ConnectTimeout != nil {
		opts.SetConnectTimeout(time.Duration(*conf.ConnectTimeout) * time.Second)
	}

	// create client
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	// ping
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	instanceMap[conf.Alias] = client
}

func Connect(alias string) *mongo.Client {
	lock.RLock()
	defer lock.RUnlock()

	client, ok := instanceMap[alias]
	if !ok {
		panic("mongo client not found: " + alias)
	}
	return client
}

func CloseAll() {
	lock.Lock()
	defer lock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for k, client := range instanceMap {
		_ = client.Disconnect(ctx)
		delete(instanceMap, k)
	}
}
