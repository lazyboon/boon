package xamqp

type Args = map[string]interface{}

type Config struct {
	Alias     string            `json:"alias"`
	Host      string            `json:"host"`
	Port      int               `json:"port"`
	User      string            `json:"user"`
	Password  string            `json:"password"`
	Vhost     string            `json:"vhost"`
	Exchanges []*ExchangeConfig `json:"exchanges"`
}

type ExchangeConfig struct {
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Durable    bool           `json:"durable"`
	AutoDelete bool           `json:"auto_delete"`
	Internal   bool           `json:"internal"`
	NoWait     bool           `json:"no_wait"`
	Args       Args           `json:"args"`
	Queues     []*QueueConfig `json:"queues"`
}

type QueueConfig struct {
	Name        string          `json:"name"`
	Durable     bool            `json:"durable"`
	AutoDelete  bool            `json:"auto_delete"`
	Exclusive   bool            `json:"exclusive"`
	NoWait      bool            `json:"no_wait"`
	Args        Args            `json:"args"`
	DLXExchange *ExchangeConfig `json:"dlx_exchange"`
	RoutingKey  string          `json:"routing_key"`
}
