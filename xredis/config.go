package xredis

import "time"

type Config struct {
	Alias        string `json:"alias"`
	Host         string `json:"host"`
	Port         uint   `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DB           uint   `json:"db"`
	DialTimeout  *uint  `json:"dial_timeout"`
	ReadTimeout  *uint  `json:"read_timeout"`
	WriteTimeout *uint  `json:"write_timeout"`
	PoolTimeout  *uint  `json:"pool_timeout"`
	MaxConnAge   *uint  `json:"max_conn_age"`
	IdleTimeout  *uint  `json:"idle_timeout"`
	PoolSize     uint   `json:"pool_size"`
	MinIdleConn  uint   `json:"min_idle_conn"`
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
	if c.Username != "" {
		ans = append(ans, WithConfig.Username(c.Username))
	}
	if c.Password != "" {
		ans = append(ans, WithConfig.Password(c.Password))
	}
	ans = append(ans, WithConfig.DB(c.DB))
	if c.DialTimeout != nil {
		ans = append(ans, WithConfig.DialTimeout(time.Duration(*c.DialTimeout)*time.Second))
	}
	if c.ReadTimeout != nil {
		ans = append(ans, WithConfig.ReadTimeout(time.Duration(*c.ReadTimeout)*time.Second))
	}
	if c.WriteTimeout != nil {
		ans = append(ans, WithConfig.WriteTimeout(time.Duration(*c.WriteTimeout)*time.Second))
	}
	if c.PoolTimeout != nil {
		ans = append(ans, WithConfig.PoolTimeout(time.Duration(*c.PoolTimeout)*time.Second))
	}
	if c.MaxConnAge != nil {
		ans = append(ans, WithConfig.MaxConnAge(time.Duration(*c.MaxConnAge)*time.Second))
	}
	if c.IdleTimeout != nil {
		ans = append(ans, WithConfig.IdleTimeout(time.Duration(*c.IdleTimeout)*time.Second))
	}
	ans = append(ans, WithConfig.PoolSize(c.PoolSize))
	ans = append(ans, WithConfig.MinIdleConn(c.MinIdleConn))
	return ans
}
