package repositories

import (
	"malakashuttle/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByID(id uint) (*entities.User, error)
	FindAll() ([]entities.User, error)
	GetAllWithPagination(page, limit int) ([]entities.User, int64, error)
	Update(user *entities.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*entities.User, error) {
	var user entities.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]entities.User, error) {
	var users []entities.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) Update(user *entities.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entities.User{}, id).Error
}

func (r *userRepository) GetAllWithPagination(page, limit int) ([]entities.User, int64, error) {
	var users []entities.User
	var total int64

	// Count total
	if err := r.db.Model(&entities.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	offset := (page - 1) * limit
	err := r.db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
