package mysql

import (
	"database/sql"
	"fmt"
	"testForum/internal/models"
)

func SignUp(str *models.User) error {
	err := checkDuplex(str)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "./internal/database/data.db")
	if err != nil {
		fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		return err
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO users (user_name, email, password) VALUES (?,?,?)")
	if err != nil {
		fmt.Println("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
		return err
	}

	_, err = statement.Exec(str.User_name, str.Email, str.Password)
	if err != nil {
		return err
	}
	return nil
}
