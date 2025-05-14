package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/szymonSadowski/chirpy/internal/auth"
	"github.com/szymonSadowski/chirpy/internal/database"
)

type refreshResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return
	}

	tokenData, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token data", err)
		return
	}
	if tokenData.ExpiresAt.Before(time.Now()) { // Or tokenData.ExpiresAt < time.Now()
		respondWithError(w, http.StatusUnauthorized, "Token expired", nil)
		return
	}
	if tokenData.RevokedAt.Valid { // .Valid is true if RevokedAt is not NULL (meaning it has been set)
		respondWithError(w, http.StatusUnauthorized, "Token revoked", nil)
		return
	}
	newToken, err := auth.MakeJWT(tokenData.UserID, cfg.jwt_secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, refreshResponse{
		Token: newToken,
	})

}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get token", err)
		return
	}
	currentTime := time.Now().UTC()
	_, err = cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{ // Correctly construct sql.NullTime
			Time:  currentTime,
			Valid: true,
		},
		UpdatedAt:    currentTime,
		RefreshToken: token,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
