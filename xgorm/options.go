package xgorm

import (
	"gorm.io/gorm"
	"time"
)

var (
	WithConfig   = withConfig{}
	WithFindOne  = withFindOne{}
	WithFindByID = withFindByID{}
	WithSearch   = withSearch{}
	WithCount    = withCount{}
)

//----------------------------------------------------------------------------------------------------------------------

type config struct {
	alias           string
	drive           string
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
		drive:           "mysql",
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

func (withConfig) Alias(alias string) ConfigOption {
	return func(c *config) {
		c.alias = alias
	}
}

func (withConfig) Drive(drive string) ConfigOption {
	return func(c *config) {
		c.drive = drive
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

type findByIDOptions struct {
	select_  []string
	unselect []string
	unscoped bool
}

func newFindByIDOptions(options ...FindByIDOption) *findByIDOptions {
	c := &findByIDOptions{}
	for _, option := range options {
		option(c)
	}
	return c
}

type FindByIDOption func(c *findByIDOptions)

type withFindByID struct{}

func (withFindByID) Select(fields ...string) FindByIDOption {
	return func(c *findByIDOptions) {
		c.select_ = fields
	}
}

func (withFindByID) Unselect(fields ...string) FindByIDOption {
	return func(c *findByIDOptions) {
		c.unselect = fields
	}
}

func (withFindByID) Unscoped(v bool) FindByIDOption {
	return func(c *findByIDOptions) {
		c.unscoped = v
	}
}

//----------------------------------------------------------------------------------------------------------------------

type findOneOptions struct {
	select_  []string
	unselect []string
	where    *Condition
	having   *Condition
	group    string
	order    string
	unscoped bool
}

func newFindOneOptions(options ...FindOneOption) *findOneOptions {
	c := &findOneOptions{}
	for _, option := range options {
		option(c)
	}
	return c
}

type FindOneOption func(c *findOneOptions)

type withFindOne struct{}

func (withFindOne) Select(fields ...string) FindOneOption {
	return func(c *findOneOptions) {
		c.select_ = fields
	}
}

func (withFindOne) Unselect(fields ...string) FindOneOption {
	return func(c *findOneOptions) {
		c.unselect = fields
	}
}

func (withFindOne) Where(condition *Condition) FindOneOption {
	return func(c *findOneOptions) {
		c.where = condition
	}
}

func (withFindOne) Having(condition *Condition) FindOneOption {
	return func(c *findOneOptions) {
		c.having = condition
	}
}

func (withFindOne) Group(group string) FindOneOption {
	return func(c *findOneOptions) {
		c.group = group
	}
}

func (withFindOne) Order(order string) FindOneOption {
	return func(c *findOneOptions) {
		c.order = order
	}
}

func (withFindOne) Unscoped(v bool) FindOneOption {
	return func(c *findOneOptions) {
		c.unscoped = v
	}
}

//----------------------------------------------------------------------------------------------------------------------

type searchOptions struct {
	limit    uint
	page     uint
	select_  []string
	unselect []string
	where    *Condition
	having   *Condition
	group    string
	order    string
	count    bool
	unscoped bool
}

func newSearchOptions(options ...SearchOption) *searchOptions {
	c := &searchOptions{
		limit: 10,
		page:  1,
		count: true,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

type SearchOption func(c *searchOptions)

type withSearch struct{}

func (withSearch) Limit(limit uint) SearchOption {
	return func(c *searchOptions) {
		c.limit = limit
	}
}

func (withSearch) Page(page uint) SearchOption {
	return func(c *searchOptions) {
		c.page = page
	}
}

func (withSearch) Select(fields ...string) SearchOption {
	return func(c *searchOptions) {
		c.select_ = fields
	}
}

func (withSearch) Unselect(fields ...string) SearchOption {
	return func(c *searchOptions) {
		c.unselect = fields
	}
}

func (withSearch) Where(condition *Condition) SearchOption {
	return func(c *searchOptions) {
		c.where = condition
	}
}

func (withSearch) Having(condition *Condition) SearchOption {
	return func(c *searchOptions) {
		c.having = condition
	}
}

func (withSearch) Group(group string) SearchOption {
	return func(c *searchOptions) {
		c.group = group
	}
}

func (withSearch) Order(order string) SearchOption {
	return func(c *searchOptions) {
		c.order = order
	}
}

func (withSearch) Count(count bool) SearchOption {
	return func(c *searchOptions) {
		c.count = count
	}
}

func (withSearch) Unscoped(v bool) SearchOption {
	return func(c *searchOptions) {
		c.unscoped = v
	}
}

//----------------------------------------------------------------------------------------------------------------------

type countOptions struct {
	where    *Condition
	having   *Condition
	group    string
	unscoped bool
}

func newCountOptions(options ...CountOption) *countOptions {
	c := &countOptions{}
	for _, option := range options {
		option(c)
	}
	return c
}

type CountOption func(c *countOptions)

type withCount struct{}

func (withCount) Where(condition *Condition) CountOption {
	return func(c *countOptions) {
		c.where = condition
	}
}

func (withCount) Having(condition *Condition) CountOption {
	return func(c *countOptions) {
		c.having = condition
	}
}

func (withCount) Group(group string) CountOption {
	return func(c *countOptions) {
		c.group = group
	}
}

func (withCount) Unscoped(v bool) CountOption {
	return func(c *countOptions) {
		c.unscoped = v
	}
}
