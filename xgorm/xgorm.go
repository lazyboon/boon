package xgorm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	lock           sync.RWMutex
	instanceMap    map[string]*gorm.DB
	connectPoolSet map[string]struct{}
)

func InitWithConfigs(configs []*Config) {
	for _, config := range configs {
		AddConnectPool(config)
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

	// add connect pool to instance map
	var dial gorm.Dialector
	switch conf.Drive {
	case "mysql":
		dial = mysql.Open(fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
			conf.User,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.DB,
			conf.Charset,
		))
	case "sqlserver":
		dial = sqlserver.Open(fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?database=%s",
			conf.User,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.DB,
		))
	default:
		panic("add connect pool error: drive unknown")
	}
	db, err := gorm.Open(dial, conf.GormConfig)
	if err != nil {
		panic(err)
	}
	if conf.Debug {
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
	if conf.MaxOpenConn != nil {
		sqlDB.SetMaxOpenConns(int(*conf.MaxOpenConn))
	}
	if conf.MaxIdleConn != nil {
		sqlDB.SetMaxIdleConns(int(*conf.MaxIdleConn))
	}
	if conf.ConnMaxLifetime != nil {
		sqlDB.SetConnMaxLifetime(time.Duration(*conf.ConnMaxLifetime) * time.Second)
	}
	if conf.ConnMaxIdleTime != nil {
		sqlDB.SetConnMaxIdleTime(time.Duration(*conf.ConnMaxIdleTime) * time.Second)
	}
	instanceMap[conf.Alias] = db
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

func connectPoolKey(c *Config) string {
	return fmt.Sprintf("%s+%d+%s", c.Host, c.Port, c.DB)
}
