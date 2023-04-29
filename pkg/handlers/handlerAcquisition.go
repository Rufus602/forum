package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testForum/pkg/models"
	"text/template"
	"time"
)

var errorMessage = "There is no such user. Maybe incorrect username or password, or you did not register"

func (app *Application) redirect(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodGet
	http.Redirect(w, r, "/signin", http.StatusPermanentRedirect)
	if _, err := w.Write([]byte("loginFirst")); err != nil {
		app.serverError(w, err)
	}
	return
}

func (app *Application) checkerSession(w http.ResponseWriter, r *http.Request) (*models.Session, error) {
	token, err := r.Cookie("session_token")
	if err != nil {
		return nil, nil
	}
	session, err := app.DB.GetUserIDByToken(token.Value)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   "",
				Expires: time.Now(),
			})
			app.redirect(w, r)
		} else {
			return nil, err
		}
	}
	return session, nil
}

/*############################################################################################################*/

func (app *Application) SignUpPost(w http.ResponseWriter, r *http.Request) {
	session, err := app.checkerSession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if session != nil {
		r.Method = http.MethodGet
		http.Redirect(w, r, "/logout", http.StatusPermanentRedirect)
		return
	}
	user := models.User{
		UserName: r.FormValue("username"),
		Gmail:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	err = app.DB.InsertUser(user)
	if err != nil {
		app.serverError(w, err)
	}
	r.Method = http.MethodGet
	http.Redirect(w, r, "/signin", http.StatusPermanentRedirect)
	return
}
func (app *Application) SignUpGet(w http.ResponseWriter, r *http.Request, s []string) {
	session, err := app.checkerSession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if session != nil {
		r.Method = http.MethodGet
		http.Redirect(w, r, "/logout", http.StatusPermanentRedirect)
		return
	}
	structure := TemplateStructure{}
	if session != nil {
		structure.Signed = true
	}
	templates, err := template.ParseFiles(s...)
	if err != nil {
		app.serverError(w, err)
	}
	if err := templates.Execute(w, structure); err != nil {
		app.serverError(w, err)
		return
	}
	return
}

/*############################################################################################################*/

func (app *Application) SignInPost(w http.ResponseWriter, r *http.Request) {
	session, err := app.checkerSession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if session != nil {
		r.Method = http.MethodGet
		http.Redirect(w, r, "/logout", http.StatusPermanentRedirect)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	session, err = app.DB.GetUser(username, password)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			r.Method = http.MethodGet
			structure := TemplateStructure{Err: errorMessage}
			if session != nil {
				structure.Signed = true
			}
			templates, err := template.ParseFiles("./ui/templates/signin.html", "./ui/templates/header.html")
			if err != nil {
				app.serverError(w, err)
			}
			if err := templates.Execute(w, structure); err != nil {
				app.serverError(w, err)
				return
			}
			return
		}
		app.serverError(w, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   session.Token,
		Expires: session.ExpirationDate,
	})
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	return
}
func (app *Application) SignInGet(w http.ResponseWriter, r *http.Request, s []string) {
	session, err := app.checkerSession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if session != nil {
		r.Method = http.MethodGet
		http.Redirect(w, r, "/logout", http.StatusPermanentRedirect)
		return
	}
	templates, err := template.ParseFiles(s...)
	if err != nil {
		app.serverError(w, err)
	}
	if err := templates.Execute(w, nil); err != nil {
		app.serverError(w, err)
		return
	}
	return
}

