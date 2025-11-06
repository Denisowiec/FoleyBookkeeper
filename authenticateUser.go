package main

import (
	"net/http"

	"github.com/Denisowiec/FoleyBookkeeper/internal/auth"
	"github.com/google/uuid"
)

func authenticateUser(r *http.Request, secret string) (uuid.UUID, string, error) {
	// This function returns the uid encoded in the jwt provided and
	// check whether it was encoded using our secret code
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.UUID{}, "", err
	}
	inputUID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		return uuid.UUID{}, "", err
	}

	return inputUID, token, nil
}
