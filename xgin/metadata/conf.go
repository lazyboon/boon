package metadata

var (
	WithConfig = withConfig{}
)

//----------------------------------------------------------------------------------------------------------------------

type config struct {
	requestID     bool
	receiveTime   bool
	responseTime  bool
	serverName    string
	serverVersion string
}

func newConfig(options ...ConfigOption) *config {
	c := &config{}
	for _, option := range options {
		option(c)
	}
	return c
}

//----------------------------------------------------------------------------------------------------------------------

type ConfigOption func(c *config)

type withConfig struct{}

func (withConfig) RequestID(b bool) ConfigOption {
	return func(c *config) {
		c.requestID = b
	}
}

func (withConfig) ReceiveTime(b bool) ConfigOption {
	return func(c *config) {
		c.receiveTime = b
	}
}

func (withConfig) ResponseTime(b bool) ConfigOption {
	return func(c *config) {
		c.responseTime = b
	}
}

func (withConfig) ServerName(name string) ConfigOption {
	return func(c *config) {
		c.serverName = name
	}
}

func (withConfig) ServerVersion(version string) ConfigOption {
	return func(c *config) {
		c.serverVersion = version
	}
}
