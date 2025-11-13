package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func commandCreateProject(cfg *config, args []string) error {
	// Command hanldes creating new projects
	// Takes the project title and client name as arguments
	if len(args) < 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	// First we need to fetch the client id from the db
	url := fmt.Sprintf("%s/api/clients", cfg.serverAddress)

	type getClientReqType struct {
		ClientName string `json:"client_name"`
	}
	getClientReq := getClientReqType{ClientName: args[1]}

	resp1, err := sendRequest(getClientReq, "GET", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp1.Body.Close()

	// The api returns a complete client object, but in this case we only care about the ID
	type getClientResponseType struct {
		ID    uuid.UUID `json:"id"`
		Error string    `json:"error"`
	}
	getClientResp := getClientResponseType{}

	err = processResponse(resp1, &getClientResp)
	if err != nil {
		return err
	}

	if resp1.StatusCode != http.StatusOK {
		return fmt.Errorf(getClientResp.Error)
	}

	// Now we can create the project
	url = fmt.Sprintf("%s/api/projects", cfg.serverAddress)

	type createProjectReqType struct {
		Title    string `json:"title"`
		ClientID string `json:"client_id"`
	}
	createProjectReq := createProjectReqType{
		Title:    args[0],
		ClientID: getClientResp.ID.String(),
	}

	resp2, err := sendRequest(createProjectReq, "POST", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	type createProjectRespType struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Title     string    `json:"title"`
		ClientID  uuid.UUID `json:"client_id"`
		Error     string    `json:"error"`
	}

	createProjectResp := createProjectRespType{}
	err = processResponse(resp2, &createProjectResp)
	if err != nil {
		return err
	}

	if resp2.StatusCode != http.StatusCreated {
		return fmt.Errorf(createProjectResp.Error)
	}

	fmt.Printf("Project %s created successfully\n", createProjectResp.Title)

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
		url := fmt.Sprintf("%s/api/clients", cfg.serverAddress)

		type getClientReqType struct {
			ClientName string `json:"client_name"`
		}
		getClientReq := getClientReqType{ClientName: args[2]}

		resp1, err := sendRequest(getClientReq, "GET", url, cfg.jwt)
		if err != nil {
			return err
		}
		defer resp1.Body.Close()

		// The api returns a complete client object, but in this case we only care about the ID
		type getClientResponseType struct {
			ID    uuid.UUID `json:"id"`
			Error string    `json:"error"`
		}
		getClientResp := getClientResponseType{}

		err = processResponse(resp1, &getClientResp)
		if err != nil {
			return err
		}

		if resp1.StatusCode != http.StatusOK {
			return fmt.Errorf(getClientResp.Error)
		}

		clientID = getClientResp.ID.String()
	}

	// Now we fetch the current project's data
	url := fmt.Sprintf("%s/api/projects", cfg.serverAddress)

	type getPrjReqType struct {
		ProjectTitle string `json:"title"`
	}
	getPrjReq := getPrjReqType{
		ProjectTitle: args[0],
	}

	resp2, err := sendRequest(getPrjReq, "GET", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()
	type getPrjRespType struct {
		ID       uuid.UUID `json:"id"`
		Title    string    `json:"title"`
		ClientID uuid.UUID `json:"client_id"`
		Error    string    `json:"error"`
	}
	getPrjResp := getPrjRespType{}

	err = processResponse(resp2, &getPrjResp)
	if err != nil {
		return err
	}

	if resp2.StatusCode != http.StatusOK {
		return fmt.Errorf(getPrjResp.Error)
	}

	// Now we can update the project
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
		updateProjectReq.ClientID = getPrjResp.ClientID.String()
	}

	resp3, err := sendRequest(updateProjectReq, "POST", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp3.Body.Close()

	type updateProjectRespType struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Title     string    `json:"title"`
		ClientID  uuid.UUID `json:"client_id"`
		Error     string    `json:"error"`
	}

	updateProjectResp := updateProjectRespType{}
	err = processResponse(resp2, &updateProjectResp)
	if err != nil {
		return err
	}

	if resp3.StatusCode != http.StatusAccepted {
		return fmt.Errorf(updateProjectResp.Error)
	}

	fmt.Printf("Project %s updated successfully", updateProjectResp.Title)

	return nil
}

func commandGetProjectInfo(cfg *config, args []string) error {
	// Command displays basic information about a given project
	if len(args) < 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	// We first fetch the project info
	url := fmt.Sprintf("%s/api/projects", cfg.serverAddress)

	type getPrjReqType struct {
		ProjectTitle string `json:"title"`
	}
	getPrjReq := getPrjReqType{
		ProjectTitle: args[0],
	}

	resp, err := sendRequest(getPrjReq, "GET", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	type getPrjRespType struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Title     string    `json:"title"`
		ClientID  uuid.UUID `json:"client_id"`
		Error     string    `json:"error"`
	}
	getPrjResp := getPrjRespType{}

	err = processResponse(resp, &getPrjResp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(getPrjResp.Error)
	}

	// Now we fetch the name of the client
	url = fmt.Sprintf("%s/api/clients/%s", cfg.serverAddress, getPrjResp.ClientID.String())

	resp2, err := sendEmptyRequest("GET", url, cfg.jwt)
	if err != nil {
		return err
	}

	type getClientRespType struct {
		ClientName string `json:"client_name"`
		Error      string `json:"error"`
	}
	getClientResp := getClientRespType{}

	err = processResponse(resp2, &getClientResp)
	if err != nil {
		return err
	}

	if resp2.StatusCode != http.StatusOK {
		return fmt.Errorf(getClientResp.Error)
	}

	// Now we format the output
	fmt.Printf("Project title: %s\nProjectID: %s\nCreated at: %v\nUpdated at: %v\nClient: %s\n",
		getPrjResp.Title, getPrjResp.ID.String(), getPrjResp.CreatedAt, getPrjResp.UpdatedAt, getClientResp.ClientName)
	return nil
}
