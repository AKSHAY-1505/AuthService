package util

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetLoggerFromContext(c *gin.Context) *logrus.Entry {
	logger, _ := c.Get("logger")

	return logger.(*logrus.Entry)
}
