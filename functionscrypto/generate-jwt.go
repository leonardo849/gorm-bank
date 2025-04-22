package functionscrypto

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(customerID uint, role int, roleUpdatedAt time.Time, bankAccountID uint) (string, error) {
	claims := jwt.MapClaims{
		"ID":   customerID,
		"bank-account_id": bankAccountID, 
		"role": role,
		"role_updated_at": roleUpdatedAt.Unix(),
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, err
}
