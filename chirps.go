package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/szymonSadowski/chirpy/internal/auth"
	"github.com/szymonSadowski/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwt_secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}
	if userID == uuid.Nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	chirp, err := cfg.ValidateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't validate chirp", err)
		return
	}

	createdChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		UserID:    uuid.MustParse(userID.String()),
		Body:      chirp,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	// Using a response struct with proper JSON field tags
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	})
}

type DBChirp struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	UserID    uuid.UUID `db:"user_id"`
	Body      string    `db:"body"`
}

func (cfg *apiConfig) handlerGetChrips(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sortOrder := strings.ToLower(r.URL.Query().Get("sort"))

	var dbChirps []database.Chirp
	var err error

	if authorID == "" {
		dbChirps, err = cfg.db.GetAllChirps(r.Context())
	} else {
		uid, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id", parseErr)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByUserID(r.Context(), uid)
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	// Convert []database.Chirp to []Chirp
	chirps := make([]Chirp, 0, len(dbChirps))
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	// Sorting logic
	switch sortOrder {
	case "desc":
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	default:
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID") // Always returns a string

	if chirpIDString == "" {
		http.Error(w, "Chirp ID cannot be empty", http.StatusBadRequest)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), uuid.MustParse(chirpIDString))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	})

}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID") // Always returns a string

	if chirpIDString == "" {
		http.Error(w, "Chirp ID cannot be empty", http.StatusBadRequest)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwt_secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}
	if userID == uuid.Nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", nil)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), uuid.MustParse(chirpIDString))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirp", err)
		return
	}
	if dbChirp.UserID != uuid.MustParse(userID.String()) {
		respondWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp", nil)
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), uuid.MustParse(chirpIDString))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't delete chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
