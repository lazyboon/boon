package recovery

import "github.com/gin-gonic/gin"

var (
	WithConfig = withConfig{}
)

//----------------------------------------------------------------------------------------------------------------------

type config struct {
	logCallback func(msg string)
	handler     func(c *gin.Context, err interface{})
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

func (withConfig) LogCallback(f func(msg string)) ConfigOption {
	return func(c *config) {
		c.logCallback = f
	}
}

func (withConfig) Handler(f func(c *gin.Context, err interface{})) ConfigOption {
	return func(c *config) {
		c.handler = f
	}
}
