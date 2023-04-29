package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
)

type errParser struct {
	Number int
	Error  string
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	myErr := errParser{}
	myErr.Number = http.StatusInternalServerError
	myErr.Error = http.StatusText(http.StatusInternalServerError)
	temp, err := template.ParseFiles("./ui/templates/error.html", "./ui/templates/header.html")
	w.WriteHeader(http.StatusInternalServerError)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = temp.Execute(w, myErr)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Помощник clientError отправляет определенный код состояния и соответствующее ее описание
// пользователю. Мы будем использовать это в следующий уроках, чтобы отправлять ответы вроде 400 "Bad
// Request", когда есть проблема с пользовательским запросом.
func (app *Application) clientError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	myErr := errParser{
		status,
		http.StatusText(status),
	}
	// myErr.Number = string(status)
	// myErr.Error = http.StatusText(status)

	temp, err := template.ParseFiles("./ui/templates/error.html", "./ui/templates/header.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = temp.Execute(w, myErr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// http.Error(w, http.StatusText(status), status)
}

// Мы также реализуем помощник notFound. Это просто
// удобная оболочка вокруг clientError, которая отправляет пользователю ответ "404 Страница не найдена".
func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
