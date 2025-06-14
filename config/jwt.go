package config

import (
	"os"
	"time"
)

func GetJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func GetJWTDuration() time.Duration {
	duration, err := time.ParseDuration(os.Getenv("JWT_EXPIRATION"))
	if err != nil {
		return 24 * time.Hour // default duration time
	}
	return duration
}
