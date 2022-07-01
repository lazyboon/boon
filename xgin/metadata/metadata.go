package metadata

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

func New(options ...ConfigOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conf := newConfig(options...)
		start := time.Now()
		ctx.Writer = &bodyWriter{
			ResponseWriter: ctx.Writer,
			Body:           bytes.NewBufferString(""),
			ResponseTime:   conf.responseTime,
		}
		if conf.requestId {
			requestId := uuid.New().String()
			ctx.Request.Header.Set("X-Request-Id", requestId)
			ctx.Writer.Header().Set("X-Request-Id", requestId)
		}
		if conf.receiveTime {
			ctx.Writer.Header().Set("X-Receive-Time", start.Format(time.RFC3339Nano))
		}
		if conf.serverName != "" {
			ctx.Writer.Header().Set("X-Server-Name", conf.serverName)
		}
		if conf.serverVersion != "" {
			ctx.Writer.Header().Set("X-Server-Version", conf.serverVersion)
		}
	}
}
