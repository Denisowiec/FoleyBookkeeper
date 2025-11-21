package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateProject(w http.ResponseWriter, r *http.Request) {
	// Function for handling requests to create projects
	// Requires authentication
	// Takes title (string) and client (string) as json input
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	type projectInputType struct {
		Title    string `json:"title"`
		ClientID string `json:"client_id"`
	}

	projectInput := projectInputType{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&projectInput)
	if err != nil {
		respondWithError(w, "Error decoding user request", http.StatusBadRequest, err)
		return
	}
	clientID, err := uuid.Parse(projectInput.ClientID)
	if err != nil {
		respondWithError(w, "Error processing user request", http.StatusBadRequest, err)
		return
	}

	createProjectParams := db.CreateProjectParams{
		Title:    projectInput.Title,
		ClientID: clientID,
	}

	prj, err := cfg.db.CreateProject(r.Context(), createProjectParams)
	if err != nil {
		respondWithError(w, "Error creating project", http.StatusInternalServerError, err)
		return
	}
	err = respondWithJSON(w, http.StatusAccepted, prj)
	if err != nil {
		respondWithError(w, "Error processing user response", http.StatusInternalServerError, err)
		return
	}
}

func (cfg *apiConfig) handlerGetProjectByTitle(w http.ResponseWriter, r *http.Request) {
	// This function returns a single project referenced by title provided in JSON input
	// If no title is provided, it returns all projects
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	type projectInputType struct {
		ProjectTitle string `json:"title"`
	}
	projectInput := projectInputType{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&projectInput)

	switch {
	case err == io.EOF:
		// If the body is empty we return a list of all projects
		list, err := cfg.db.GetAllProjects(r.Context())
		if err != nil {
			respondWithError(w, "Error contacting database", http.StatusInternalServerError, err)
			return
		}
		err = respondWithJSON(w, http.StatusOK, list)
		if err != nil {
			respondWithError(w, "Unable to process response data", http.StatusInternalServerError, err)
			return
		}
	case err != nil:
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	default:
		// Non-empty body and no errors
		prj, err := cfg.db.GetProjectByTitle(r.Context(), projectInput.ProjectTitle)
		if err != nil {
			respondWithError(w, "Project not found", http.StatusNotFound, err)
			return
		}
		err = respondWithJSON(w, http.StatusOK, prj)
		if err != nil {
			respondWithError(w, "Unable to process response data", http.StatusInternalServerError, err)
			return
		}
	}
}

func (cfg *apiConfig) handlerGetProjectByID(w http.ResponseWriter, r *http.Request) {
	// Function handles requests to get project data
	// Requires authentication
	// Takes project UUID as input in the url
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthrized", http.StatusUnauthorized, err)
		return
	}

	projectID, err := uuid.Parse(r.PathValue("projectid"))
	if err != nil {
		respondWithError(w, "Error processing user request", http.StatusBadRequest, err)
		return
	}

	prj, err := cfg.db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		respondWithError(w, "Project not found", http.StatusNotFound, err)
		return
	}

	err = respondWithJSON(w, http.StatusOK, prj)
	if err != nil {
		respondWithError(w, "Error processing user response", http.StatusInternalServerError, err)
		return
	}
}

func (cfg *apiConfig) handlerUpdateProject(w http.ResponseWriter, r *http.Request) {
	// This handler requests for modifying project data
	// Requires authentification
	// Takes title (string) and client (string) as input
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthrized", http.StatusUnauthorized, err)
		return
	}

	type projectInputType struct {
		ID       uuid.UUID `json:"id"`
		Title    string    `json:"title"`
		ClientID string    `json:"client_id"`
	}
	projectInput := projectInputType{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&projectInput)
	if err != nil {
		respondWithError(w, "Error decoding user request", http.StatusBadRequest, err)
		return
	}
	clientID, err := uuid.Parse(projectInput.ClientID)
	if err != nil {
		respondWithError(w, "Error decoding user request", http.StatusBadRequest, err)
		return
	}

	updateProjectParams := db.UpdateProjectParams{
		ID:       projectInput.ID,
		Title:    projectInput.Title,
		ClientID: clientID,
	}

	prj, err := cfg.db.UpdateProject(r.Context(), updateProjectParams)
	if err != nil {
		respondWithError(w, "Project not found", http.StatusNotFound, err)
		return
	}

	err = respondWithJSON(w, http.StatusAccepted, prj)
	if err != nil {
		respondWithError(w, "Error processing user response", http.StatusInternalServerError, err)
		return
	}
}

func (cfg *apiConfig) handlerDeleteProject(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	projectID, err := uuid.Parse(r.PathValue("projectid"))
	if err != nil {
		respondWithError(w, "Error parsing user input", http.StatusBadRequest, err)
		return
	}

	prj, err := cfg.db.DeleteProject(r.Context(), projectID)
	if err != nil {
		respondWithError(w, "Project not found", http.StatusNotFound, err)
		return
	}
	err = respondWithJSON(w, http.StatusAccepted, prj)
	if err != nil {
		respondWithError(w, "Error processing user response", http.StatusInternalServerError, err)
		return
	}
}
