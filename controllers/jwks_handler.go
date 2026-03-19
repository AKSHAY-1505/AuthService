package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/AKSHAY-1505/auth-service/auth"
	"github.com/gin-gonic/gin"
)

func JWKSHandler(c *gin.Context) {
	set := auth.GetJWKS()
	jsonbuf, err := json.Marshal(set)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to serialize jwks",
		})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonbuf)
}
