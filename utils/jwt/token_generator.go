package jwt

import (
	"ncrypt/utils/logger"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken() (string, error) {
	logger.Log.Println("Generating JWT token")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"is_authorized": true,
		"expiry":        time.Now().Add(time.Minute * 20).Unix(), //Token valid for 20 mins
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(os.Getenv("MASTER_PASSWORD_KEY")))
}

func ShortLivedToken() (string, error) {
	logger.Log.Println("Generating JWT token")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"is_authorized": true,
		"expiry":        time.Now().Add(time.Minute * 5).Unix(), //Token valid for 5 mins
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(os.Getenv("MASTER_PASSWORD_KEY")))
}
