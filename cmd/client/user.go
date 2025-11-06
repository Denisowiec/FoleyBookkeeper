package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func commandCreateUser(cfg *config, args []string) error {
	// Command handles creating new users
	if len(args) < 3 {
		return fmt.Errorf("Invalid number of arguments")
	}
	url := fmt.Sprintf("%s/api/users", cfg.serverAddress)

	type reqBodyType struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	reqBody := reqBodyType{
		Username: args[0],
		Email:    args[1],
		Password: args[2],
	}

	dat, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	payload := bytes.NewBuffer(dat)

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	type createUserResponseType struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Error     string    `json:"error"`
	}

	createUserResponse := createUserResponseType{}
	err = json.Unmarshal(respBody, &createUserResponse)
	if err != nil {
		return err
	}

	// If the REST api respond with something other than 201, somnething went wrong
	// There should be an error message in the response payload
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error processing request: %s", createUserResponse.Error)
	}

	fmt.Printf("User %s created successfully\n", createUserResponse.Username)
	return nil
}

func commandLogin(cfg *config, args []string) error {
	// Command handles logging in, saves the JWT in the config struct
	if len(args) < 2 {
		return fmt.Errorf("invalid number of arguments")
	}
	url := fmt.Sprintf("%s/api/login", cfg.serverAddress)

	type reqBodyType struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	reqBody := reqBodyType{
		Username: args[0],
		Email:    args[1],
	}
	dat, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	payload := bytes.NewBuffer(dat)

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type loginResponseType struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Username     string    `json:"username"`
		Email        string    `json:"email"`
		JWT          string    `json:"jwt"`
		RefreshToken string    `json:"refresh_token"`
		Error        string    `json:"error"`
	}
	loginResponse := loginResponseType{}
	err = json.Unmarshal(respBody, &loginResponse)
	if err != nil {
		return err
	}

	// If the REST api didn't respond with 202 it means something went wrong and
	// there should be an error message waiting in the response body
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error processing request: %s", loginResponse.Error)
	}

	// We save the jwt in our cfg struct
	cfg.jwt = loginResponse.JWT
	cfg.username = loginResponse.Username
	cfg.email = loginResponse.Email

	fmt.Printf("Logged in as %s\n", loginResponse.Username)

	return nil
}
