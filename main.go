package main

import (
	"github.com/AKSHAY-1505/auth-service/auth"
	"github.com/AKSHAY-1505/auth-service/handlers"
	"github.com/AKSHAY-1505/auth-service/initializers"
	"github.com/AKSHAY-1505/auth-service/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.InitLogger()
	initializers.InitEnvVariables()
	auth.InitJWTSecretKeys()
	initializers.InitDB()
}

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// auth endpoints
	authEndpoints := r.Group("/auth", middlewares.RequestLoggerMiddleware)
	{
		authEndpoints.POST("/register", handlers.Register)
		authEndpoints.POST("/login", handlers.Login)

		adminEndpoints := authEndpoints.Group("/admins", middlewares.AdminAuthMiddleware)
		{
			adminEndpoints.POST("", handlers.RegisterAdmin)
		}
	}

	r.GET("/.well-known/jwks.json", handlers.JWKSHandler)

	r.Run()
}
