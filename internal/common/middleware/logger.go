package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

// StructuredLog c.Next()前：请求执行前
// c.Next()后：请求执行完成，但还没有结束
func StructuredLog(l *logrus.Entry) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		elapsed := time.Since(t)
		l.WithFields(logrus.Fields{
			"time_elapsed_ms": elapsed.Milliseconds(),
			"request_uri":     c.Request.RequestURI,
			"client_ip":       c.ClientIP(),
		}).Info("request_out")
	}
}
