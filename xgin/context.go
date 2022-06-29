package xgin

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyboon/boon/response"
	"github.com/lazyboon/boon/xgin/bind"
)

var (
	BeforeResponseCallback func(ctx *gin.Context, handler response.Handler) (stop bool)
	AfterResponseCallback  func(ctx *gin.Context, handler response.Handler)
)

//----------------------------------------------------------------------------------------------------------------------

type Context struct {
	*gin.Context
}

func NewContext(ctx *gin.Context) *Context {
	return &Context{Context: ctx}
}

func (c *Context) JSONReq() interface{} {
	return c.Context.MustGet(string(bind.KeyJSON))
}

func (c *Context) XMLReq() interface{} {
	return c.Context.MustGet(string(bind.KeyXML))
}

func (c *Context) FormReq() interface{} {
	return c.Context.MustGet(string(bind.KeyForm))
}

func (c *Context) QueryReq() interface{} {
	return c.Context.MustGet(string(bind.KeyQuery))
}

func (c *Context) FormPostReq() interface{} {
	return c.Context.MustGet(string(bind.KeyFormPost))
}

func (c *Context) FormMultipartReq() interface{} {
	return c.Context.MustGet(string(bind.KeyFormMultipart))
}

func (c *Context) ProtoBufReq() interface{} {
	return c.Context.MustGet(string(bind.KeyProtoBuf))
}

func (c *Context) MsgPackReq() interface{} {
	return c.Context.MustGet(string(bind.KeyMsgPack))
}

func (c *Context) YAMLReq() interface{} {
	return c.Context.MustGet(string(bind.KeyYAML))
}

func (c *Context) HeaderReq() interface{} {
	return c.Context.MustGet(string(bind.KeyHeader))
}

func (c *Context) TOMLReq() interface{} {
	return c.Context.MustGet(string(bind.KeyTOML))
}

func (c *Context) UriReq() interface{} {
	return c.Context.MustGet(string(bind.KeyUri))
}

func (c *Context) JSON(handler response.Handler) {
	c.response(handler, func() {
		r := handler()
		for key, val := range r.Header {
			c.Context.Header(key, val)
		}
		c.Context.JSON(r.StatusCode, gin.H{
			"code": r.Code,
			"msg":  r.Msg,
			"data": r.Data,
		})
	})
}

func (c *Context) response(handler response.Handler, f func()) {
	if BeforeResponseCallback != nil {
		stop := BeforeResponseCallback(c.Context, handler)
		if stop {
			return
		}
	}
	f()
	if AfterResponseCallback != nil {
		AfterResponseCallback(c.Context, handler)
	}
}

//----------------------------------------------------------------------------------------------------------------------

func Wrap(handler func(c *Context)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := &Context{Context: ctx}
		handler(c)
		ctx.Next()
	}
}
