package pkg

import (
	"errors"
	"fmt"
	"strings"
	"testForum/internal/models"
)

var (
	ErrInvalidEmail    = errors.New("invalid email address")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidUsername = errors.New("invalid username")
)

func CheckUserInfo(user models.User) error {
	// if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(user.Email) {
	// 	fmt.Println("111")
	// 	return ErrInvalidEmail
	// }

	for _, w := range user.User_name {
		if w < 32 || w > 126 {
			fmt.Println("222")
			return ErrInvalidUsername
		}
	}

	if !checkPassword(user.Password) {
		fmt.Println("333")
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
		fmt.Println("444")
		return false
	}

	if !contains(password, numbers) || !contains(password, lowerCase) || !contains(password, upperCase) || !contains(password, symbols) {
		fmt.Println("555")
		return false
	}

	for _, w := range symbols {
		if w < 32 || w > 126 {
			fmt.Println("666")
			return false
		}
	}
	return true
}

func contains(s, checkSymbols string) bool {
	for _, w := range checkSymbols {
		if strings.Contains(s, string(w)) {
			fmt.Println("777")
			return true
		}
	}
	fmt.Println("888")
	return false
}
