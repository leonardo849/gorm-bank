package functionscrypto

import "golang.org/x/crypto/bcrypt"

func CompareHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}