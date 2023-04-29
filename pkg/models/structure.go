package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("Sign In")

var SignIn = errors.New("SignUp first")

type Post struct {
	PostId   int
	UserId   int
	UserName string
	Title    string
	Text     string
	Category string
	Likes    int
	Dislikes int
}

type Comment struct {
	CommentId int
	UserId    int
	UserName  string
	PostId    int
	Text      string
	Likes     int
	Dislikes  int
}

type User struct {
	UserId   int
	UserName string
	Gmail    string
	Password string
}

type Session struct {
	UserID         int
	UserName       string
	Token          string
	ExpirationDate time.Time
}

type Signed struct {
	Sign bool
}
