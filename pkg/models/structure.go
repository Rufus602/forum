package models

import "errors"

var ErrNoRecord = errors.New("Create User")

type Post struct {
	PostId   int
	UserId   int
	UserName string
	Title    string
	Text     string
	Category string
	Likes    int
	Dislikes int
	Reaction int
}
type Comment struct {
	CommentId int
	UserId    int
	PostId    int
	UserName  int
	Text      string
	Likes     int
	Dislikes  int
	Reaction  int
}
type User struct {
	UserId   int
	UserName string
	Gmail    string
	Password string
}
type ReactionPost struct {
	UserId   int
	PostId   int
	Reaction int
}
type ReactionComment struct {
	UserId    int
	CommentId int
	Reaction  int
}

type Category int
