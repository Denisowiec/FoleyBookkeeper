package main

import (
	"fmt"
	"net/http"

	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
)

func getProjectByName(cfg *config, name string) (db.Project, error) {
	// A helper function, since fetching a project by name is a common necessity in a couple commands
	type getPrjType struct {
		ProjectTitle string `json:"title"`
	}
	getPrjReq := getPrjType{
		ProjectTitle: name,
	}

	return getThing(cfg, "/api/projects", getPrjReq, db.Project{})
}

func commandCreateProject(cfg *config, args []string) error {
	// Command hanldes creating new projects
	// Takes the project title and client name as arguments
	if len(args) < 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	// First we need to fetch the client id from the db
	client, err := getClientByName(cfg, args[1])
	if err != nil {
		return err
	}

	// Now we can create the project
	url := fmt.Sprintf("%s/api/projects", cfg.serverAddress)

	type createProjectReqType struct {
		Title    string `json:"title"`
		ClientID string `json:"client_id"`
	}
	createProjectReq := createProjectReqType{
		Title:    args[0],
		ClientID: client.ID.String(),
	}

	resp, err := sendRequest(createProjectReq, "POST", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return processErrorResponse(resp)
	}

	prj := db.Project{}
	err = processResponse(resp, &prj)
	if err != nil {
		return err
	}

	fmt.Printf("Project %s created successfully\n", prj.Title)

	return nil
}

func commandUpdateProject(cfg *config, args []string) error {
	// Command hanldes creating new projects
	// Takes the project's old title, the new title and optionally client name as arguments
	if len(args) < 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	// First we need to fetch the client id from the db, if a new client name was given as a command argument
	clientID := ""

	if len(args) >= 3 {
		client, err := getClientByName(cfg, args[2])
		if err != nil {
			return err
		}

		clientID = client.ID.String()
	}

	// Now we fetch the current project's data
	oldPrj, err := getProjectByName(cfg, args[0])
	if err != nil {
		return err
	}

	// Now we can update the project

	url := fmt.Sprintf("%s/api/projects", cfg.serverAddress)

	type updateProjectReqType struct {
		Title    string `json:"title"`
		ClientID string `json:"client_id"`
	}
	updateProjectReq := updateProjectReqType{
		Title: args[1],
	}

	if clientID != "" {
		updateProjectReq.ClientID = clientID
	} else {
		updateProjectReq.ClientID = oldPrj.ClientID.String()
	}

	resp, err := sendRequest(updateProjectReq, "POST", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return processErrorResponse(resp)
	}

	prj := db.Project{}

	err = processResponse(resp, &prj)
	if err != nil {
		return err
	}

	fmt.Printf("Project %s updated successfully", prj.Title)

	return nil
}

func commandGetProjectInfo(cfg *config, args []string) error {
	// Command displays basic information about a given project
	if len(args) < 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	// We first fetch the project info
	prj, err := getProjectByName(cfg, args[0])
	if err != nil {
		return err
	}

	// Now we fetch the name of the client
	client, err := getThingByID(cfg, "/api/clients/", prj.ClientID.String(), db.Client{})
	if err != nil {
		return err
	}

	// Now we format the output
	fmt.Printf("Project title: %s\nProjectID: %s\nCreated at: %v\nUpdated at: %v\nClient: %s\n",
		prj.Title, prj.ID.String(), prj.CreatedAt, prj.UpdatedAt, client.ClientName)
	return nil
}

func commandGetAllProjects(cfg *config, args []string) error {
	url := fmt.Sprintf("%s/api/projects", cfg.serverAddress)

	var list []db.Project

	resp, err := sendEmptyRequest("GET", url, cfg.jwt)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return processErrorResponse(resp)
	}

	err = processResponse(resp, &list)
	if err != nil {
		return err
	}

	for _, item := range list {
		client, err := getThingByID(cfg, "/api/clients", item.ClientID.String(), db.Client{})
		if err != nil {
			return err
		}
		fmt.Printf("Name: %s, client: %s\n", item.Title, client.ClientName)
	}
	return nil
}
