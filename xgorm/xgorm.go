package xgorm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	lock           sync.RWMutex
	instanceMap    map[string]*gorm.DB
	connectPoolSet map[string]struct{}
)

func InitWithConfigs(configs []*Config, cs ...map[string]*gorm.Config) {
	m := mergeGormConfig(cs...)
	for _, cfg := range configs {
		opts := cfg.ToOptions()
		if c, ok := m[cfg.Alias]; ok {
			opts = append(opts, WithConfig.GormConfig(c))
		}
		AddConnectPool(opts...)
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

	// add connect pool to instance map
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		conf.user,
		conf.password,
		conf.host,
		conf.port,
		conf.db,
		conf.charset,
	)
	db, err := gorm.Open(mysql.Open(dsn), conf.gormConfig)
	if err != nil {
		panic(err)
	}
	if conf.debug {
		db = db.Debug()
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}
	if conf.maxOpenConn != nil {
		sqlDB.SetMaxOpenConns(int(*conf.maxOpenConn))
	}
	if conf.maxIdleConn != nil {
		sqlDB.SetMaxIdleConns(int(*conf.maxIdleConn))
	}
	if conf.connMaxLifetime != nil {
		sqlDB.SetConnMaxLifetime(*conf.connMaxLifetime)
	}
	if conf.connMaxIdleTime != nil {
		sqlDB.SetConnMaxIdleTime(*conf.connMaxIdleTime)
	}
	instanceMap[conf.alias] = db
}

func Connect(alias ...string) *gorm.DB {
	k := ""
	if len(alias) > 0 {
		k = alias[len(alias)-1]
	}
	return instanceMap[k]
}

func initInstancesContainer() {
	if instanceMap == nil {
		instanceMap = make(map[string]*gorm.DB)
	}
	if connectPoolSet == nil {
		connectPoolSet = make(map[string]struct{})
	}
}

func connectPoolKey(c *config) string {
	return fmt.Sprintf("%s+%d+%s", c.host, c.port, c.db)
}

func mergeGormConfig(cs ...map[string]*gorm.Config) map[string]*gorm.Config {
	ans := map[string]*gorm.Config{}
	for _, item := range cs {
		for k, v := range item {
			ans[k] = v
		}
	}
	return ans
}
