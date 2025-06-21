package services

import (
	"errors"
	"malakashuttle/dto"
	"malakashuttle/repositories"
	"malakashuttle/utils"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(params utils.PaginationParams) (*utils.PaginationResponse, error) {
	users, total, err := s.userRepo.GetAllWithPagination(params.Page, params.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to DTO response format
	data := make([]dto.UserListResponse, len(users))
	for i, user := range users {
		data[i] = *dto.NewUserListResponseFromEntity(&user)
	}

	response := utils.CreatePaginationResponse(data, total, params)
	return &response, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return dto.NewUserResponseFromEntity(user), nil
}

// CreateUser creates a new user (only staff or user role allowed)
func (s *UserService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Validate role (only staff and user allowed)
	if req.Role != "staff" && req.Role != "user" {
		return nil, errors.New("invalid role: only 'staff' and 'user' roles are allowed")
	}

	// Check if user with email already exists
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Convert DTO to entity
	user := req.ToUserEntity()

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return dto.NewUserResponseFromEntity(user), nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Find existing user
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Prevent updating to admin role
	if req.Role == "admin" {
		return nil, errors.New("cannot update user to admin role")
	}

	// Prevent updating admin user
	if user.Role == "admin" {
		return nil, errors.New("cannot update admin user")
	}

	// Check email uniqueness if email is being updated
	if req.Email != "" && req.Email != user.Email {
		_, err := s.userRepo.FindByEmail(req.Email)
		if err == nil {
			return nil, errors.New("user with this email already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// Apply updates from DTO
	req.ApplyToEntity(user)

	// Hash password if it's being updated
	if req.Password != "" {
		if err := user.HashPassword(); err != nil {
			return nil, errors.New("failed to hash password")
		}
	}

	// Update user
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return dto.NewUserResponseFromEntity(user), nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id uint) error {
	// Find user first
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Prevent deleting admin user
	if user.Role == "admin" {
		return errors.New("cannot delete admin user")
	}

	return s.userRepo.Delete(id)
}
