package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"testForum/pkg/models"
	"time"
)

func (app *Application) MiddleWare(handle http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {

			http.Redirect(w, r, "/signin", http.StatusPermanentRedirect)
			if _, err = w.Write([]byte("Please login")); err != nil {
				app.serverError(w, err)
			}
			return
		}
		token := cookie.Value
		session, err := app.DB.GetUserIDByToken(token)
		if err != nil {
			fmt.Println("3")
			if errors.Is(err, models.ErrNoRecord) {
				http.Redirect(w, r, "/logout", http.StatusPermanentRedirect)
				if _, err = w.Write([]byte("There is no such session")); err != nil {
					app.serverError(w, err)
				}
				return
			}
			app.serverError(w, err)
			return
		}
		if session.ExpirationDate.Before(time.Now()) {
			if session != nil {
				err = app.DB.DeleteToken(token)
				if err != nil {
					fmt.Println("4")
					app.serverError(w, err)
					return
				}
			}
			http.Redirect(w, r, "/signin", http.StatusPermanentRedirect)
			return
		}
		if session != nil && (r.URL.Path == "/signin" || r.URL.Path == "/signup") {
			fmt.Println("5")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		handle.ServeHTTP(w, r)
	}
}
