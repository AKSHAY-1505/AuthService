package controllers

import "github.com/gin-gonic/gin"

func Login(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.ShouldBindBodyWithJSON(&requestBody)
}
