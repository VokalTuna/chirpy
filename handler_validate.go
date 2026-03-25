package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerValidation(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type response struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(rw, http.StatusBadRequest, "Couldn't decode parameters", err)
		rw.WriteHeader(500)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(rw, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	output := getCleanBody(params.Body, badWords)
	respondWithJSON(rw, http.StatusOK, response{Cleaned_body: output})
}

func getCleanBody(msg string, badWords map[string]struct{}) string {
	words := strings.Split(msg, "")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, ok := badWords[lowerWord]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
