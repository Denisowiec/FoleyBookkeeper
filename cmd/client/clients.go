package main

import (
	"fmt"
	"net/http"

	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
)

func getClientByName(cfg *config, name string) (db.Client, error) {
	// A helper function, since fetching a client by name is a common necessity in a couple commands

	type getClientReqType struct {
		ClientName string `json:"client_name"`
	}
	getClientReq := getClientReqType{
		ClientName: name,
	}

	return getThing(cfg, "/api/clients", getClientReq, db.Client{})

	/*
		url := fmt.Sprintf("%s/api/clients", cfg.serverAddress)

		resp1, err := sendRequest(getClientReq, "GET", url, cfg.jwt)
		if err != nil {
			return db.Client{}, err
		}
		defer resp1.Body.Close()

		if resp1.StatusCode != http.StatusOK {
			return db.Client{}, processErrorResponse(resp1)
		}

		clientResponse := db.Client{}
		err = processResponse(resp1, &clientResponse)
		if err != nil {
			return db.Client{}, err
		}

		return clientResponse, nil*/
}

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

	if resp.StatusCode != http.StatusCreated {
		return processErrorResponse(resp)
	}

	createClientResponse := db.Client{}
	err = processResponse(resp, &createClientResponse)
	if err != nil {
		return err
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

	client, err := getClientByName(cfg, args[0])
	if err != nil {
		return err
	}

	// Now we update the client data
	url := fmt.Sprintf("%s/api/clients/%s", cfg.serverAddress, client.ID.String())

	type updClientReqType struct {
		ClientNewName string `json:"client_name"`
		Email         string `json:"email"`
		Notes         string `json:"notes"`
	}
	updClientReq := updClientReqType{
		ClientNewName: args[1],
	}

	// We only change the things that were provided as arguments. We leave the rest as it was.
	if len(args) >= 2 {
		updClientReq.Email = args[2]
	} else {
		updClientReq.Email = client.Email.String
	}
	if len(args) >= 3 {
		updClientReq.Notes = args[3]
	} else {
		updClientReq.Notes = client.Notes.String
	}

	resp2, err := sendRequest(updClientReq, "PUT", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusAccepted {
		return processErrorResponse(resp2)
	}

	updClientResp := db.Client{}

	err = processResponse(resp2, &updClientResp)
	if err != nil {
		return err
	}

	fmt.Printf("Client %s updated successfully\n", updClientResp.ClientName)
	return nil
}

func commandGetClient(cfg *config, args []string) error {
	// This function show information about a given client
	if len(args) == 0 {
		return fmt.Errorf("invalid number of arguments")
	}

	client, err := getClientByName(cfg, args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\nID: %s\nCreated at: %v\nUpdated at: %v\nEmail: %s\nNotes: %s\n",
		client.ClientName, client.ID.String(), client.CreatedAt, client.UpdatedAt, client.Email.String, client.Notes.String)

	return nil
}

func commandGetAllClients(cfg *config, args []string) error {
	url := fmt.Sprintf("%s/api/clients", cfg.serverAddress)

	resp, err := sendEmptyRequest("GET", url, cfg.jwt)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return processErrorResponse(resp)
	}

	var list []db.Client
	err = processResponse(resp, &list)
	if err != nil {
		return err
	}

	for _, item := range list {
		fmt.Printf("%s, email: %s, notes: %s\n", item.ClientName, item.Email.String, item.Notes.String)
	}

	return nil
}
