package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Model struct {
	DB *sql.DB
}

/*############################################################################################################*/
func (m *Model) GetUser() {

}
func (m *Model) GetPost(postId int) (*Post, error) {
	row := m.DB.QueryRow(`SELECT post_id, user_name, title, text, category FROM Posts where post_id=?`, postId)
	post := &Post{}
	err := row.Scan(&post.PostId, &post.UserName, &post.Title, &post.Text, &post.Category)
	if err != nil {
		return nil, err
	}
	post.Likes, post.Dislikes, err = m.GetReactionCountPost(post.PostId)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (m *Model) GetPostAll() ([]*Post, error) {
	query := `SELECT post_id, user_name, title, text, category FROM Posts`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []*Post{}
	err = m.PostShorter(rows, posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
func (m *Model) GetPostCategories(category Category) ([]*Post, error) {
	query := `SELECT post_id, user_name, title, text, category FROM Posts where category=?`
	rows, err := m.DB.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []*Post{}
	err = m.PostShorter(rows, posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *Model) GetPostCreated(userId int) ([]*Post, error) {
	query := `SELECT post_id, user_name, title, text, category FROM Posts where user_id=?`
	rows, err := m.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []*Post{}
	err = m.PostShorter(rows, posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *Model) GetPostLiked(userId int) ([]*Post, error) {
	query := `SELECT p.post_id, p.user_name, p.title, p.text, p.category FROM posts p join PostReactions l on p.post_id=l.post_id where l.user_id=? and l.reaction=1`
	rows, err := m.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []*Post{}
	err = m.PostShorter(rows, posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *Model) GetComments(postId int) ([]*Comment, error) {

	query := `SELECT comment_id, user_id, post_id, user_name, text FROM Comments where post_id=?`
	rows, err := m.DB.Query(query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*Comment{}
	for rows.Next() {
		comment := &Comment{}
		err = rows.Scan(&comment.CommentId, &comment.UserId, &comment.PostId, &comment.UserName, &comment.Text)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNoRecord
			}
			return nil, err
		}
		comment.Likes, comment.Dislikes, err = m.GetReactionCountComment(comment.CommentId)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

/*############################################################################################################*/

func (m *Model) ReactPost(userId int, postId int, reaction int) error {
	result, err := m.GetUserPostReaction(userId, postId)
	if err != nil {
		return err
	}
	if result == reaction {
		_, err := m.DB.Exec("DELETE FROM PostReactions WHERE post_id = $1 AND user_id = $2", postId, userId)
		if err != nil {
			return err
		}
	} else if result == 0 {
		query := `INSERT INTO PostReactions ( user_id, post_id, reaction) 
			VALUES(?,?,?)`

		_, err := m.DB.Exec(query, userId, postId, reaction)
		if err != nil {
			return err
		}

	} else {
		_, err = m.DB.Exec("UPDATE PostReactions SET reaction = $1, WHERE user_id = $2 and post_id=$3", reaction, userId, postId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) ReactComment(userId int, commentId int, reaction int) error {
	result, err := m.GetUserPostReaction(userId, commentId)
	if err != nil {
		return err
	}
	if result == reaction {
		_, err := m.DB.Exec("DELETE FROM CommentReactions WHERE post_id = $1 AND comment_id = $2", commentId, userId)
		if err != nil {
			return err
		}
	} else if result == 0 {
		query := `INSERT INTO CommentReactions ( user_id, comment_id, reaction) 
			VALUES(?,?,?)`

		_, err := m.DB.Exec(query, userId, commentId, reaction)
		if err != nil {
			return err
		}

	} else {
		_, err = m.DB.Exec("UPDATE CommentReactions SET reaction = $1, WHERE user_id = $2 and comment_id=$3", reaction, userId, commentId)
		if err != nil {
			return err
		}
	}
	return nil
}

/*############################################################################################################*/
func (m *Model) InsertComment(postId int, userId int, userName string, text string) error {
	query := `INSERT INTO Comments ( user_id, post_id, user_name, text) 
			VALUES(?,?,?,?)`

	_, err := m.DB.Exec(query, userId, postId, userName, text)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) InsertPost(post Post) error {
	query := `INSERT INTO Posts(post_id, user_id, user_name, title, text, category) 
			VALUES(?,?,?,?,?,?,?)`
	_, err := m.DB.Exec(query, post.UserId, post.UserName, post.Title, post.Text, post.Category)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) InsertUser(u User) error {
	query := `INSERT INTO Users ( user_name, gmail, password) 
			VALUES(?,?,?)`

	_, err := m.DB.Exec(query, u.UserName, u.Gmail, u.Password)
	if err != nil {
		return err
	}
	return nil
}

type ErrorMsg struct {
	Status int
	Msg    string
}

/*######################################################Acquisition#######################################################################################*/

func (m *Model) GetReactionCountPost(postId int) (likes, dislikes int, err error) {
	rowLikes := m.DB.QueryRow("SELECT COUNT(*) from PostReactions where post_id=$1 and reaction=1", postId)
	rowDislikes := m.DB.QueryRow("SELECT COUNT(*) from PostReactions where post_id=$1 and reaction=2", postId)

	err = rowLikes.Scan(&likes)
	if err != nil {
		return 0, 0, err
	}
	err = rowDislikes.Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}
	return likes, dislikes, nil
}
func (m *Model) GetReactionCountComment(postId int) (likes, dislikes int, err error) {
	rowLikes := m.DB.QueryRow("SELECT COUNT(*) from CommentReactions where comment_id=$1 and reaction=1", postId)
	rowDislikes := m.DB.QueryRow("SELECT COUNT(*) from CommentReactions where comment_id=$1 and reaction=2", postId)

	err = rowLikes.Scan(&likes)
	if err != nil {
		return 0, 0, err
	}
	err = rowDislikes.Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}
	return likes, dislikes, nil
}

func (m *Model) PostShorter(rows *sql.Rows, posts []*Post) (err error) {
	for rows.Next() {
		post := &Post{}
		err = rows.Scan(&post.PostId, &post.UserName, &post.Title, &post.Text, &post.Category)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNoRecord
			}
			return err
		}
		post.Likes, post.Dislikes, err = m.GetReactionCountPost(post.PostId)
		if err != nil {
			return err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func (m *Model) GetUserPostReaction(userId int, postId int) (result int, err error) {
	row := m.DB.QueryRow(`select reaction from PostReactions where user_id=? and post_id=?`, userId, postId)
	err = row.Scan(&result)
	if err != nil {
		fmt.Println("error is right here")
		return 0, err
	}
	return result, nil
}

func (m *Model) GetUserCommentReaction(userId int, comment int) (result int, err error) {
	row := m.DB.QueryRow(`select reaction from CommentReactions where user_id=? and comment_id=?`, userId, comment)
	err = row.Scan(&result)
	if err != nil {
		fmt.Println("error is right here")
		return 0, err
	}
	return result, nil
}

/*#####################################################################################################################################################*/
