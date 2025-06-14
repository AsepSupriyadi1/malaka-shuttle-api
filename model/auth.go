package model

type (
	Login struct {
		GrantType string `json:"grant_type"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	RegisterResponse struct {
		Email string `json:"email"`
	}

	LoginResponse struct {
		AccessToken string `json:"access_token"`
	}
)
