package service

import (
	"errors"
	"forum/internal/models"
	"regexp"
	"strings"
)

var (
	ErrInvalidEmail    = errors.New("invalid email address")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidUsername = errors.New("invalid username")
)

func checkUserInfo(user models.User) error {
	if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(user.Email) {
		return ErrInvalidEmail
	}

	for _, w := range user.Username {
		if w < 32 || w > 126 {
			return ErrInvalidUsername
		}
	}

	if !checkPassword(user.Password) {
		return ErrInvalidPassword
	}

	return nil
}

func checkPassword(password string) bool {
	numbers := "0123456789"
	lowerCase := "qwertyuiopasdfghjklzxcvbnm"
	upperCase := "QWERTYUIOPASDFGHJKLZXCVBNM"
	symbols := "!@#$%^&*()_-+={[}]|\\:;<,>.?/"

	if len(password) < 8 || len(password) > 20 {
		return false
	}

	if !contains(password, numbers) || !contains(password, lowerCase) || !contains(password, upperCase) || !contains(password, symbols) {
		return false
	}

	for _, w := range symbols {
		if w < 32 || w > 126 {
			return false
		}
	}
	return true
}

func contains(s, checkSymbols string) bool {
	for _, w := range checkSymbols {
		if strings.Contains(s, string(w)) {
			return true
		}
	}
	return false
}
