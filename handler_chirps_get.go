package main

import (
	"net/http"
	"sort"

	"github.com/VokalTuna/chirpy/internal/database"
	"github.com/google/uuid"
)

func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIDstring := r.URL.Query().Get("author_id")
	if authorIDstring == "" {
		return uuid.Nil, nil
	}
	authorID, err := uuid.Parse(authorIDstring)
	if err != nil {
		return uuid.Nil, nil
	}
	return authorID, nil
}

func isAscending(r *http.Request) bool {
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "desc" {
		return false
	}
	return true
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authID, err := authorIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
	}

	var dbChirps []database.Chirp

	if authID != uuid.Nil {
		dbChirps, err = cfg.db.GetChirpsByUser(r.Context(), authID)
	} else {
		dbChirps, err = cfg.db.GetAllChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}
	if !isAscending(r) {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	// Retrieves a singel chirp.
	chirpsIdString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpsIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}
