package controllers

import (
	"malakashuttle/dto"
	"malakashuttle/services"
	"malakashuttle/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	// Parse pagination params
	params := utils.GetPaginationParams(ctx)
	// Get users from service
	result, err := c.userService.GetAllUsers(params)
	if err != nil {
		utils.Response.InternalServerError(ctx, "Failed to get users", nil)
		return
	}

	utils.Response.OK(ctx, "Users retrieved successfully", result)
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	// Parse user ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Response.BadRequest(ctx, "Invalid user ID", nil)
		return
	}

	// Get user from service
	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			utils.Response.NotFound(ctx, "User not found", nil)
			return
		}
		utils.Response.InternalServerError(ctx, "Failed to get user", nil)
		return
	}

	utils.Response.OK(ctx, "User retrieved successfully", user)
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Response.BadRequest(ctx, "Invalid request data", nil)
		return
	}

	// Create user via service
	user, err := c.userService.CreateUser(req)
	if err != nil {
		utils.Response.BadRequest(ctx, err.Error(), nil)
		return
	}

	utils.Response.Created(ctx, "User created successfully", user)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	// Parse user ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Response.BadRequest(ctx, "Invalid user ID", nil)
		return
	}

	var req dto.UpdateUserRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Response.BadRequest(ctx, "Invalid request data", nil)
		return
	}

	// Update user via service
	user, err := c.userService.UpdateUser(uint(id), req)
	if err != nil {
		if err.Error() == "user not found" {
			utils.Response.NotFound(ctx, "User not found", nil)
			return
		}
		utils.Response.BadRequest(ctx, err.Error(), nil)
		return
	}

	utils.Response.OK(ctx, "User updated successfully", user)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	// Parse user ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Response.BadRequest(ctx, "Invalid user ID", nil)
		return
	}

	// Delete user via service
	err = c.userService.DeleteUser(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			utils.Response.NotFound(ctx, "User not found", nil)
			return
		}
		utils.Response.BadRequest(ctx, err.Error(), nil)
		return
	}

	utils.Response.OK(ctx, "User deleted successfully", nil)
}
