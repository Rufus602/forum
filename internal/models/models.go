package models

// import (
// 	"errors"
// 	"time"
// )

// var ErrNoRecord = errors.New("models: no matching records found")

// type Post struct {
// 	ID      int
// 	Title   string
// 	Content string
// 	Created time.Time
// 	Expires time.Time
// }

type User struct {
	Id        int
	User_name string
	Email     string
	Password  string
}

func NewUser(username, email, pass string) *User {
	return &User{
		User_name: username,
		Email:     email,
		Password:  pass,
	}
}
