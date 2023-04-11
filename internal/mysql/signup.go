package mysql

import (
	"database/sql"
	"fmt"
	"testForum/internal/models"
)

func SignUp(str *models.User) {
	db, err := sql.Open("sqlite3", "./internal/database/data.db")
	if err != nil {
		fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	statement, err := db.Prepare("INSERT INTO users (user_name, email, password) VALUES (?,?,?)")
	if err != nil {
		fmt.Println("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	}

	a, err := statement.Exec(str.User_name, str.Email, str.Password)
	if err != nil {
		fmt.Println(a)
	}
}
