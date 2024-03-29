package handlers

import (
	"net/http"
	"time"
)

func (app *Application) signIn(w http.ResponseWriter, r *http.Request) {
	strings := []string{"./ui/templates/signin.html", "./ui/templates/header.html", "./ui/templates/footer.html"}
	if r.Method == http.MethodPost {
		app.SignInPost(w, r, strings)
	} else if r.Method == http.MethodGet {
		app.SignInGet(w, r, strings)
	} else {
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		app.notFound(w)
	}
	return
}

func (app *Application) signUp(w http.ResponseWriter, r *http.Request) {
	strings := []string{"./ui/templates/signup.html", "./ui/templates/header.html", "./ui/templates/footer.html"}
	if r.Method == http.MethodPost {
		app.SignUpPost(w, r, strings)
	} else if r.Method == http.MethodGet {
		app.SignUpGet(w, r, strings)
	} else {
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		app.notFound(w)
	}
	return
}

func (app *Application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Minute),
	})

	r.Method = http.MethodGet
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

/*############################################################################################################*/
func (app *Application) createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		app.CreatePostPost(w, r)
	} else if r.Method == http.MethodGet {
		strings := []string{
			"./ui/templates/write.html",
			"./ui/templates/header.html",
			"./ui/templates/footer.html",
		}
		app.CreatePostGet(w, r, strings)
	} else {
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		app.notFound(w)
	}
	return
}

func (app *Application) likedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.notFound(w)
		return
	}
	strings := []string{
		"./ui/templates/page-myposts.html",
		"./ui/templates/header.html",

		"./ui/templates/category-topik.html",

		"./ui/templates/footer.html",
	}
	app.LikedPostGet(w, r, strings)
	return
}

func (app *Application) createdPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.notFound(w)
		return
	}
	strings := []string{
		"./ui/templates/page-myposts.html",
		"./ui/templates/header.html",

		"./ui/templates/category-topik.html",

		"./ui/templates/footer.html",
	}
	app.CreatedPostGet(w, r, strings)
	return
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.notFound(w)
		return
	}
	strings := []string{
		"./ui/templates/index.html",
		"./ui/templates/header.html",
		"./ui/templates/footer.html",
		"./ui/templates/buttons.html",
		"./ui/templates/category-topik.html",
	}
	app.HomeGet(w, r, strings)
	return
}

func (app *Application) post(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		app.PostPost(w, r)
	} else if r.Method == http.MethodGet {
		strings := []string{
			"./ui/templates/category-page.html",
			"./ui/templates/header.html",
			"./ui/templates/footer.html",
			"./ui/templates/comment.html",
			"./ui/templates/buttons.html",
			"./ui/templates/category-topik.html",
		}
		app.PostGet(w, r, strings)
	} else {
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		app.notFound(w)
	}
	return
}
