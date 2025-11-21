package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateEpisode(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthirized", http.StatusUnauthorized, err)
		return
	}

	type episodeInputType struct {
		Title         string `json:"title"`
		EpisodeNumber int    `json:"episode_number"`
		ProjectID     string `json:"project_id"`
	}

	episodeInput := episodeInputType{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&episodeInput)
	if err != nil {
		respondWithError(w, "Error decoding user request", http.StatusBadRequest, err)
		return
	}

	projectID, err := uuid.Parse(episodeInput.ProjectID)
	if err != nil {
		respondWithError(w, "Error processing user request", http.StatusBadRequest, err)
		return
	}

	createEpisodeParams := db.CreateEpisodeParams{}
	createEpisodeParams.ProjectID = projectID
	createEpisodeParams.EpisodeNumber = int32(episodeInput.EpisodeNumber)
	if episodeInput.Title != "" {
		createEpisodeParams.Title = sql.NullString{String: episodeInput.Title, Valid: true}
	} else {
		createEpisodeParams.Title.Valid = false
	}

	ep, err := cfg.db.CreateEpisode(r.Context(), createEpisodeParams)
	if err != nil {
		respondWithError(w, "Error creating episode", http.StatusInternalServerError, err)
		return
	}

	err = respondWithJSON(w, http.StatusAccepted, ep)
	if err != nil {
		respondWithError(w, "Error processing response data", http.StatusInternalServerError, err)
		return
	}
}

func (cfg *apiConfig) handlerUpdateEpisode(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthirized", http.StatusUnauthorized, err)
		return
	}

	type episodeInputType struct {
		Title         string `json:"title"`
		EpisodeNumber int    `json:"episode_number"`
		ProjectID     string `json:"project_id"`
	}

	episodeInput := episodeInputType{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&episodeInput)
	if err != nil {
		respondWithError(w, "Error decoding user request", http.StatusBadRequest, err)
		return
	}

	episodeID, err := uuid.Parse(r.PathValue("episodeid"))
	if err != nil {
		respondWithError(w, "Error processing user request", http.StatusBadRequest, err)
		return
	}
	projectID, err := uuid.Parse(episodeInput.ProjectID)
	if err != nil {
		respondWithError(w, "Error processing user request", http.StatusBadRequest, err)
		return
	}

	updateEpisodeParams := db.UpdateEpisodeParams{}
	updateEpisodeParams.ID = episodeID
	updateEpisodeParams.ProjectID = projectID
	updateEpisodeParams.EpisodeNumber = int32(episodeInput.EpisodeNumber)
	if episodeInput.Title != "" {
		updateEpisodeParams.Title = sql.NullString{String: episodeInput.Title, Valid: true}
	} else {
		updateEpisodeParams.Title.Valid = false
	}

	ep, err := cfg.db.UpdateEpisode(r.Context(), updateEpisodeParams)
	if err != nil {
		respondWithError(w, "Error updating episode", http.StatusInternalServerError, err)
		return
	}
	err = respondWithJSON(w, http.StatusAccepted, ep)
	if err != nil {
		respondWithError(w, "Error processing user response", http.StatusInternalServerError, err)
		return
	}
}

func (cfg *apiConfig) handlerGetEpisodeByID(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthrized", http.StatusUnauthorized, err)
		return
	}

	episodeID, err := uuid.Parse(r.PathValue("episodeid"))
	if err != nil {
		respondWithError(w, "Error processing user request", http.StatusBadRequest, err)
		return
	}

	ep, err := cfg.db.GetEpisodeByID(r.Context(), episodeID)
	if err != nil {
		respondWithError(w, "Episode not found", http.StatusNotFound, err)
		return
	}

	err = respondWithJSON(w, http.StatusOK, ep)
	if err != nil {
		respondWithError(w, "Error processing user response", http.StatusInternalServerError, err)
		return
	}
}

func (cfg *apiConfig) handlerGetEpisodesForProject(w http.ResponseWriter, r *http.Request) {
	// Handles requests to list episodes for a project. An episode number can be provided
	// in which case it will return a specific episode
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthrized", http.StatusUnauthorized, err)
		return
	}

	type episodesInputType struct {
		ProjectID     string `json:"project_id"`
		EpisodeNumber int    `json:"episode_number"`
	}
	episodesInput := episodesInputType{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&episodesInput)
	if err != nil {
		respondWithError(w, "Error processing user request", http.StatusBadRequest, err)
		return
	}

	projectID, err := uuid.Parse(episodesInput.ProjectID)
	if err != nil {
		respondWithError(w, "Error processing user request", http.StatusBadRequest, err)
		return
	}

	if episodesInput.EpisodeNumber != 0 {
		// If a number is provided in the input it means that the client want to get as single episode of the given number
		getEpisodeParams := db.GetEpisodeByNumberParams{
			ProjectID:     projectID,
			EpisodeNumber: int32(episodesInput.EpisodeNumber),
		}
		ep, err := cfg.db.GetEpisodeByNumber(r.Context(), getEpisodeParams)
		if err != nil {
			respondWithError(w, "Error getting data from database", http.StatusInternalServerError, err)
			return
		}

		err = respondWithJSON(w, http.StatusOK, ep)
		if err != nil {
			respondWithError(w, "Error processing response data", http.StatusInternalServerError, err)
			return
		}
	} else {
		// Otherwise, we return all of the episodes for the given project
		eps, err := cfg.db.GetAllEpisodesForProject(r.Context(), projectID)
		if err != nil {
			respondWithError(w, "Error getting data from database", http.StatusInternalServerError, err)
			return
		}

		err = respondWithJSON(w, http.StatusOK, eps)
		if err != nil {
			respondWithError(w, "Error processing response data", http.StatusInternalServerError, err)
			return
		}
	}
}

func (cfg *apiConfig) handlerDeleteEpisode(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	episodeID, err := uuid.Parse(r.PathValue("episodeid"))
	if err != nil {
		respondWithError(w, "Error processing user input", http.StatusBadRequest, err)
		return
	}

	ep, err := cfg.db.DeleteEpisode(r.Context(), episodeID)
	if err != nil {
		respondWithError(w, "Episode not found", http.StatusNotFound, err)
		return
	}
	err = respondWithJSON(w, http.StatusAccepted, ep)
	if err != nil {
		respondWithError(w, "Error processing user response", http.StatusInternalServerError, err)
		return
	}
}
