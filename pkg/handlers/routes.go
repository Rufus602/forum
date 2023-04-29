package handlers

import (
	"net/http"
)

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/createPost", app.MiddleWare(app.createPost))
	mux.HandleFunc("/createdPosts", app.MiddleWare(app.createdPosts))
	mux.HandleFunc("/likedPosts", app.MiddleWare(app.likedPosts))
	mux.HandleFunc("/signin", app.signIn)
	mux.HandleFunc("/signup", app.signUp)
	mux.HandleFunc("/logout", app.logout)
	mux.HandleFunc("/post", app.post)

	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fileServer))
	return mux
}
