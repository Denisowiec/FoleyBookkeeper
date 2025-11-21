package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
	"github.com/google/uuid"
)

// These functions convert string values in json input into enums
func strToPart(input string) (db.Part, error) {
	switch input {
	case "props":
		return db.PartProps, nil
	case "footsteps":
		return db.PartFootsteps, nil
	case "movements":
		return db.PartMovements, nil
	case "dialogue":
		return db.PartDialogue, nil
	case "adr":
		return db.PartAdr, nil
	case "music":
		return db.PartMusic, nil
	case "background":
		return db.PartBackground, nil
	case "other":
		return db.PartOther, nil
	default:
		return "", fmt.Errorf("part unknown")
	}
}

func strToActivity(input string) (db.Activity, error) {
	switch input {
	case "record":
		return db.ActivityRecord, nil
	case "edit":
		return db.ActivityEdit, nil
	case "service":
		return db.ActivityService, nil
	case "spotting":
		return db.ActivitySpotting, nil
	case "other":
		return db.ActivityOther, nil
	default:
		return "", fmt.Errorf("activity unknown")
	}
}

func (cfg *apiConfig) handlerCreateSession(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	type sessionInputType struct {
		Duration     int64  `json:"duration"`
		EpisodeID    string `json:"episode_id"`
		PartWorkedOn string `json:"part_worked_on"`
		ActivityDone string `json:"activity_done"`
	}

	sessionInput := sessionInputType{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&sessionInput)
	if err != nil {
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	}

	createSessionParams := db.CreateSessionParams{}
	createSessionParams.Duration = sessionInput.Duration
	createSessionParams.EpisodeID, err = uuid.Parse(sessionInput.EpisodeID)
	if err != nil {
		respondWithError(w, "error decoding user input", http.StatusBadRequest, err)
		return
	}
	createSessionParams.PartWorkedOn, err = strToPart(sessionInput.PartWorkedOn)
	if err != nil {
		respondWithError(w, "error decoding user input", http.StatusBadRequest, err)
		return
	}
	createSessionParams.ActivityDone, err = strToActivity(sessionInput.ActivityDone)
	if err != nil {
		respondWithError(w, "error decoding user input", http.StatusBadRequest, err)
		return
	}

	session, err := cfg.db.CreateSession(r.Context(), createSessionParams)
	if err != nil {
		respondWithError(w, "Error contacting database", http.StatusInternalServerError, err)
		return
	}

	err = respondWithJSON(w, http.StatusAccepted, session)
	if err != nil {
		respondWithError(w, "Error processing user response", http.StatusInternalServerError, err)
		return
	}
}

func (cfg *apiConfig) handlerAddUsersToSession(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	type addUsersToSessionInputType struct {
		UserIDs []string `json:"user_ids"`
	}
	input := addUsersToSessionInputType{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&input)
	if err != nil {
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	}

	sessionID, err := uuid.Parse(r.PathValue("sessionid"))

	if err != nil {
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	}

	for _, user := range input.UserIDs {
		id, err := uuid.Parse(user)
		if err != nil {
			respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
			return
		}
		addusersParams := db.AddUserToSessionParams{
			UserID:    id,
			SessionID: sessionID,
		}
		_, err = cfg.db.AddUserToSession(r.Context(), addusersParams)
		if err != nil {
			respondWithError(w, "Error adding user to session", http.StatusInternalServerError, err)
			return
		}
	}

	err = respondWithJSON(w, http.StatusAccepted, input)
	if err != nil {
		respondWithError(w, "Error processing user response", http.StatusInternalServerError, err)
		return
	}
}

func (cfg *apiConfig) handlerGetSession(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	sessionID, err := uuid.Parse(r.PathValue("sessionid"))
	if err != nil {
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	}

	session, err := cfg.db.GetSession(r.Context(), sessionID)
	if err != nil {
		respondWithError(w, "Session not found", http.StatusNotFound, err)
		return
	}

	users, err := cfg.db.GetUsersForSession(r.Context(), sessionID)
	if err != nil {
		respondWithError(w, "Error contacting database", http.StatusInternalServerError, err)
		return
	}

	type getSesRespType struct {
		db.GetSessionRow
		Users []uuid.UUID `json:"users"`
	}

	getSesResp := getSesRespType{
		GetSessionRow: session,
		Users:         users,
	}

	err = respondWithJSON(w, http.StatusOK, getSesResp)
	if err != nil {
		respondWithError(w, "Error processing user response", http.StatusInternalServerError, err)
		return
	}
}
