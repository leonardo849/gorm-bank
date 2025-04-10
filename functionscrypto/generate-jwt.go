package functionscrypto

import (
	"banco/utils"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(customerID uint, role utils.Role, roleUpdatedAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"ID":   customerID,
		"role": role,
		"role_updated_at": roleUpdatedAt.Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, err
}
