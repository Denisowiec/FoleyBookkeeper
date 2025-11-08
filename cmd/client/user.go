package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func commandCreateUser(cfg *config, args []string) error {
	// Command handles creating new users
	if len(args) < 3 {
		return fmt.Errorf("invalid number of arguments")
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

	resp, err := sendRequest(reqBody, "POST", url, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type createUserResponseType struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Error     string    `json:"error"`
	}

	createUserResponse := createUserResponseType{}
	err = processResponse(resp, &createUserResponse)
	if err != nil {
		return err
	}

	// If the REST api respond with something other than 201, somnething went wrong
	// There should be an error message in the response payload
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf(createUserResponse.Error)
	}

	fmt.Printf("User %s created successfully\n", createUserResponse.Username)
	return nil
}

func commandLogin(cfg *config, args []string) error {
	// Command handles logging in, saves the JWT in the config struct
	if len(args) < 2 {
		// If there are no arguments the command displays login status:
		if cfg.jwt == "" {
			fmt.Println("Not logged in.")
			return nil
		} else {
			fmt.Printf("Logged in as %s.\n", cfg.username)
			return nil
		}
	}
	url := fmt.Sprintf("%s/api/login", cfg.serverAddress)

	type reqBodyType struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	reqBody := reqBodyType{
		Email:    args[0],
		Password: args[1],
	}

	resp, err := sendRequest(reqBody, "POST", url, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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

	err = processResponse(resp, &loginResponse)
	if err != nil {
		return err
	}

	// If the REST api didn't respond with 202 it means something went wrong and
	// there should be an error message waiting in the response body
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf(loginResponse.Error)
	}

	// We save the jwt in our cfg struct
	cfg.jwt = loginResponse.JWT
	cfg.username = loginResponse.Username
	cfg.email = loginResponse.Email

	fmt.Printf("Logged in as %s.\n", loginResponse.Username)

	return nil
}

func commandUpdateUser(cfg *config, args []string) error {
	// This function updates username and email address, but not the password
	if len(args) < 2 {
		return fmt.Errorf("invalid number of arguments")
	}
	type reqBodyType struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	reqBody := reqBodyType{
		Username: args[0],
		Email:    args[1],
	}

	url := fmt.Sprintf("%s/api/users", cfg.serverAddress)

	// PUT is the difference between updating user info and creating a new user
	// We need to attach the jwt to authorize ourselves for this operation
	resp, err := sendRequest(reqBody, "PUT", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type updateUserResponseType struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Error     string    `json:"error"`
	}

	updateUserResponse := updateUserResponseType{}

	err = processResponse(resp, &updateUserResponse)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf(updateUserResponse.Error)
	}

	cfg.username = updateUserResponse.Username
	cfg.email = updateUserResponse.Email

	fmt.Printf("Username changed to %s,\nEmail changed to %s.\n", cfg.username, cfg.email)
	return nil
}

func commandUpdatePassword(cfg *config, args []string) error {
	// This function updates the user's password
	if len(args) < 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	type reqBodyType struct {
		Password string `json:"password"`
	}
	reqBody := reqBodyType{
		Password: args[0],
	}

	url := fmt.Sprintf("%s/api/users", cfg.serverAddress)

	resp, err := sendRequest(reqBody, "PUT", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type updatePasswordResponseType struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Error     string    `json:"error"`
	}
	updatePasswordResponse := updatePasswordResponseType{}

	err = processResponse(resp, &updatePasswordResponse)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf(updatePasswordResponse.Error)
	}

	fmt.Println("Password changed successfully.")
	return nil
}
