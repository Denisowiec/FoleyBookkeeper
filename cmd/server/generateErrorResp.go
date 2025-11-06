package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func generateErrorResp(s string) []byte {
	errBody := ErrorResponse{
		Error: s,
	}
	dat, err := json.Marshal(errBody)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %s", err)
	}
	return dat
}

func respondWithError(w http.ResponseWriter, message string, errorCode int, err error) {
	// Generates a JSON error message and logs the error to stdout
	if err != nil {
		log.Println(message, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	dat := generateErrorResp(message)
	w.Write(dat)
}
