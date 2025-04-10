package functionscrypto

import "golang.org/x/crypto/bcrypt"

func StringToHash(password string) ([]byte, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}