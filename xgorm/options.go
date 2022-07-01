package xgorm

import (
	"gorm.io/gorm"
	"time"
)

var (
	WithConfig  = withConfig{}
	WithFindOne = withFindOne{}
	WithSearch  = withSearch{}
	WithCount   = withCount{}
)

//----------------------------------------------------------------------------------------------------------------------

type config struct {
	host            string
	port            uint
	db              string
	user            string
	password        string
	charset         string
	maxIdleConn     *uint
	maxOpenConn     *uint
	connMaxLifetime *time.Duration
	connMaxIdleTime *time.Duration
	debug           bool
	gormConfig      *gorm.Config
}

func newConfig(options ...ConfigOption) *config {
	c := &config{
		host:            "127.0.0.1",
		port:            3306,
		db:              "mysql",
		user:            "root",
		password:        "",
		charset:         "utf8mb4",
		maxIdleConn:     nil,
		maxOpenConn:     nil,
		connMaxLifetime: nil,
		connMaxIdleTime: nil,
		debug:           false,
		gormConfig:      &gorm.Config{},
	}
	for _, option := range options {
		option(c)
	}
	return c

}

type ConfigOption func(c *config)

type withConfig struct{}

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

func (withConfig) DB(db string) ConfigOption {
	return func(c *config) {
		c.db = db
	}
}

func (withConfig) User(user string) ConfigOption {
	return func(c *config) {
		c.user = user
	}
}

func (withConfig) Password(password string) ConfigOption {
	return func(c *config) {
		c.password = password
	}
}

func (withConfig) Charset(charset string) ConfigOption {
	return func(c *config) {
		c.charset = charset
	}
}

func (withConfig) MaxIdleConn(maxIdleConn uint) ConfigOption {
	return func(c *config) {
		c.maxIdleConn = &maxIdleConn
	}
}

func (withConfig) MaxOpenConn(maxOpenConn uint) ConfigOption {
	return func(c *config) {
		c.maxOpenConn = &maxOpenConn
	}
}

func (withConfig) ConnMaxLifetime(d time.Duration) ConfigOption {
	return func(c *config) {
		c.connMaxLifetime = &d
	}
}

func (withConfig) ConnMaxIdleTime(d time.Duration) ConfigOption {
	return func(c *config) {
		c.connMaxIdleTime = &d
	}
}

func (withConfig) Debug(v bool) ConfigOption {
	return func(c *config) {
		c.debug = v
	}
}

func (withConfig) GormConfig(conf *gorm.Config) ConfigOption {
	return func(c *config) {
		c.gormConfig = conf
	}
}

//----------------------------------------------------------------------------------------------------------------------

type findOneConfig struct {
	select_  []string
	unselect []string
	where    *Condition
	having   *Condition
	group    string
	order    string
}

func newFindOneConfig(options ...FindOneOption) *findOneConfig {
	c := &findOneConfig{}
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

func (withFindOne) Where(condition *Condition) FindOneOption {
	return func(c *findOneConfig) {
		c.where = condition
	}
}

func (withFindOne) Having(condition *Condition) FindOneOption {
	return func(c *findOneConfig) {
		c.having = condition
	}
}

func (withFindOne) Group(group string) FindOneOption {
	return func(c *findOneConfig) {
		c.group = group
	}
}

func (withFindOne) Order(order string) FindOneOption {
	return func(c *findOneConfig) {
		c.order = order
	}
}

//----------------------------------------------------------------------------------------------------------------------

type searchConfig struct {
	limit    uint
	page     uint
	select_  []string
	unselect []string
	where    *Condition
	having   *Condition
	group    string
	order    string
	count    bool
}

func newSearchConfig(options ...SearchOption) *searchConfig {
	c := &searchConfig{
		count: true,
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

func (withSearch) Where(condition *Condition) SearchOption {
	return func(c *searchConfig) {
		c.where = condition
	}
}

func (withSearch) Having(condition *Condition) SearchOption {
	return func(c *searchConfig) {
		c.having = condition
	}
}

func (withSearch) Group(group string) SearchOption {
	return func(c *searchConfig) {
		c.group = group
	}
}

func (withSearch) Order(order string) SearchOption {
	return func(c *searchConfig) {
		c.order = order
	}
}

func (withSearch) Count(count bool) SearchOption {
	return func(c *searchConfig) {
		c.count = count
	}
}

//----------------------------------------------------------------------------------------------------------------------

type countConfig struct {
	where  *Condition
	having *Condition
	group  string
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

func (withCount) Where(condition *Condition) CountOption {
	return func(c *countConfig) {
		c.where = condition
	}
}

func (withCount) Having(condition *Condition) CountOption {
	return func(c *countConfig) {
		c.having = condition
	}
}

func (withCount) Group(group string) CountOption {
	return func(c *countConfig) {
		c.group = group
	}
}
