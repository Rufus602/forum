package handlers

import (
	"groupie_tracker/pkg/structure"
	"html/template"
	"net/http"
	"strconv"
)

//var Api func() []structure.Artist

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	Information, err := structure.JsonReader()
	if err != nil {
		app.serverError(w, err)
		return
	}
	temp, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		app.notFound(w)
		return
	}
	err = temp.ExecuteTemplate(w, "inx", Information)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *Application) artist(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/artist/" {
		app.notFound(w)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	//Information := Api()
	temp, err := template.ParseFiles("./ui/html/artPage.html")
	if err != nil {
		app.notFound(w)
		return
	}
	Information, err := structure.JsonReader()
	if err != nil {
		app.serverError(w, err)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || (id < 1 || id > len(Information)) {
		app.notFound(w)
		return
	}
	err = temp.ExecuteTemplate(w, "artPage", Information[id-1])
	if err != nil {
		app.serverError(w, err)
		return
	}
}

//
//func Jsongiver(app *Application, w http.ResponseWriter) func() []structure.Artist {
//	Information, err := structure.JsonReader()
//	if err != nil {
//		app.serverError(w, err)
//	}
//	return func() []structure.Artist {
//		return Information
//	}
//}
