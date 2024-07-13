package main

import (
	"fmt"
	"net/http"

	"github.com/tunahandag/rss-aggregator/internal/auth"
	"github.com/tunahandag/rss-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey,err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := cfg.db.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Could not get user: %v", err))
			return
		}
		handler(w,r,user)
	}
}