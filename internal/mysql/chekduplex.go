package mysql

import (
	"database/sql"
	"errors"
	"testForum/internal/models"
)

func checkDuplex(str *models.User) error {
	db, err := sql.Open("sqlite3", "./internal/database/data.db")
	if err != nil {
		return err
	}

	defer db.Close()

	count := 0

	query := "SELECT COUNT(*) FROM users WHERE user_name = ? OR email = ?"

	err = db.QueryRow(query, str.User_name, str.Email).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("User with this username or e-mail already exists")
	}

	return nil
}
