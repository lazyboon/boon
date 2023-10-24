package xgin

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
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

func Wrap(handler func(c *Context) response.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := &Context{Context: ctx}
		h := handler(c)
		c.response(h, func() {
			r := h()
			for key, val := range r.Header {
				c.Context.Header(key, val)
			}
			normal := gin.H{
				"code": r.Code,
				"msg":  r.Msg,
				"data": r.Data,
			}
			switch r.Type {
			case response.ContentTypeJSON:
				c.Context.JSON(r.StatusCode, normal)
			case response.ContentTypeIndentedJSON:
				c.Context.IndentedJSON(r.StatusCode, normal)
			case response.ContentTypeSecureJSON:
				c.Context.SecureJSON(r.StatusCode, normal)
			case response.ContentTypeJsonpJSON:
				c.Context.JSONP(r.StatusCode, normal)
			case response.ContentTypeAsciiJSON:
				c.Context.AsciiJSON(r.StatusCode, normal)
			case response.ContentTypePureJSON:
				c.Context.PureJSON(r.StatusCode, normal)
			case response.ContentTypeProtoBuf:
				c.Context.ProtoBuf(r.StatusCode, normal)
			case response.ContentTypeTOML:
				c.Context.TOML(r.StatusCode, normal)
			case response.ContentTypeXML:
				c.Context.XML(r.StatusCode, normal)
			case response.ContentTypeYAML:
				c.Context.YAML(r.StatusCode, normal)
			case response.ContentTypeMsgPack:
				c.Context.Render(r.StatusCode, render.MsgPack{Data: normal})
			case response.ContentTypeRedirect:
				switch r.Data.(type) {
				case string:
					c.Context.Redirect(r.StatusCode, r.Data.(string))
				}
			case response.ContentTypeString:
				switch r.Data.(type) {
				case string:
					c.Context.String(r.StatusCode, r.Data.(string))
				}
			case response.ContentTypeHTML:
				c.Context.HTML(r.StatusCode, r.HTMLPath, r.Data)
			}
			ctx.Next()
		})
	}
}
