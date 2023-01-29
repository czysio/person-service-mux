package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	sqlc "github.com/czysio/person-service-mux/db/sqlc"
	"github.com/czysio/person-service-mux/schemas"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type PersonController struct {
	Router  *mux.Router
	Queries sqlc.Queries
	Ctx     context.Context
}

func (pc *PersonController) getPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid person ID")
		return
	}

	person, err := pc.Queries.GetPersonById(pc.Ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Person not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, person)
}

func (pc *PersonController) getPeople(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	offset, _ := strconv.Atoi(r.FormValue("offset"))

	if limit > 10 || limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	args := sqlc.GetPeopleParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	people, err := pc.Queries.GetPeople(pc.Ctx, args)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, people)
}

func (pc *PersonController) createPerson(w http.ResponseWriter, r *http.Request) {
	var p schemas.CreatePerson
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	now := time.Now()
	args := sqlc.CreatePersonParams{
		FirstName: p.FirstName,
		Surname:   p.Surname,
		Email:     p.Email,
		Nickname:  p.Nickname,
		CreatedAt: now,
		UpdatedAt: now,
	}

	person, err := pc.Queries.CreatePerson(pc.Ctx, args)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, person)
}

func (pc *PersonController) updatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid person ID")
		return
	}

	var p schemas.UpdatePerson
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	now := time.Now()
	args := sqlc.UpdatePersonParams{
		ID:        id,
		FirstName: sql.NullString{String: p.FirstName, Valid: p.FirstName != ""},
		Surname:   sql.NullString{String: p.Surname, Valid: p.Surname != ""},
		Email:     sql.NullString{String: p.Email, Valid: p.Email != ""},
		Nickname:  sql.NullString{String: p.Nickname, Valid: p.Nickname != ""},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}

	person, err := pc.Queries.UpdatePerson(pc.Ctx, args)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, person)
}

func (pc *PersonController) deletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	if err := pc.Queries.DeletePerson(pc.Ctx, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (pc *PersonController) InitializeRoutes() {
	pc.Router.HandleFunc("/people", pc.getPeople).Methods("GET")
	pc.Router.HandleFunc("/person", pc.createPerson).Methods("POST")
	pc.Router.HandleFunc("/person/{id:[0-9]+}", pc.getPerson).Methods("GET")
	pc.Router.HandleFunc("/person/{id:[0-9]+}", pc.updatePerson).Methods("PUT")
	pc.Router.HandleFunc("/person/{id:[0-9]+}", pc.deletePerson).Methods("DELETE")
}
