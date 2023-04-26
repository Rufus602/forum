package handlers

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"testForum/pkg/models"
	"time"
)

func (app *Application) middleWare(handle http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_cookie")
		if err != nil {
			redirectToken := uuid.NewString()
			http.SetCookie(w, &http.Cookie{
				Name:    "redirect_cookie",
				Value:   redirectToken,
				Expires: time.Now().Add(2 * time.Minute),
			})
			http.Redirect(w, r, "/signIn", http.StatusTemporaryRedirect)
			w.Write([]byte("Please login"))
			return
		}
		token := cookie.Value
		session, err := app.DB.GetUserIDByToken(token)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.Redirect(w, r, "/logout", http.StatusMovedPermanently)
				w.Write([]byte("There is no such session"))
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
			redirectToken := uuid.NewString()
			http.SetCookie(w, &http.Cookie{
				Name:    "redirect_cookie",
				Value:   redirectToken,
				Expires: time.Now().Add(2 * time.Minute),
			})
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
	}
}
