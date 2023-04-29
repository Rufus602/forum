package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testForum/pkg/models"
	"text/template"
	"time"
)

var (
	ErrInvalidEmail    = errors.New("Email must consists from letters, at and dot")
	ErrInvalidPassword = errors.New("Password have to contain at least 8 characters but no more than 20; include number, lowercase, uppercase and symbol")
	ErrInvalidUsername = errors.New("Name must consists from letters")
	errorSignIn        = errors.New("There is no such user. Maybe incorrect username or password, or you did not register")
)

func checkInput(u, p string) error {
	if u == "" || p == "" {
		return errorSignIn
	}
	return nil
}

func checkUserInfo(user models.User) error {
	if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(user.Gmail) {
		return ErrInvalidEmail
	}

	for _, w := range user.UserName {
		if w < 32 || w > 126 {
			return ErrInvalidUsername
		}
	}

	if !checkPassword(user.Password) {
		return ErrInvalidPassword
	}

	return nil
}

func checkPassword(password string) bool {
	numbers := "0123456789"
	lowerCase := "qwertyuiopasdfghjklzxcvbnm"
	upperCase := "QWERTYUIOPASDFGHJKLZXCVBNM"
	symbols := "!@#$%^&*()_-+={[}]|\\:;<,>.?/"

	if len(password) < 8 || len(password) > 20 {
		return false
	}

	if !contains(password, numbers) || !contains(password, lowerCase) || !contains(password, upperCase) || !contains(password, symbols) {
		return false
	}

	for _, w := range symbols {
		if w < 32 || w > 126 {
			return false
		}
	}
	return true
}

func contains(s, checkSymbols string) bool {
	for _, w := range checkSymbols {
		if strings.Contains(s, string(w)) {
			return true
		}
	}
	return false
}

var errorMessage = "There is no such user. Maybe incorrect username or password, or you did not register"

func (app *Application) redirect(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodGet

	app.signIn(w, r)
	return
}

func (app *Application) checkerSession(w http.ResponseWriter, r *http.Request) (*models.Session, error) {
	token, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		}
		return nil, err
	}
	session, err := app.DB.GetUserIDByToken(token.Value)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   "",
				Expires: time.Now().Add(-1 * time.Minute),
			})
			return nil, nil
		} else {
			return nil, err
		}
	}
	return session, nil
}

/*############################################################################################################*/

func (app *Application) SignUpPost(w http.ResponseWriter, r *http.Request, s []string) {
	session, err := app.checkerSession(w, r)
	if err != nil {

		app.serverError(w, err)
		return
	}
	if session != nil {
		r.Method = http.MethodGet
		app.logout(w, r)
		return
	}
	user := models.User{
		UserName: r.FormValue("username"),
		Gmail:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	if err = checkUserInfo(user); err != nil {
		r.Method = http.MethodGet
		structure := TemplateStructure{Err: err.Error()}

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
	if err = app.DB.CheckDuplex(user); err != nil {
		r.Method = http.MethodGet
		structure := TemplateStructure{Err: err.Error()}

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

	err = app.DB.InsertUser(user)
	if err != nil {
		app.serverError(w, err)
	}
	r.Method = http.MethodGet
	app.signIn(w, r)
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
		app.logout(w, r)
		return
	}
	structure := TemplateStructure{}
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

func (app *Application) SignInPost(w http.ResponseWriter, r *http.Request, s []string) {
	session, err := app.checkerSession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// r.Method = http.MethodGet
	if session != nil {

		app.logout(w, r)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if err = checkInput(username, password); err != nil {
		r.Method = http.MethodGet
		structure := TemplateStructure{Err: err.Error()}

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

	session, err = app.DB.GetUser(username, password)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {

			templates, err := template.ParseFiles("./ui/templates/signin.html", "./ui/templates/header.html", "./ui/templates/footer.html")
			if err != nil {
				app.serverError(w, err)
			}
			if err := templates.Execute(w, nil); err != nil {
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
	r.Method = http.MethodGet
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
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
		app.logout(w, r)
		return
	}

	templates, err := template.ParseFiles(s...)
	if err != nil {
		app.serverError(w, err)
	}

	structure := TemplateStructure{}
	if session != nil {
		structure.Signed = &models.Signed{true}
	}
	if err := templates.Execute(w, structure); err != nil {
		app.serverError(w, err)
		return
	}
	return
}

/*############################################################################################################*/

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
		Text:     r.FormValue("content"),
		Category: r.FormValue("category"),

		Title: r.FormValue("title"),
	}
	if post.Text == "" || post.Category == "" || post.Title == "" {
		r.Method = http.MethodGet

		http.Redirect(w, r, "/createPost", http.StatusSeeOther)
		return
	}
	err = app.DB.InsertPost(post)
	if err != nil {
		app.serverError(w, err)
	}
	r.Method = http.MethodGet

	app.createPost(w, r)
	return
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
		structure.Signed = &models.Signed{true}
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
	tag := r.URL.Query().Get("tag")
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
		r.Method = http.MethodGet
		url := fmt.Sprintf("?tag=%s", tag)
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
	structure := TemplateStructure{}
	structure.Tag = tag

	if tag == "" {
		structure.Posts, err = app.DB.GetPostAll()

		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				structure.Posts = nil
			} else {
				app.serverError(w, err)
				return
			}
		}

	} else if tag != "" {
		if tag == "golang" || tag == "rust" || tag == "js" {
			structure.Posts, err = app.DB.GetPostCategories(tag)
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
		structure.Signed = &models.Signed{true}
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
	if session == nil {
		r.Method = http.MethodGet
		app.signIn(w, r)
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
	if comment.Text == "" {
	} else if err = app.DB.InsertComment(postId, comment.UserId, comment.UserName, comment.Text); err != nil {
		app.serverError(w, err)
		return
	}
	r.Method = http.MethodGet

	url := fmt.Sprintf("/post?postId=%s", postIdStr)
	http.Redirect(w, r, url, http.StatusSeeOther)
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
		http.Redirect(w, r, url, http.StatusSeeOther)
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
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
	structure := TemplateStructure{}
	if session != nil {
		structure.Signed = &models.Signed{true}
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

		http.Redirect(w, r, "/likedPosts", http.StatusSeeOther)
	}
	structure := TemplateStructure{}
	if session != nil {
		structure.Signed = &models.Signed{true}
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

		http.Redirect(w, r, "/createdPosts", http.StatusSeeOther)
	}
	structure := TemplateStructure{}
	if session != nil {
		structure.Signed = &models.Signed{true}
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
