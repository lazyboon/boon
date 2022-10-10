package xgorm

import (
	"time"
)

type Config struct {
	Alias           string `json:"alias"`
	Host            string `json:"host"`
	Port            uint   `json:"port"`
	DB              string `json:"db"`
	User            string `json:"user"`
	Password        string `json:"password"`
	Charset         string `json:"charset"`
	MaxIdleConn     *uint  `json:"max_idle_conn"`
	MaxOpenConn     *uint  `json:"max_open_conn"`
	ConnMaxLifetime *uint  `json:"conn_max_lifetime"`
	ConnMaxIdleTime *uint  `json:"conn_max_idle_time"`
	Debug           bool   `json:"debug"`
}

func (c *Config) ToOptions() []ConfigOption {
	ans := make([]ConfigOption, 0)
	if c.Alias != "" {
		ans = append(ans, WithConfig.Alias(c.Alias))
	}
	if c.Host != "" {
		ans = append(ans, WithConfig.Host(c.Host))
	}
	if c.Port != 0 {
		ans = append(ans, WithConfig.Port(c.Port))
	}
	if c.DB != "" {
		ans = append(ans, WithConfig.DB(c.DB))
	}
	if c.User != "" {
		ans = append(ans, WithConfig.User(c.User))
	}
	if c.Password != "" {
		ans = append(ans, WithConfig.Password(c.Password))
	}
	if c.Charset != "" {
		ans = append(ans, WithConfig.Charset(c.Charset))
	}
	if c.MaxIdleConn != nil {
		ans = append(ans, WithConfig.MaxIdleConn(*c.MaxIdleConn))
	}
	if c.MaxOpenConn != nil {
		ans = append(ans, WithConfig.MaxOpenConn(*c.MaxOpenConn))
	}
	if c.ConnMaxLifetime != nil {
		ans = append(ans, WithConfig.ConnMaxLifetime(time.Duration(*c.ConnMaxLifetime)*time.Second))
	}
	if c.ConnMaxIdleTime != nil {
		ans = append(ans, WithConfig.ConnMaxIdleTime(time.Duration(*c.ConnMaxIdleTime)*time.Second))
	}
	ans = append(ans, WithConfig.Debug(c.Debug))
	return ans
}
