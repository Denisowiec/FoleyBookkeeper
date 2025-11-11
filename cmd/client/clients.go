package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func commandCreateClient(cfg *config, args []string) error {
	// Command handles creating new clients
	if len(args) < 1 {
		return fmt.Errorf("invalid number of arguments")
	}
	url := fmt.Sprintf("%s/api/clients", cfg.serverAddress)

	type reqBodyType struct {
		ClientName string `json:"client_name"`
		Email      string `json:"email"`
		Notes      string `json:"notes"`
	}
	reqBody := reqBodyType{}
	reqBody.ClientName = args[0]

	if len(args) == 2 {
		reqBody.Email = args[1]
	} else if len(args) >= 3 {
		reqBody.Notes = args[2]
	}

	resp, err := sendRequest(reqBody, "POST", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type createClientResponseType struct {
		ID         uuid.UUID      `json:"id"`
		CreatedAt  time.Time      `json:"created_at"`
		UpdatedAt  time.Time      `json:"updated_at"`
		ClientName string         `json:"client_name"`
		Email      sql.NullString `json:"email"`
		Notes      sql.NullString `json:"notes"`
		Error      string         `json:"error"`
	}

	createClientResponse := createClientResponseType{}
	err = processResponse(resp, &createClientResponse)
	if err != nil {
		return err
	}

	// If the REST api respond with something other than 201, somnething went wrong
	// There should be an error message in the response payload
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf(createClientResponse.Error)
	}

	fmt.Printf("Client %s created successfully\n", createClientResponse.ClientName)
	return nil

}

func commandUpdateClient(cfg *config, args []string) error {
	// This function updates a client
	// Takes the old name, new name, email and notes as parameters
	if len(args) < 2 {
		return fmt.Errorf("invalid number of arguments")
	}
	type getClientReqType struct {
		ClientName string `json:"client_name"`
	}
	getClientReq := getClientReqType{
		ClientName: args[0],
	}

	// First, we fetch the original client info from the database
	url := fmt.Sprintf("%s/api/clients", cfg.serverAddress)

	resp1, err := sendRequest(getClientReq, "GET", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp1.Body.Close()

	type getClientResponseType struct {
		ID         uuid.UUID      `json:"id"`
		CreatedAt  time.Time      `json:"created_at"`
		UpdatedAt  time.Time      `json:"updated_at"`
		ClientName string         `json:"client_name"`
		Email      sql.NullString `json:"email"`
		Notes      sql.NullString `json:"notes"`
		Error      string         `json:"error"`
	}

	getClientResponse := getClientResponseType{}

	err = processResponse(resp1, &getClientResponse)
	if err != nil {
		return err
	}

	if resp1.StatusCode != http.StatusOK {
		return fmt.Errorf(getClientResponse.Error)
	}

	// Now we update the client data
	url = fmt.Sprintf("%s/api/clients/%s", cfg.serverAddress, getClientResponse.ID.String())

	type updClientReqType struct {
		ClientNewName string         `json:"client_name"`
		Email         sql.NullString `json:"email"`
		Notes         sql.NullString `json:"notes"`
	}
	updClientReq := updClientReqType{
		ClientNewName: args[1],
	}

	// We only change the things that were provided as arguments. We leave the rest as it was.
	if len(args) >= 2 {
		updClientReq.Email = sql.NullString{String: args[2], Valid: true}
	} else {
		updClientReq.Email = getClientResponse.Email
	}
	if len(args) >= 3 {
		updClientReq.Notes = sql.NullString{String: args[3], Valid: true}
	} else {
		updClientReq.Notes = getClientResponse.Notes
	}

	resp2, err := sendRequest(updClientReq, "PUT", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	updClientResp := getClientResponseType{}

	err = processResponse(resp2, &updClientResp)
	if err != nil {
		return err
	}

	if resp2.StatusCode != http.StatusAccepted {
		return fmt.Errorf(updClientResp.Error)
	}

	fmt.Printf("Client %s updated successfully\n", updClientResp.ClientName)
	return nil
}
