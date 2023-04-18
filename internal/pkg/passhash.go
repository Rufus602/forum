package pkg

import "golang.org/x/crypto/bcrypt"

func PassSecurity(s string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(s), 10)
	if err != nil {
		return "", err
	}
	return string(pass), err
}
