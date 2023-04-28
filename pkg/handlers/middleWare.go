package handlers

import (
	"errors"
	"net/http"
	"testForum/pkg/models"
	"time"
)

func (app *Application) MiddleWare(handle http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_cookie")
		if err != nil {
			http.Redirect(w, r, "/signIn", http.StatusPermanentRedirect)
			if _, err = w.Write([]byte("Please login")); err != nil {
				app.serverError(w, err)
			}
			return
		}
		token := cookie.Value
		session, err := app.DB.GetUserIDByToken(token)
		if err != nil {
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
					app.serverError(w, err)
					return
				}
			}
			http.Redirect(w, r, "/signin", http.StatusPermanentRedirect)
			return
		}
		if session != nil && (r.URL.Path == "/signIn" || r.URL.Path == "/signup") {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		handle.ServeHTTP(w, r)
	}
}
