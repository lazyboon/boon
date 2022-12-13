package recovery

import "github.com/gin-gonic/gin"

type Option struct {
	LogCallback func(msg string)
	Handler     func(c *gin.Context, err interface{})
}

func NewOption() *Option {
	return &Option{}
}

func (o *Option) SetLogCallback(v func(msg string)) *Option {
	o.LogCallback = v
	return o
}

func (o *Option) SetHandler(v func(c *gin.Context, err interface{})) *Option {
	o.Handler = v
	return o
}

func mergeOptions(options ...*Option) *Option {
	ans := NewOption()
	for _, item := range options {
		if item.LogCallback != nil {
			ans.LogCallback = item.LogCallback
		}
		if item.Handler != nil {
			ans.Handler = item.Handler
		}
	}
	return ans
}
