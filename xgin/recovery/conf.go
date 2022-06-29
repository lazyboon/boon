package recovery

import "github.com/gin-gonic/gin"

type Conf struct {
	LogCallback func(msg string)
	Handle      func(c *gin.Context, err interface{})
}
