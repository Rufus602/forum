package mysql

import (
	"database/sql"
	"fmt"
	"testForum/internal/models"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO `
	fmt.Println(stmt)
	return 0, nil
}

func (m *PostModel) Get(id int) ([]*models.Post, error) {
	return nil, nil
}

func (m *PostModel) Latest() ([]*models.Post, error) {
	return nil, nil
}
