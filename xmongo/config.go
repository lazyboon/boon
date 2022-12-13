package xmongo

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

func (c *Config) init() {
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}
	if c.Port == 0 {
		c.Port = 27017
	}
}
