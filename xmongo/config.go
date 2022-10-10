package xmongo

import "time"

type Config struct {
	Alias           string  `json:"alias"`
	Host            string  `json:"host"`
	Port            uint    `json:"port"`
	MaxPoolSize     *uint64 `json:"max_pool_size"`
	MinPoolSize     *uint64 `json:"min_pool_size"`
	MaxConnIdleTime *uint   `json:"max_conn_idle_time"`
	Username        *string `json:"username"`
	Password        *string `json:"password"`
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
	if c.MaxPoolSize != nil {
		ans = append(ans, WithConfig.MaxPoolSize(*c.MaxPoolSize))
	}
	if c.MinPoolSize != nil {
		ans = append(ans, WithConfig.MinPoolSize(*c.MinPoolSize))
	}
	if c.MaxConnIdleTime != nil {
		ans = append(ans, WithConfig.MaxConnIdleTime(time.Duration(*c.MaxConnIdleTime)*time.Second))
	}
	if c.Username != nil {
		ans = append(ans, WithConfig.Username(*c.Username))
	}
	if c.Password != nil {
		ans = append(ans, WithConfig.Password(*c.Password))
	}
	return ans
}
