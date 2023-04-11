package server

import (
	"net/http"
	"testForum/internal/mysql"
	"testForum/internal/transport"
)

func Server() error {
	err := mysql.CreateDB()
	if err != nil {
		return err
	}

	transport.Handlers()

	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		return err
	}

	return nil
}
