package response

import (
	"net/http"
	"strings"
)

type ContentType int8

const (
	ContentTypeJSON = iota
	ContentTypeIndentedJSON
	ContentTypeSecureJSON
	ContentTypeJsonpJSON
	ContentTypeAsciiJSON
	ContentTypePureJSON
	ContentTypeMsgPack
	ContentTypeProtoBuf
	ContentTypeRedirect
	ContentTypeString
	ContentTypeTOML
	ContentTypeXML
	ContentTypeYAML
	ContentTypeHTML
)

type Response struct {
	StatusCode int               `json:"status_code"`
	Header     map[string]string `json:"header"`
	Code       int               `json:"code"`
	Msg        string            `json:"msg"`
	Data       interface{}       `json:"data"`
	Errs       []error           `json:"errs"`
	Type       ContentType       `json:"type"`
	HTMLPath   string            `json:"html_path"`
}

type Handler func() *Response

func New() Handler {
	return func() *Response {
		return &Response{
			Errs:   make([]error, 0),
			Header: make(map[string]string),
			Type:   ContentTypeJSON,
		}
	}
}

func NewWithStatusCode(statusCode int) Handler {
	return func() *Response {
		return &Response{
			StatusCode: statusCode,
			Code:       statusCode,
			Msg:        http.StatusText(statusCode),
			Errs:       make([]error, 0),
			Header:     make(map[string]string),
		}
	}
}

func (h Handler) WithStatusCode(statusCode int) Handler {
	return func() *Response {
		a := h()
		a.StatusCode = statusCode
		return a
	}
}

func (h Handler) WithHeader(key string, val string) Handler {
	return func() *Response {
		a := h()
		a.Header[key] = val
		return a
	}
}

func (h Handler) WithCode(code int) Handler {
	return func() *Response {
		a := h()
		a.Code = code
		return a
	}
}

func (h Handler) WithMsg(msg string) Handler {
	return func() *Response {
		a := h()
		a.Msg = msg
		return a
	}
}

func (h Handler) WithData(data interface{}) Handler {
	return func() *Response {
		a := h()
		a.Data = data
		return a
	}
}

func (h Handler) WithErr(err error) Handler {
	return func() *Response {
		a := h()
		a.Errs = append(a.Errs, err)
		return a
	}
}

func (h Handler) WithHTMLPath(path string) Handler {
	return func() *Response {
		a := h()
		a.HTMLPath = path
		return a
	}
}

func (h Handler) JSON() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeJSON
		return a
	}
}

func (h Handler) IndentedJSON() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeIndentedJSON
		return a
	}
}

func (h Handler) SecureJSON() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeSecureJSON
		return a
	}
}

func (h Handler) JsonpJSON() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeJsonpJSON
		return a
	}
}

func (h Handler) AsciiJSON() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeAsciiJSON
		return a
	}
}

func (h Handler) PureJSON() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypePureJSON
		return a
	}
}

func (h Handler) MsgPack() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeMsgPack
		return a
	}
}

func (h Handler) ProtoBuf() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeProtoBuf
		return a
	}
}

func (h Handler) Redirect() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeRedirect
		return a
	}
}

func (h Handler) String() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeString
		return a
	}
}

func (h Handler) TOML() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeTOML
		return a
	}
}

func (h Handler) XML() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeXML
		return a
	}
}

func (h Handler) YAML() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeYAML
		return a
	}
}

func (h Handler) HTML() Handler {
	return func() *Response {
		a := h()
		a.Type = ContentTypeHTML
		return a
	}
}

func (h Handler) Error() string {
	a := h()
	var builder strings.Builder
	limit := len(a.Errs) - 1
	for idx, err := range a.Errs {
		builder.WriteString(err.Error())
		if idx < limit {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func (h Handler) StatusCode() int {
	return h().StatusCode
}

func (h Handler) Header() map[string]string {
	return h().Header
}

func (h Handler) Code() int {
	return h().Code
}

func (h Handler) Msg() string {
	return h().Msg
}

func (h Handler) Data() interface{} {
	return h().Data
}

func (h Handler) Errs() []error {
	return h().Errs
}
