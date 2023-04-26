package handlers

import (
	"log"
	"os"
	"testForum/pkg/models"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	DB       *models.Model
}

func LoggerCreater() (*log.Logger, *log.Logger) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	return infoLog, errorLog
}
