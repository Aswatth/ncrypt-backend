package jwt

import (
	"ncrypt/utils/logger"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(token_validity_in_minutes int) (string, error) {
	logger.Log.Println("Generating JWT token")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"is_authorized": true,
		"expiry":        time.Now().Add(time.Duration(token_validity_in_minutes) * time.Minute).Unix(), //Token valid for 20 mins
	})

	return token.SignedString([]byte(os.Getenv("MASTER_PASSWORD_KEY")))
}
