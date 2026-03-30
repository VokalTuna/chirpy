package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/VokalTuna/chirpy/internal/auth"
	"github.com/VokalTuna/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}
	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(param.Password)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't hass password", err)
	}

	dbUser, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		HashedPassword: hashedPassword,
		Email:          param.Email,
	})
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(rw, http.StatusCreated, response{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
	})

}
