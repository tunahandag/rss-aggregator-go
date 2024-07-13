package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tunahandag/rss-aggregator/internal/database"
)

func (apiCfg *apiConfig)handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error parsing json")
	}
	user , err := apiCfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: params.Name,

	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating user: %v", err))
	}
	respondWithJSON(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUser(user))
}