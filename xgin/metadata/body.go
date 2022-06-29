package metadata

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"time"
)

type bodyWriter struct {
	gin.ResponseWriter
	Body         *bytes.Buffer
	ResponseTime bool
}

func (b bodyWriter) Write(bs []byte) (int, error) {
	b.Body.Write(bs)
	if b.ResponseTime {
		b.Header().Set("X-Response-Time", time.Now().Format(time.RFC3339Nano))
	}
	return b.ResponseWriter.Write(bs)
}

