package metadata

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

func New(options ...*Option) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conf := mergeOptions(options...)
		start := time.Now()
		ctx.Writer = &bodyWriter{
			ResponseWriter: ctx.Writer,
			Body:           bytes.NewBufferString(""),
			ResponseTime:   conf.ResponseTime != nil && *conf.ResponseTime,
		}
		if conf.RequestID != nil && *conf.RequestID {
			requestID := uuid.New().String()
			ctx.Request.Header.Set("X-Request-ID", requestID)
			ctx.Writer.Header().Set("X-Request-ID", requestID)
		}
		if conf.ReceiveTime != nil && *conf.ReceiveTime {
			ctx.Writer.Header().Set("X-Receive-Time", start.Format(time.RFC3339Nano))
		}
		if conf.ServerName != nil && *conf.ServerName != "" {
			ctx.Writer.Header().Set("X-Server-Name", *conf.ServerName)
		}
		if conf.ServerVersion != nil && *conf.ServerVersion != "" {
			ctx.Writer.Header().Set("X-Server-Version", *conf.ServerVersion)
		}
	}
}
