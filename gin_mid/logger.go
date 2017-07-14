package gin_mid

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Now().Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if method == "HEAD" {
			return
		}

		entry := logrus.WithField("mod", "gin").WithField(
			"latency", latency.String()).WithField(
			"ip", clientIP).WithField("method", method).WithField("path", path)
		if comment != "" {
			entry = entry.WithField("err", comment)
		}
		entry.Infoln(statusCode)
	}
}
