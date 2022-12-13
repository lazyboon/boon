package access

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func New(handler func(entity *Entity), options ...*Option) gin.HandlerFunc {
	if handler == nil {
		panic("access handler must not nil")
	}
	conf := mergeOptions(options...)
	skip := sliceToSet(conf.SkipPaths)
	return func(ctx *gin.Context) {
		mps := NewMethodPath(ctx.Request.Method, ctx.FullPath()).String()
		if _, ok := skip[mps]; ok {
			return
		}

		writer := &bodyWriter{
			ResponseWriter: ctx.Writer,
			Body:           bytes.NewBufferString(""),
		}
		ctx.Writer = writer
		start := time.Now()

		requestID := ctx.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			ctx.Request.Header.Set("X-Request-ID", requestID)
			ctx.Writer.Header().Set("X-Request-ID", requestID)
		}

		// request
		buildRequestEntity := func(requestHeader bool, requestBody bool) *RequestEntity {
			if !requestHeader && !requestBody {
				return nil
			}
			ans := &RequestEntity{}
			if requestHeader {
				ans.Header = httpHeaderToMap(ctx.Request.Header)
			}
			if requestBody {
				contentType := ctx.ContentType()
				switch contentType {
				case gin.MIMEJSON:
					body := make(map[string]interface{})
					err := ctx.ShouldBindBodyWith(&body, binding.JSON)
					if err == nil {
						ans.Body = body
						if cb, ok := ctx.Get(gin.BodyBytesKey); ok {
							if cbb, ok := cb.([]byte); ok {
								ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(cbb))
							}
						}
					}
				case gin.MIMEMultipartPOSTForm, gin.MIMEPOSTForm:
					ctx.GetPostForm("")
					ans.Body = postFormToMap(ctx.Request.PostForm)
				}
			}
			return ans
		}
		var requestEntity *RequestEntity
		if c, ok := conf.SpecificPath[mps]; ok {
			requestEntity = buildRequestEntity(c.RequestHeader != nil && *c.RequestHeader, c.RequestBody != nil && *c.RequestBody)
		} else {
			requestEntity = buildRequestEntity(conf.RequestHeader != nil && *conf.RequestHeader, conf.RequestBody != nil && *conf.RequestBody)
		}

		ctx.Next()

		// response
		buildResponseEntity := func(responseHeader bool, responseBody bool) *ResponseEntity {
			ans := &ResponseEntity{
				Status: writer.Status(),
			}
			if responseHeader {
				ans.Header = httpHeaderToMap(writer.Header())
			}
			if responseBody {
				ans.Body = writer.Body.String()
			}
			return ans
		}

		var responseEntity *ResponseEntity
		if c, ok := conf.SpecificPath[mps]; ok {
			responseEntity = buildResponseEntity(c.ResponseHeader != nil && *c.ResponseHeader, c.ResponseBody != nil && *c.ResponseBody)
		} else {
			responseEntity = buildResponseEntity(conf.ResponseHeader != nil && *conf.ResponseHeader, conf.ResponseBody != nil && *conf.ResponseBody)
		}

		// entity
		latency := time.Now().Sub(start)
		if latency > time.Minute {
			latency = latency - latency%time.Second
		}
		conf.Handler(&Entity{
			Method:     ctx.Request.Method,
			Path:       ctx.Request.RequestURI,
			RemoteAddr: ctx.Request.RemoteAddr,
			Proto:      ctx.Request.Proto,
			Request:    requestEntity,
			Response:   responseEntity,
			Latency:    fmt.Sprintf("%s", latency),
			RequestID:  requestID,
		})
	}
}

func httpHeaderToMap(header http.Header) map[string]interface{} {
	ans := make(map[string]interface{}, len(header))
	for key, val := range header {
		if len(val) > 0 {
			ans[strings.ToLower(key)] = val[0]
		}
	}
	return ans
}

func sliceToSet(data []string) map[string]struct{} {
	ans := make(map[string]struct{})
	for _, item := range data {
		ans[item] = struct{}{}
	}
	return ans
}

func postFormToMap(form url.Values) map[string]interface{} {
	ans := make(map[string]interface{}, len(form))
	for key, val := range form {
		if len(val) > 0 {
			ans[strings.ToLower(key)] = val[0]
		}
	}
	return ans
}