func (app *Application) CreatePostPost(w http.ResponseWriter, r *http.Request) {
	session, err := app.checkerSession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if session == nil {
		app.redirect(w, r)
		return
	}

	post := models.Post{
		UserId:   session.UserID,
		UserName: session.UserName,
		Text:     r.FormValue("text"),
		Category: r.FormValue("category"),

		Title: r.FormValue("title"),
	}
	err = app.DB.InsertPost(post)
	if err != nil {
		app.serverError(w, err)
	}
	r.Method = http.MethodGet
	url := fmt.Sprintf("/createPost")
	http.Redirect(w, r, url, http.StatusPermanentRedirect)

}
func (app *Application) CreatePostGet(w http.ResponseWriter, r *http.Request, s []string) {
	session, err := app.checkerSession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if session == nil {
		app.redirect(w, r)
		return
	}
	templates, err := template.ParseFiles(s...)
	if err != nil {
		app.serverError(w, err)
	}
	structure := TemplateStructure{}
	if session != nil {
		structure.Signed = true
	}
	if err := templates.Execute(w, structure); err != nil {
		app.serverError(w, err)
		return
	}
	return
}

func (app *Application) HomeGet(w http.ResponseWriter, r *http.Request, s []string) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	session, err := app.checkerSession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	action := r.URL.Query().Get("action")
	category := r.URL.Query().Get("category")
	postIdStr := r.URL.Query().Get("postId")
	reactStr := r.URL.Query().Get("reaction")
	if action == "reaction" {
		if session != nil {
			postId, err := strconv.Atoi(postIdStr)
			if err == nil {
				reaction, err := strconv.Atoi(reactStr)
				if err == nil {
					err = app.DB.ReactPost(session.UserID, postId, reaction)
					if err != nil {
						app.serverError(w, err)
						return
					}
				}
			}
		} else {
			app.redirect(w, r)
			return
		}
		url := fmt.Sprintf("/?category=%s", category)
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
	structure := TemplateStructure{}
	if category == "" {
		structure.Posts, err = app.DB.GetPostAll()
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				structure.Posts = nil
			} else {
				app.serverError(w, err)
				return
			}
		}
	} else if category != "" {
		if category == "Golang" || category == "Rust" || category == "JS" {
			structure.Posts, err = app.DB.GetPostCategories(category)
			if err != nil {
				if errors.Is(err, models.ErrNoRecord) {
					structure.Posts = nil
				} else {
					app.serverError(w, err)
					return
				}
			}
		} else {
			app.notFound(w)
			return
		}
	}
	templates, err := template.ParseFiles(s...)
	if err != nil {
		app.serverError(w, err)
	}
	if session != nil {
		structure.Signed = true
	}
	if structure.Posts == nil {
		structure.Err = "There is no posts yet"
	}
	if err := templates.Execute(w, structure); err != nil {
		app.serverError(w, err)
		return
	}
	return
}

func (app *Application) PostPost(w http.ResponseWriter, r *http.Request) {
	session, err := app.checkerSession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	comment := models.Comment{
		UserId:   session.UserID,
		UserName: session.UserName,
		Text:     r.FormValue("text"),
	}
	postIdStr := r.URL.Query().Get("postId")
	if postIdStr == "" {
		app.clientError(w, http.StatusNotFound)
		return
	}
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}
	if err = app.DB.InsertComment(postId, comment.UserId, comment.UserName, comment.Text); err != nil {
		app.serverError(w, err)
		return
	}
	r.Method = http.MethodGet
	url := fmt.Sprintf("/post?postId=%s", postIdStr)
	http.Redirect(w, r, url, http.StatusPermanentRedirect)

}
func (app *Application) PostGet(w http.ResponseWriter, r *http.Request, s []string) {
	session, err := app.checkerSession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	action := r.URL.Query().Get("action")
	postIdStr := r.URL.Query().Get("postId")
	if postIdStr == "" {
		app.clientError(w, http.StatusNotFound)
		return
	}
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}
	reactStr := r.URL.Query().Get("reaction")
	if action == "reactionPost" {
		if session != nil {

			reaction, err := strconv.Atoi(reactStr)
			if err == nil {
				err = app.DB.ReactPost(session.UserID, postId, reaction)
				if err != nil {
					app.serverError(w, err)
					return
				}
			}

		} else {
			app.redirect(w, r)
			return
		}
		url := fmt.Sprintf("/post?postId=%s", postIdStr)
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	} else if action == "reactionComment" {

		if session != nil {
			commentIdStr := r.URL.Query().Get("commentId")
			commentId, err := strconv.Atoi(commentIdStr)
			if err == nil {
				reaction, err := strconv.Atoi(reactStr)
				if err == nil {
					err = app.DB.ReactComment(session.UserID, commentId, reaction)
					if err != nil {
						app.serverError(w, err)
						return
					}
				}
			}
		} else {
			app.redirect(w, r)
			return
		}
		url := fmt.Sprintf("/post?postId=%s", postIdStr)
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
	structure := TemplateStructure{}
	if session != nil {
		structure.Signed = true
	}
	structure.Post, err = app.DB.GetPost(postId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	structure.Comments, err = app.DB.GetComments(postId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
		} else {
			app.serverError(w, err)
			return
		}
	}
	templates, err := template.ParseFiles(s...)
	if err != nil {
		app.serverError(w, err)
	}
	if err := templates.Execute(w, structure); err != nil {
		app.serverError(w, err)
		return
	}
	return
}

