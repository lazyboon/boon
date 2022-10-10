package access

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type bodyWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (b bodyWriter) Write(bs []byte) (int, error) {
	b.Body.Write(bs)
	return b.ResponseWriter.Write(bs)
}
