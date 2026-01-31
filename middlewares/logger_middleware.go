package middlewares

import (
	"github.com/AKSHAY-1505/auth-service/initializers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestLoggerMiddleware(c *gin.Context) {
	requestId := c.GetHeader("X-REQUEST-ID")
	if requestId == "" {
		requestId = uuid.NewString()
	}

	// create request-scoped logger
	reqLogger := initializers.Log.WithField("request_id", requestId)
	
	// attach to context
	c.Set("logger", reqLogger)

	// add request-id to response headers
	c.Writer.Header().Set("X-REQUEST-ID", requestId)

	// proceed with middleware chain
	c.Next()
}
