package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func DB() (*sql.DB, error) {
	db, err := createDB()
	if err != nil {
		return nil, err
	}
	if err = CreateTables(db); err != nil {
		return nil, err
	}
	return db, err
}

func createDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "Forum.db")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTables(b *sql.DB) error {
	var quesries []string = []string{Users, Post, Comment, PostReaction, CommentReaction, Session}
	for _, each := range quesries {

		_, err := b.Exec(each)
		if err != nil {
			fmt.Println("db has not created")

			return err
		}
	}
	return nil
}

const (
	Users = `CREATE TABLE IF NOT EXISTS Users(
		user_id INTEGER PRIMARY KEY,
		user_name TEXT UNIQUE,
		gmail TEXT UNIQUE,
		password TEXT
		);`
	Post = `CREATE TABLE IF NOT EXISTS Posts (
			post_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER ,
			user_name TEXT,
			title TEXT,
			text TEXT,
			category TEXT,
			FOREIGN KEY (user_id) REFERENCES Users(user_id)
		);`

	Comment = `CREATE TABLE IF NOT EXISTS Comments (
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		user_name TEXT,
		text TEXT ,
		FOREIGN KEY (user_id) REFERENCES Users(user_id),
		FOREIGN KEY (post_id) REFERENCES Posts(post_id)
	);`
	Session = `CREATE TABLE IF NOT EXISTS Sessions (
		session_id INTEGER PRIMARY KEY,
		user_name TEXT,
		user_id INTEGER,
		token TEXT UNIQUE,
		expiration_date TIMESTAMP
	);`
	PostReaction = `CREATE TABLE IF NOT EXISTS PostReactions (
		user_id INTEGER,
		post_id INTEGER,
		reaction INTEGER,
		FOREIGN KEY (user_id) REFERENCES Users(user_id),
		FOREIGN KEY (post_id) REFERENCES Posts(post_id)
	);`

	CommentReaction = `CREATE TABLE IF NOT EXISTS CommentReactions (
		user_id INTEGER,
		comment_id INTEGER,
		reaction INTEGER,
		FOREIGN KEY (user_id) REFERENCES Users(user_id),
		FOREIGN KEY (comment_id) REFERENCES Comments(comment_id)
	);`
)
