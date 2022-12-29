package xgorm

import "gorm.io/gorm"

type Config struct {
	Alias           string  `json:"alias"`
	Drive           string  `json:"drive"`
	Host            string  `json:"host"`
	Port            uint    `json:"port"`
	DB              string  `json:"db"`
	User            string  `json:"user"`
	Password        string  `json:"password"`
	Charset         string  `json:"charset"`
	Loc             *string `json:"loc"`
	ParseTime       *bool   `json:"parse_time"`
	MaxIdleConn     *uint   `json:"max_idle_conn"`
	MaxOpenConn     *uint   `json:"max_open_conn"`
	ConnMaxLifetime *uint   `json:"conn_max_lifetime"`
	ConnMaxIdleTime *uint   `json:"conn_max_idle_time"`
	Debug           bool    `json:"debug"`
	GormConfig      *gorm.Config
}

func (c *Config) init() {
	if c.Drive == "" {
		c.Drive = "mysql"
	}
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}
	if c.Port == 0 {
		c.Port = 3306
	}
	if c.User == "" {
		c.User = "root"
	}
	if c.Charset == "" {
		c.Charset = "utf8mb4"
	}
}
