package entities

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `gorm:"size:100;uniqueIndex"`
	Password    string
	Role        string `gorm:"type:enum('user','admin', 'staff');default:'user'"`
	FirstName   string `gorm:"size:50"`
	LastName    string `gorm:"size:50"`
	PhoneNumber string `gorm:"size:20"`
}

// HashPassword hashes the user's password using bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the user's hashed password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
