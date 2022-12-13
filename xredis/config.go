package xredis

type Config struct {
	Alias           string `json:"alias"`
	Host            string `json:"host"`
	Port            uint   `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	DB              uint   `json:"db"`
	DialTimeout     *uint  `json:"dial_timeout"`
	ReadTimeout     *uint  `json:"read_timeout"`
	WriteTimeout    *uint  `json:"write_timeout"`
	PoolTimeout     *uint  `json:"pool_timeout"`
	PoolSize        *uint  `json:"pool_size"`
	MinIdleConn     *uint  `json:"min_idle_conn"`
	MaxIdleConn     *uint  `json:"max_idle_conn"`
	ConnMaxIdleTime *uint  `json:"conn_max_idle_time"`
	ConnMaxLifetime *uint  `json:"conn_max_lifetime"`
}

func (c *Config) init() {
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}
	if c.Port == 0 {
		c.Port = 6379
	}
}
