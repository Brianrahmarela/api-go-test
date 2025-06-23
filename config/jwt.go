package config

import (
	"os"
	"time"
)

// ambil JWT Secret key
func GetJwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

// ambil durasi expired dari JWT
func GetJwtExpirationDuration() time.Duration {
	duration, err := time.ParseDuration(os.Getenv("JWT_EXPIRES_IN"))
	//set default 24 jam jika JWT_EXPIRES_IN == nil
	if err != nil {
		return time.Hour * 24
	}
	return duration
}
