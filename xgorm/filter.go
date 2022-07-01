package xgorm

import (
	"fmt"
	"strings"
)

type Condition struct {
	Query string
	Args  []interface{}
}

type Filter struct {
	Queries []string
	Args    []interface{}
}

func NewFilter() *Filter {
	return &Filter{
		Queries: make([]string, 0),
		Args:    make([]interface{}, 0),
	}
}

func (f *Filter) Add(query string, args ...interface{}) *Filter {
	f.Queries = append(f.Queries, query)
	f.Args = append(f.Args, args...)
	return f
}

func (f *Filter) QueryAnd() string {
	return fmt.Sprintf("(%s)", strings.Join(f.Queries, " AND "))
}

func (f *Filter) QueryOr() string {
	return fmt.Sprintf("(%s)", strings.Join(f.Queries, " OR "))
}

func (f *Filter) ToAndCondition() *Condition {
	return &Condition{
		Query: f.QueryAnd(),
		Args:  f.Args,
	}
}

func (f *Filter) ToOrCondition() *Condition {
	return &Condition{
		Query: f.QueryOr(),
		Args:  f.Args,
	}
}
