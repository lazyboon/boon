package access

var (
	WithConfig     = withConfig{}
	WithBaseConfig = withBaseConfig{}
)

//----------------------------------------------------------------------------------------------------------------------

type baseConfig struct {
	requestHeader  bool
	requestBody    bool
	responseHeader bool
	responseBody   bool
}

func newBaseConfig(options ...BaseConfigOption) *baseConfig {
	c := &baseConfig{}
	for _, option := range options {
		option(c)
	}
	return c
}

//----------------------------------------------------------------------------------------------------------------------

type BaseConfigOption func(c *baseConfig)

type withBaseConfig struct{}

func (withBaseConfig) RequestHeader(v bool) BaseConfigOption {
	return func(c *baseConfig) {
		c.requestHeader = v
	}
}

func (withBaseConfig) RequestBody(v bool) BaseConfigOption {
	return func(c *baseConfig) {
		c.requestBody = v
	}
}

func (withBaseConfig) ResponseHeader(v bool) BaseConfigOption {
	return func(c *baseConfig) {
		c.responseHeader = v
	}
}

func (withBaseConfig) ResponseBody(v bool) BaseConfigOption {
	return func(c *baseConfig) {
		c.responseBody = v
	}
}

//----------------------------------------------------------------------------------------------------------------------

type config struct {
	*baseConfig
	skipPaths    []string
	specificPath map[string]*baseConfig
	handler      func(entity *Entity)
}

func newConfig(handler func(entity *Entity), options ...ConfigOption) *config {
	c := &config{
		baseConfig:   &baseConfig{},
		skipPaths:    make([]string, 0),
		specificPath: make(map[string]*baseConfig),
		handler:      handler,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

//----------------------------------------------------------------------------------------------------------------------

type ConfigOption func(c *config)

type withConfig struct{}

func (withConfig) SkipPaths(paths ...*MethodPath) ConfigOption {
	return func(c *config) {
		if c.skipPaths == nil {
			c.skipPaths = make([]string, 0)
		}
		p := make([]string, 0, len(paths))
		for _, path := range paths {
			p = append(p, path.String())
		}
		c.skipPaths = append(c.skipPaths, p...)
	}
}

func (withConfig) SpecificPath(path *MethodPath, options ...BaseConfigOption) ConfigOption {
	return func(c *config) {
		if c.specificPath == nil {
			c.specificPath = make(map[string]*baseConfig)
		}
		c.specificPath[path.String()] = newBaseConfig(options...)
	}
}

func (withConfig) RequestHeader(v bool) ConfigOption {
	return func(c *config) {
		if c.baseConfig == nil {
			c.baseConfig = &baseConfig{}
		}
		c.requestHeader = v
	}
}

func (withConfig) RequestBody(v bool) ConfigOption {
	return func(c *config) {
		if c.baseConfig == nil {
			c.baseConfig = &baseConfig{}
		}
		c.requestBody = v
	}
}

func (withConfig) ResponseHeader(v bool) ConfigOption {
	return func(c *config) {
		if c.baseConfig == nil {
			c.baseConfig = &baseConfig{}
		}
		c.responseHeader = v
	}
}

func (withConfig) ResponseBody(v bool) ConfigOption {
	return func(c *config) {
		if c.baseConfig == nil {
			c.baseConfig = &baseConfig{}
		}
		c.responseBody = v
	}
}
