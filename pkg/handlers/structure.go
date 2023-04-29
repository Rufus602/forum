package handlers

import "testForum/pkg/models"

type TemplateStructure struct {
	Signed   bool
	User     *models.User
	Post     *models.Post
	Posts    []*models.Post
	Comments []*models.Comment
	Err      string
	Tag      string
}

type ErrMessage struct{}
