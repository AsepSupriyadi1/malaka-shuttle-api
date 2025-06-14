package models

import (
	"time"

	"malakashuttle/entities"

	"gorm.io/gorm"
)

// User model untuk business logic layer
type User struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserWithPassword model untuk internal operations yang memerlukan password
type UserWithPassword struct {
	User
	Password string `json:"password,omitempty"`
}

// UserHandler untuk converter operations
type UserHandler struct{}

// NewUserHandler membuat instance baru dari UserHandler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// FromEntity mengkonversi entity.User ke model.User
func (h *UserHandler) FromEntity(entity *entities.User) *User {
	if entity == nil {
		return nil
	}

	return &User{
		ID:          entity.ID,
		Email:       entity.Email,
		FirstName:   entity.FirstName,
		LastName:    entity.LastName,
		PhoneNumber: entity.PhoneNumber,
		Role:        entity.Role,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// ToEntity mengkonversi model.User ke entity.User
func (h *UserHandler) ToEntity(model *User) *entities.User {
	if model == nil {
		return nil
	}

	return &entities.User{
		Model: gorm.Model{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		Email:       model.Email,
		FirstName:   model.FirstName,
		LastName:    model.LastName,
		PhoneNumber: model.PhoneNumber,
		Role:        model.Role,
	}
}

// FromEntityWithPassword mengkonversi entity.User ke model.UserWithPassword
func (h *UserHandler) FromEntityWithPassword(entity *entities.User) *UserWithPassword {
	if entity == nil {
		return nil
	}

	return &UserWithPassword{
		User: User{
			ID:          entity.ID,
			Email:       entity.Email,
			FirstName:   entity.FirstName,
			LastName:    entity.LastName,
			PhoneNumber: entity.PhoneNumber,
			Role:        entity.Role,
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
		},
		Password: entity.Password,
	}
}

// ToEntityWithPassword mengkonversi model.UserWithPassword ke entity.User
func (h *UserHandler) ToEntityWithPassword(model *UserWithPassword) *entities.User {
	if model == nil {
		return nil
	}

	return &entities.User{
		Model: gorm.Model{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		Email:       model.Email,
		Password:    model.Password,
		FirstName:   model.FirstName,
		LastName:    model.LastName,
		PhoneNumber: model.PhoneNumber,
		Role:        model.Role,
	}
}

// FromEntityList mengkonversi slice entity.User ke slice model.User
func (h *UserHandler) FromEntityList(entities []*entities.User) []*User {
	if entities == nil {
		return nil
	}

	models := make([]*User, len(entities))
	for i, entity := range entities {
		models[i] = h.FromEntity(entity)
	}
	return models
}

// Global handler instance
var UserHandlerInstance = NewUserHandler()
