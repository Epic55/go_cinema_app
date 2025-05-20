package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execute handlers

		status := c.Writer.Status()

		logger.WithFields(logrus.Fields{
			"status": status,
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("request completed")
	}
}
