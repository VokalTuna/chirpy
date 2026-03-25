package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func (cfg *apiConfig) handlerValidation(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		rw.WriteHeader(500)
		return
	}
	if len(params.Body) > 140 {
		respondWithError(rw, 400, "Chirp is too long")
		return
	}

	output := profaneFilter(params.Body)

	type response struct {
		Cleaned_body string `json:"cleaned_body"`
	}
	respondWithJSON(rw, http.StatusOK, response{Cleaned_body: output})
}

func profaneFilter(msg string) string {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	spltMsg := strings.Split(msg, " ")
	compMsg := strings.ToLower(msg)
	compArr := strings.Split(compMsg, " ")

	for i, val := range compArr {
		if slices.Contains(profanity, val) {
			spltMsg[i] = "****"
		}
	}

	return strings.Join(spltMsg, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnVal struct {
		Error string `json:"error"`
	}
	respBody := returnVal{
		Error: msg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error mashalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
