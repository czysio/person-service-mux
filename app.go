package main

import (
	"context"
	"database/sql"
	"fmt"
	pc "github.com/czysio/person-service-mux/controllers"
	sqlc "github.com/czysio/person-service-mux/db/sqlc"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	router *mux.Router
	db     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	ctx := context.TODO()
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	queries := sqlc.New(a.db)
	pc := pc.PersonController{
		Router:  a.router,
		Queries: *queries,
		Ctx:     ctx,
	}

	a.router = mux.NewRouter()

	pc.InitializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", a.router))
}
