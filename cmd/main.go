package main

import (
	"flag"
	"net/http"
	"testForum/db"
	"testForum/pkg/handlers"
)

func main() {
	addr := flag.String("addr", ":8000", "Сетевой адрес веб-страницы")
	flag.Parse()
	infoLog, errorLog := handlers.LoggerCreater()
	db, err := db.DB()
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &handlers.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		DB:       db,
	}
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}
	infoLog.Printf("Запуск сервера на localhost%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
