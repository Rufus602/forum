package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"testForum/pkg/models"
	"text/template"
)

//var Api func() []structure.Artist

var errorMessage = ErrMessage{Err: "There is no such user. Maybe incorrect username or password, or you did not register"}

func (app *Application) createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		app.CreatePostPost(w, r)
	} else if r.Method == http.MethodGet {
		strings := []string{"./ui/templates/header.gohtml", "./ui/templates/createPost.gohtml"}
		app.CreatePostGet(w, r, strings)
	} else {
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		app.notFound(w)
	}
	return
}
func (app *Application) signIn(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("./ui/templates/signin.html", "./ui/templates/header.html")
	if err != nil {
		app.notFound(w)
		return
	}
	if r.Method != http.MethodPost {
		temp.Execute(w, nil)
		return
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		app.serverError(w, err)
	}
	if resp.StatusCode == http.StatusNotAcceptable {
		temp.Execute(w, err)
	}
	// Close response body as required.
	defer resp.Body.Close()

	info := models.User{
		UserName: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	session, err := app.DB.GetUser(info.UserName, info.Password)
	if err != nil {
		fmt.Println("srgr")
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/signIn", http.StatusNotAcceptable)
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
}
func (app *Application) signUp(w http.ResponseWriter, r *http.Request) {
}
func (app *Application) logout(w http.ResponseWriter, r *http.Request) {
	err := app.checkerSession(w, r)

}

/*############################################################################################################*/
func (app *Application) likedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.notFound(w)
		return
	}
	strings := []string{"./ui/templates/header.gohtml", "./ui/templates/posts.gohtml"}
	app.LikedPostGet(w, r, strings)
	return
}

func (app *Application) createdPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.notFound(w)
		return
	}
	strings := []string{"./ui/templates/header.gohtml", "./ui/templates/posts.gohtml"}
	app.CreatedPostGet(w, r, strings)
	return
}
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.notFound(w)
		return
	}
	strings := []string{"./ui/templates/header.gohtml", "./ui/templates/category.gohtml", "./ui/templates/posts.gohtml"}
	app.HomeGet(w, r, strings)
	return
}
func (app *Application) post(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		app.PostPost(w, r)
	} else if r.Method == http.MethodGet {
		strings := []string{"./ui/templates/header.gohtml", "./ui/templates/comment.gohtml", "./ui/templates/commentPost.gohtml", "./ui/templates/posts.gohtml"}
		app.PostGet(w, r, strings)
	} else {
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		app.notFound(w)
	}
	return
}
