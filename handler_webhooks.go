package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	param := Parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Not able to decode.", err)
		return
	}
	if param.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), param.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "No such user.", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
	}
	w.WriteHeader(http.StatusNoContent)
}
