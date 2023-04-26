package handlers

import "testForum/pkg/models"

type TemplateStructure struct {
	User     *models.User
	Post     *models.Post
	Posts    []*models.Post
	Comments []*models.Comment
}
