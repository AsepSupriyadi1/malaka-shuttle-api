package dto

import (
	"malakashuttle/entities"
	"time"
)

// CreateUserRequest represents request for creating new user
type CreateUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	Role        string `json:"role" binding:"required,oneof=staff user"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number"`
}

// UpdateUserRequest represents request for updating user
type UpdateUserRequest struct {
	Email       string `json:"email" binding:"omitempty,email"`
	Password    string `json:"password" binding:"omitempty,min=8"`
	Role        string `json:"role" binding:"omitempty,oneof=staff user"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

// UserResponse represents user data in response (without password)
type UserResponse struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserListResponse represents user data in list response (simplified)
type UserListResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// NewUserResponseFromEntity creates UserResponse from User entity
func NewUserResponseFromEntity(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Role:        user.Role,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// NewUserListResponseFromEntity creates UserListResponse from User entity
func NewUserListResponseFromEntity(user *entities.User) *UserListResponse {
	return &UserListResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

// ToUserEntity converts CreateUserRequest to User entity
func (r *CreateUserRequest) ToUserEntity() *entities.User {
	return &entities.User{
		Email:       r.Email,
		Password:    r.Password,
		Role:        r.Role,
		FirstName:   r.FirstName,
		LastName:    r.LastName,
		PhoneNumber: r.PhoneNumber,
	}
}

// ApplyToEntity applies UpdateUserRequest to existing User entity
func (r *UpdateUserRequest) ApplyToEntity(user *entities.User) {
	if r.Email != "" {
		user.Email = r.Email
	}
	if r.Password != "" {
		user.Password = r.Password
	}
	if r.Role != "" {
		user.Role = r.Role
	}
	if r.FirstName != "" {
		user.FirstName = r.FirstName
	}
	if r.LastName != "" {
		user.LastName = r.LastName
	}
	if r.PhoneNumber != "" {
		user.PhoneNumber = r.PhoneNumber
	}
}
