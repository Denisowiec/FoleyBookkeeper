package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateClient(w http.ResponseWriter, r *http.Request) {
	// Function for handling requests to create clients
	// Requires authentication
	// Takes a name for the client (string), an email-address (string), and optionally some notes
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	type clientInputType struct {
		ClientName string `json:"client_name"`
		Email      string `json:"email"`
		Notes      string `json:"notes"`
	}

	clientInput := clientInputType{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&clientInput)
	if err != nil {
		respondWithError(w, "Error decoding user request", http.StatusBadRequest, err)
		return
	}

	createClientParams := db.CreateClientParams{
		ClientName: clientInput.ClientName,
		Email:      sql.NullString{String: clientInput.ClientName, Valid: true},
		Notes:      sql.NullString{String: clientInput.Notes, Valid: true},
	}

	client, err := cfg.db.CreateClient(r.Context(), createClientParams)
	if err != nil {
		respondWithError(w, "Error creating project", http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	dat, err := json.Marshal(client)
	if err != nil {
		dat = []byte{}
	}

	w.Write(dat)
}

func (cfg *apiConfig) handlerGetClientByID(w http.ResponseWriter, r *http.Request) {
	// Function for handling requests to get client data by it's id
	// Requires authentification
	// Takes an id for the client in the url
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	clientID, err := uuid.Parse(r.PathValue("clientid"))
	if err != nil {
		respondWithError(w, "Error parsing request", http.StatusBadRequest, err)
		return
	}

	client, err := cfg.db.GetClientByID(r.Context(), clientID)
	if err != nil {
		respondWithError(w, "Error fetching data from database", http.StatusInternalServerError, err)
		return
	}

	dat, err := json.Marshal(client)
	if err != nil {
		respondWithError(w, "Unable to process response data", http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (cfg *apiConfig) handlerGetClientByName(w http.ResponseWriter, r *http.Request) {
	// This returns client info by it's name given in the json input
	// If no name is provided, it returns a list of all clients
	// Requires authentification
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}
	var dat []byte

	type clientInputType struct {
		ClientName string `json:"client_name"`
	}
	clientInput := clientInputType{}

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&clientInput)

	switch {
	case err == io.EOF:
		// If the body is empty we return a list of all clients
		list, err := cfg.db.GetAllClients(r.Context())
		if err != nil {
			respondWithError(w, "Error contacting the database", http.StatusInternalServerError, err)
			return
		}
		dat, err = json.Marshal(list)
		if err != nil {
			respondWithError(w, "Error processing data", http.StatusInternalServerError, err)
			return
		}
	case err != nil:
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	default:
		// Body not empty, no errors found
		client, err := cfg.db.GetClientByName(r.Context(), clientInput.ClientName)
		if err != nil {
			respondWithError(w, "Client not found", http.StatusNotFound, err)
			return
		}
		dat, err = json.Marshal(client)
		if err != nil {
			respondWithError(w, "Error processing data", http.StatusInternalServerError, err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (cfg *apiConfig) handlerUpdateClient(w http.ResponseWriter, r *http.Request) {
	// Function handling updates to client info
	// Requires authentification
	// Takes the client's ID, a name for the client (string), an email-address (string)
	// and optionally some notes
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	clientID, err := uuid.Parse(r.PathValue("clientid"))
	if err != nil {
		respondWithError(w, "Error processing user request", http.StatusBadRequest, err)
		return
	}

	type clientInputType struct {
		ClientName string `json:"client_name"`
		Email      string `json:"email"`
		Notes      string `json:"notes"`
	}

	clientInput := clientInputType{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&clientInput)
	if err != nil {
		respondWithError(w, "Error decoding user request", http.StatusBadRequest, err)
		return
	}

	updateClientParams := db.UpdateClientParams{
		ID:         clientID,
		ClientName: clientInput.ClientName,
		Email:      sql.NullString{String: clientInput.Email, Valid: true},
		Notes:      sql.NullString{String: clientInput.Notes, Valid: true},
	}

	client, err := cfg.db.UpdateClient(r.Context(), updateClientParams)
	if err != nil {
		respondWithError(w, "Error creating project", http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	dat, err := json.Marshal(client)
	if err != nil {
		dat = []byte{}
	}

	w.Write(dat)
}

func (cfg *apiConfig) handlerDeleteClient(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	clientID, err := uuid.Parse(r.PathValue("clientid"))
	if err != nil {
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	}

	client, err := cfg.db.DeleteClient(r.Context(), clientID)
	if err != nil {
		respondWithError(w, "Client not found", http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	dat, err := json.Marshal(client)
	if err != nil {
		dat = []byte{}
	}

	w.Write(dat)
}
