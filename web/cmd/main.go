package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	DB       *sql.DB
	tokens   *MapTokens
}

func main() {
	app := &application{
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		tokens:   newMapTokens(),
	}

	db, err := app.openDB()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	app.DB = db
	defer func() {
		err = app.DB.Close()
		if err != nil {
			app.errorLog.Println(err)
		}
	}()

	srv := &http.Server{
		Addr:     ":8080",
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Println("Start web-server on http://127.0.0.1:8080")
	app.errorLog.Fatal(srv.ListenAndServe())
}
