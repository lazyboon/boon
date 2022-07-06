package xmongo

import "time"

var (
	WithConfig   = withConfig{}
	WithFindByID = withFindByID{}
	WithFindOne  = withFindOne{}
)

//----------------------------------------------------------------------------------------------------------------------

type config struct {
	alias           string
	host            string
	port            uint
	maxPoolSize     *uint64
	minPoolSize     *uint64
	maxConnIdleTime *time.Duration
	username        *string
	password        *string
}

func newConfig(options ...ConfigOption) *config {
	c := &config{
		alias:       "",
		host:        "127.0.0.1",
		port:        27017,
		maxPoolSize: nil,
		minPoolSize: nil,
		username:    nil,
		password:    nil,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

type ConfigOption func(c *config)

type withConfig struct{}

func (withConfig) Alias(alias string) ConfigOption {
	return func(c *config) {
		c.alias = alias
	}
}

func (withConfig) Host(host string) ConfigOption {
	return func(c *config) {
		c.host = host
	}
}

func (withConfig) Port(port uint) ConfigOption {
	return func(c *config) {
		c.port = port
	}
}

func (withConfig) MaxPoolSize(size uint64) ConfigOption {
	return func(c *config) {
		c.maxPoolSize = &size
	}
}

func (withConfig) MinPoolSize(size uint64) ConfigOption {
	return func(c *config) {
		c.minPoolSize = &size
	}
}

func (withConfig) MaxConnIdleTime(t time.Duration) ConfigOption {
	return func(c *config) {
		c.maxConnIdleTime = &t
	}
}

func (withConfig) Username(val string) ConfigOption {
	return func(c *config) {
		c.username = &val
	}
}

func (withConfig) Password(val string) ConfigOption {
	return func(c *config) {
		c.password = &val
	}
}

//----------------------------------------------------------------------------------------------------------------------

type findByIDConfig struct {
	select_  []string
	unselect []string
	unscoped bool
}

func newFindByIDConfig(options ...FindByIDOption) *findByIDConfig {
	c := &findByIDConfig{}
	for _, option := range options {
		option(c)
	}
	return c
}

type FindByIDOption func(c *findByIDConfig)

type withFindByID struct{}

func (withFindByID) Select(fields ...string) FindByIDOption {
	return func(c *findByIDConfig) {
		c.select_ = fields
	}
}

func (withFindByID) Unselect(fields ...string) FindByIDOption {
	return func(c *findByIDConfig) {
		c.unselect = fields
	}
}

func (withFindByID) Unscoped(v bool) FindByIDOption {
	return func(c *findByIDConfig) {
		c.unscoped = v
	}
}

//----------------------------------------------------------------------------------------------------------------------

type findOneConfig struct {
	condition map[string]interface{}
	select_   []string
	unselect  []string
	order     D
	unscoped  bool
}

func newFindOneConfig(options ...FindOneOption) *findOneConfig {
	c := &findOneConfig{
		condition: make(map[string]interface{}),
	}
	for _, option := range options {
		option(c)
	}
	return c
}

type FindOneOption func(c *findOneConfig)

type withFindOne struct{}

func (withFindOne) Select(fields ...string) FindOneOption {
	return func(c *findOneConfig) {
		c.select_ = fields
	}
}

func (withFindOne) Unselect(fields ...string) FindOneOption {
	return func(c *findOneConfig) {
		c.unselect = fields
	}
}

func (withFindOne) Condition(condition map[string]interface{}) FindOneOption {
	return func(c *findOneConfig) {
		c.condition = condition
	}
}

func (withFindOne) Order(order D) FindOneOption {
	return func(c *findOneConfig) {
		c.order = order
	}
}

func (withFindOne) Unscoped(v bool) FindOneOption {
	return func(c *findOneConfig) {
		c.unscoped = v
	}
}

//----------------------------------------------------------------------------------------------------------------------

type searchConfig struct {
	limit     uint
	page      uint
	condition map[string]interface{}
	select_   []string
	unselect  []string
	order     D
	count     bool
	unscoped  bool
}

func newSearchConfig(options ...SearchOption) *searchConfig {
	c := &searchConfig{
		condition: make(map[string]interface{}),
		limit:     10,
		page:      1,
		count:     true,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

type SearchOption func(c *searchConfig)

type withSearch struct{}

func (withSearch) Limit(limit uint) SearchOption {
	return func(c *searchConfig) {
		c.limit = limit
	}
}

func (withSearch) Page(page uint) SearchOption {
	return func(c *searchConfig) {
		c.page = page
	}
}

func (withSearch) Condition(condition map[string]interface{}) SearchOption {
	return func(c *searchConfig) {
		c.condition = condition
	}
}

func (withSearch) Select(fields ...string) SearchOption {
	return func(c *searchConfig) {
		c.select_ = fields
	}
}

func (withSearch) Unselect(fields ...string) SearchOption {
	return func(c *searchConfig) {
		c.unselect = fields
	}
}

func (withSearch) Order(order D) SearchOption {
	return func(c *searchConfig) {
		c.order = order
	}
}

func (withSearch) Count(v bool) SearchOption {
	return func(c *searchConfig) {
		c.count = v
	}
}

func (withSearch) Unscoped(v bool) SearchOption {
	return func(c *searchConfig) {
		c.unscoped = v
	}
}

//----------------------------------------------------------------------------------------------------------------------

type countConfig struct {
	condition map[string]interface{}
	unscoped  bool
}

func newCountConfig(options ...CountOption) *countConfig {
	c := &countConfig{}
	for _, option := range options {
		option(c)
	}
	return c
}

type CountOption func(c *countConfig)

type withCount struct{}

func (withCount) Condition(condition map[string]interface{}) CountOption {
	return func(c *countConfig) {
		c.condition = condition
	}
}

func (withCount) Unscoped(v bool) CountOption {
	return func(c *countConfig) {
		c.unscoped = v
	}
}
