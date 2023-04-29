package handlers

import "testForum/pkg/models"

type TemplateStructure struct {
	Signed   *models.Signed
	User     *models.User
	Post     *models.Post
	Posts    []*models.Post
	Comments []*models.Comment
	Err      string
	Tag      string
}

type ErrMessage struct{}
