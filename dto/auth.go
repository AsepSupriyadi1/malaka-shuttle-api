package dto

// RegisterResponse response untuk register
type (
	RegisterRequest struct {
		FirstName   string `json:"first_name" binding:"required"`
		LastName    string `json:"last_name" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=8"`
	}

	LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	RegisterResponse struct {
		Email string `json:"email"`
	}

	// LoginResponse response untuk login
	LoginResponse struct {
		Token string `json:"token"`
	}
)
