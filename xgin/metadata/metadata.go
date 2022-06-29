package metadata

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

func New(conf Conf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Writer = &bodyWriter{
			ResponseWriter: ctx.Writer,
			Body:           bytes.NewBufferString(""),
			ResponseTime:   conf.ResponseTime,
		}
		if conf.RequestId {
			requestId := uuid.New().String()
			ctx.Request.Header.Set("X-Request-Id", requestId)
			ctx.Writer.Header().Set("X-Request-Id", requestId)
		}
		if conf.ReceiveTime {
			ctx.Writer.Header().Set("X-Receive-Time", start.Format(time.RFC3339Nano))
		}
		if conf.ServerName != "" {
			ctx.Writer.Header().Set("X-Server-Name", conf.ServerName)
		}
		if conf.ServerVersion != "" {
			ctx.Writer.Header().Set("X-Server-Version", conf.ServerVersion)
		}
	}
}
