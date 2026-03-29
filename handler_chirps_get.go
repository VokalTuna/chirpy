package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	// Retrieve all chirps
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response := []Chirp{}
	for _, dbChirp := range dbChirps {
		response = append(response, Chirp{
			ID:         dbChirp.ID,
			Created_at: dbChirp.CreatedAt,
			Updated_at: dbChirp.UpdatedAt,
			Body:       dbChirp.Body,
			User_id:    dbChirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	// Retrieves a singel chirp.
	chirpsIdString := r.PathValue("ChirpID")
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
		ID:         dbChirp.ID,
		Created_at: dbChirp.CreatedAt,
		Updated_at: dbChirp.UpdatedAt,
		Body:       dbChirp.Body,
		User_id:    dbChirp.UserID,
	})
}
