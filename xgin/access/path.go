package access

import (
	"fmt"
	"strings"
)

type MethodPath struct {
	method string
	path   string
}

func NewMethodPath(method string, path string) *MethodPath {
	return &MethodPath{
		method: strings.ToUpper(method),
		path:   path,
	}
}

func (m *MethodPath) Method() string {
	return m.method
}

func (m *MethodPath) Path() string {
	return m.path
}

func (m *MethodPath) String() string {
	return fmt.Sprintf("method: %s, path: %s", m.method, m.path)
}
