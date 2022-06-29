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

func New(conf Conf) gin.HandlerFunc {
	if conf.Handler == nil {
		panic("log callback must not nil")
	}
	return func(ctx *gin.Context) {
		skip := sliceToSet(conf.SkipPaths)
		if _, ok := skip[ctx.Request.URL.Path]; ok {
			return
		}

		writer := &bodyWriter{
			ResponseWriter: ctx.Writer,
			Body:           bytes.NewBufferString(""),
		}
		ctx.Writer = writer
		start := time.Now()

		requestId := ctx.Request.Header.Get("X-Request-Id")
		if requestId == "" {
			requestId = uuid.New().String()
			ctx.Request.Header.Set("X-Request-Id", requestId)
			ctx.Writer.Header().Set("X-Request-Id", requestId)
		}

		// request
		buildRequestEntry := func(requestHeader bool, requestBody bool) *RequestEntry {
			if !requestHeader && !requestBody {
				return nil
			}
			ans := &RequestEntry{}
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
		var requestEntry *RequestEntry
		if c, ok := conf.SpecificPath[ctx.Request.URL.Path]; ok {
			requestEntry = buildRequestEntry(c.RequestHeader, c.RequestBody)
		} else {
			requestEntry = buildRequestEntry(conf.RequestHeader, conf.RequestBody)
		}

		ctx.Next()

		// response
		buildResponseEntry := func(responseHeader bool, responseBody bool) *ResponseEntry {
			ans := &ResponseEntry{
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

		var responseEntry *ResponseEntry
		if c, ok := conf.SpecificPath[ctx.Request.URL.Path]; ok {
			responseEntry = buildResponseEntry(c.ResponseHeader, c.ResponseBody)
		} else {
			responseEntry = buildResponseEntry(conf.ResponseHeader, conf.ResponseBody)
		}

		// entry
		latency := time.Now().Sub(start)
		if latency > time.Minute {
			latency = latency - latency%time.Second
		}
		conf.Handler(&Entry{
			Method:     ctx.Request.Method,
			Path:       ctx.Request.RequestURI,
			RemoteAddr: ctx.Request.RemoteAddr,
			Proto:      ctx.Request.Proto,
			Request:    requestEntry,
			Response:   responseEntry,
			Latency:    fmt.Sprintf("%s", latency),
			RequestId:  requestId,
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
