package main

import (
	"github.com/AKSHAY-1505/auth-service/auth"
	"github.com/AKSHAY-1505/auth-service/controllers"
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
		authEndpoints.POST("/register", controllers.Register)
		authEndpoints.POST("/login", controllers.Login)
		authEndpoints.POST("/register-admin", controllers.RegisterAdmin)
	}

	r.GET("/.well-known/jwks.json", controllers.JWKSHandler)

	r.Run()
}
