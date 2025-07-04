package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestController struct{}

func NewTestController() *TestController {
	return &TestController{}
}

// Simple ping test
func (tc *TestController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
