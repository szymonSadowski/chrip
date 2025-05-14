package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/szymonSadowski/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	token, err := auth.GetAPIKey(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get key", err)
		return
	}

	if token != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "Invalid key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid user ID format", err)
		return
	}
	fmt.Printf("User ID: %s\n", userID.String())
	err = cfg.db.UpgradeUser(context.Background(),
		userID,
	)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't upgrade user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
