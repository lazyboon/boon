package response

import (
	"net/http"
	"strings"
)

type Response struct {
	StatusCode int               `json:"status_code"`
	Header     map[string]string `json:"header"`
	Code       int               `json:"code"`
	Msg        string            `json:"msg"`
	Data       interface{}       `json:"data"`
	Errs       []error           `json:"errs"`
}

type Handler func() *Response

func New() Handler {
	return func() *Response {
		return &Response{
			Errs:   make([]error, 0),
			Header: make(map[string]string),
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
