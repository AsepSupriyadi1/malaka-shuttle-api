// filepath: d:\source_code\go_workspace\malakashuttle_api\services\auth_service.go
package services

import (
	"malakashuttle/dto"
	"malakashuttle/models"
	"malakashuttle/repositories"
	"malakashuttle/utils"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, utils.NewBadRequestErrorWithDetails("User with this email already exists", nil, req)
	}

	// Create user model from request
	userModel := &models.UserWithPassword{
		User: models.User{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			PhoneNumber: req.PhoneNumber,
			Email:       req.Email,
			Role:        "user", // Default role
		},
		Password: req.Password,
	}

	// Convert to entity for database operations
	userEntity := models.Handlers.User.ToEntityWithPassword(userModel)

	// Hash password
	if err := userEntity.HashPassword(); err != nil {
		return nil, utils.NewInternalServerError("Failed to hash password", err)
	}

	// Save user to database
	if err := s.userRepo.Create(userEntity); err != nil {
		return nil, utils.NewInternalServerError("Failed to create user", err)
	}

	// Prepare response using model
	savedUserModel := models.Handlers.User.FromEntity(userEntity)
	response := &dto.RegisterResponse{
		Email: savedUserModel.Email,
	}

	return response, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by email
	userEntity, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, utils.NewUnauthorizedError("Invalid email or password", nil)
	}

	// Check password
	if err := userEntity.CheckPassword(req.Password); err != nil {
		return nil, utils.NewUnauthorizedError("Invalid email or password", nil)
	}

	// Generate token
	token, err := utils.GenerateToken(userEntity.ID)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to generate token", err)
	}

	// Prepare response
	response := &dto.LoginResponse{
		Token: token,
	}

	return response, nil
}
