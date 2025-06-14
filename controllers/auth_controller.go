package controllers

import (
	"log"
	"malakashuttle/dto"
	"malakashuttle/services"
	"malakashuttle/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Validation error: %v", err)
		utils.Response.BadRequest(c, "Invalid input data", req)
		return
	}

	// Call service
	response, err := ac.authService.Register(req)
	if err != nil {
		utils.Response.BuildErrorResponse(c, err)
		return
	}

	utils.Response.Created(c, "User registered successfully", response)
}

func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response.HandleValidationError(c, err, req)
		return
	}

	// Call service
	response, err := ac.authService.Login(req)
	if err != nil {
		utils.Response.BuildErrorResponse(c, err)
		return
	}

	utils.Response.OK(c, "Login successful", response)
}
