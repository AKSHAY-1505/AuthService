package middlewares

import (
	"net/http"

	"github.com/AKSHAY-1505/auth-service/auth"
	"github.com/AKSHAY-1505/auth-service/services"
	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware(c *gin.Context) {
	accessToken, err := auth.ExtractAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !services.IsAdmin(accessToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not authorised to perform this action"})
		return
	}

	// proceed with middleware chain
	c.Next()
}
