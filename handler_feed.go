package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tunahandag/rss-aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error parsing json")
	}
	feed , err := apiCfg.db.CreateFeed(r.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: params.Name,
		Url: params.URL,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating feed: %v", err))
	}
	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds , err := apiCfg.db.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %v", err))
	}
	respondWithJSON(w, 200, databaseFeedsToFeeds(feeds))
}