package main

import (
	"encoding/json"
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
		Title  string `json:"title"`
		Client string `json:"client"`
	}

	projectInput := projectInputType{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&projectInput)
	if err != nil {
		respondWithError(w, "Error decoding user request", http.StatusBadRequest, err)
		return
	}

	createProjectParams := db.CreateProjectParams{
		Title:  projectInput.Title,
		Client: projectInput.Client,
	}

	prj, err := cfg.db.CreateProject(r.Context(), createProjectParams)
	if err != nil {
		respondWithError(w, "Error creating project", http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	dat, err := json.Marshal(prj)
	if err != nil {
		dat = []byte{}
	}

	w.Write(dat)
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

	projectID := r.PathValue("projectid")

	prj, err := cfg.db.GetProjectByTitle(r.Context(), projectID)
	if err != nil {
		respondWithError(w, "Project not found", http.StatusNotFound, err)
		return
	}

	dat, err := json.Marshal(prj)
	if err != nil {
		respondWithError(w, "Unable to process response data", http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(dat)
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
		ID     uuid.UUID `json:"id"`
		Title  string    `json:"title"`
		Client string    `json:"client"`
	}
	projectInput := projectInputType{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&projectInput)
	if err != nil {
		respondWithError(w, "Error decoding user request", http.StatusBadRequest, err)
		return
	}

	updateProjectParams := db.UpdateProjectParams{
		ID:     projectInput.ID,
		Title:  projectInput.Title,
		Client: projectInput.Client,
	}

	prj, err := cfg.db.UpdateProject(r.Context(), updateProjectParams)
	if err != nil {
		respondWithError(w, "Project not found", http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	dat, err := json.Marshal(prj)
	if err != nil {
		dat = []byte{}
	}
	w.Write(dat)
}

func (cfg *apiConfig) handlerGetAllProjects(w http.ResponseWriter, r *http.Request) {
	// This returns all projects
	// Requires authentification
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	prjs, err := cfg.db.GetAllProjects(r.Context())
	if err != nil {
		respondWithError(w, "Error fetching project from database", http.StatusInternalServerError, err)
		return
	}

	dat, err := json.Marshal(prjs)
	if err != nil {
		respondWithError(w, "Error processing data", http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