func (app *Application) LikedPostGet(w http.ResponseWriter, r *http.Request, s []string) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	session, err := app.checkerSession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	action := r.URL.Query().Get("action")
	postIdStr := r.URL.Query().Get("postId")
	reactStr := r.URL.Query().Get("reaction")
	if action == "reaction" {
		if session != nil {
			postId, err := strconv.Atoi(postIdStr)
			if err == nil {
				reaction, err := strconv.Atoi(reactStr)
				if err == nil {
					err = app.DB.ReactPost(session.UserID, postId, reaction)
					if err != nil {
						app.serverError(w, err)
						return
					}
				}
			}
		} else {
			app.redirect(w, r)
			return
		}

		http.Redirect(w, r, "/likedPosts", http.StatusPermanentRedirect)
	}
	structure := TemplateStructure{}
	if session != nil {
		structure.Signed = true
	}
	if structure.Posts == nil {
		structure.Err = "There is no posts yet"
	}

	structure.Posts, err = app.DB.GetPostLiked(session.UserID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			structure.Posts = nil
		} else {
			app.serverError(w, err)
			return
		}
	}
	templates, err := template.ParseFiles(s...)
	if err != nil {
		app.serverError(w, err)
	}
	if err := templates.Execute(w, structure); err != nil {
		app.serverError(w, err)
		return
	}
	return
}

func (app *Application) CreatedPostGet(w http.ResponseWriter, r *http.Request, s []string) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	session, err := app.checkerSession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	action := r.URL.Query().Get("action")
	postIdStr := r.URL.Query().Get("postId")
	reactStr := r.URL.Query().Get("reaction")
	if action == "reaction" {
		if session != nil {
			postId, err := strconv.Atoi(postIdStr)
			if err == nil {
				reaction, err := strconv.Atoi(reactStr)
				if err == nil {
					err = app.DB.ReactPost(session.UserID, postId, reaction)
					if err != nil {
						app.serverError(w, err)
						return
					}
				}
			}
		} else {
			app.redirect(w, r)
			return
		}

		http.Redirect(w, r, "/createdPosts", http.StatusPermanentRedirect)
	}
	structure := TemplateStructure{}
	if session != nil {
		structure.Signed = true
	}
	if structure.Posts == nil {
		structure.Err = "There is no posts yet"
	}
	structure.Posts, err = app.DB.GetPostCreated(session.UserID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			structure.Posts = nil
		} else {
			app.serverError(w, err)
			return
		}
	}
	templates, err := template.ParseFiles(s...)
	if err != nil {
		app.serverError(w, err)
	}
	if err := templates.Execute(w, structure); err != nil {
		app.serverError(w, err)
		return
	}
	return
}

/*############################################################################################################*/
